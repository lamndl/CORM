package backend

import (
	"context"
	"database/sql"
	"fmt"
)

type RepertoireManager struct {
	db          *sql.DB
	selectedRep int64
	currentFEN  string
}

func NewRepertoireManager(db *sql.DB) *RepertoireManager {
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	return &RepertoireManager{
		db:          db,
		selectedRep: 1,        // no repertoire selected yet
		currentFEN:  startFEN, // ✅ default starting position
	}
}

// Create a new repertoire
func (m *RepertoireManager) Create(name, color string, elo int) (int64, error) {
	// Insert repertoire row
	res, err := m.db.ExecContext(context.Background(),
		`INSERT INTO repertoire (name, color, elo, coverage) VALUES (?, ?, ?, 0.0)`,
		name, color, elo)
	if err != nil {
		return 0, err
	}
	repID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Insert the start node (initial chess position FEN)
	startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	_, err = m.db.ExecContext(context.Background(),
		`INSERT INTO nodes (fen, rep_id, sr_index, due, last_review) VALUES (?, ?, 0, NULL, NULL)`,
		startFEN, repID)
	if err != nil {
		return 0, err
	}

	return repID, nil
}

// List all repertoires
func (m *RepertoireManager) List() ([]Repertoire, error) {
	rows, err := m.db.QueryContext(context.Background(),
		`SELECT id, name, color, elo, coverage FROM repertoire ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reps := make([]Repertoire, 0)
	for rows.Next() {
		var r Repertoire
		if err := rows.Scan(&r.ID, &r.Name, &r.Color, &r.Elo, &r.Coverage); err != nil {
			return nil, err
		}
		reps = append(reps, r)
	}
	return reps, rows.Err()
}

// Update an existing repertoire
func (m *RepertoireManager) Update(r Repertoire) error {
	_, err := m.db.ExecContext(context.Background(),
		`UPDATE repertoire SET name=?, color=?, elo=?, coverage=? WHERE id=?`,
		r.Name, r.Color, r.Elo, r.Coverage, r.ID)
	return err
}

// Delete a repertoire
func (m *RepertoireManager) Delete(id int64) error {
	_, err := m.db.ExecContext(context.Background(),
		`DELETE FROM repertoire WHERE id=?`, id)
	return err
}

// Set the currently selected repertoire ID
func (m *RepertoireManager) SetCurrentID(id int64) {
	m.selectedRep = id
}

// Get the currently selected repertoire ID
func (m *RepertoireManager) GetCurrentID() int64 {
	return m.selectedRep
}

// Select a repertoire and set its start node
func (m *RepertoireManager) SelectRepertoire(id int64) {
	m.selectedRep = id
	m.currentFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
}

// Get current selected repertoire ID
func (m *RepertoireManager) GetSelectedID() int64 {
	return m.selectedRep
}

// Get current FEN
func (m *RepertoireManager) GetCurrentFEN() string {
	return m.currentFEN
}

// Set current FEN (e.g. after extending)
func (m *RepertoireManager) SetCurrentFEN(fen string) {
	m.currentFEN = fen
}

func (m *RepertoireManager) GetCurrentWinrates() (PositionWinrate, error) {
	fen := m.currentFEN
	if fen == "" {
		return PositionWinrate{}, fmt.Errorf("no current FEN set")

	}
	elo, err := m.GetCurrentElo()
	if err != nil {
		return PositionWinrate{}, fmt.Errorf("failed to get current elo: %w", err)
	}
	data, err := FetchExplorerData(fen, elo)
	if err != nil {
		return PositionWinrate{}, err
	}

	pos := PositionWinrate{}
	pos.Total = data.White + data.Black + data.Draws
	if pos.Total > 0 {
		pos.WhiteRate = float64(data.White) / float64(pos.Total) * 100
		pos.BlackRate = float64(data.Black) / float64(pos.Total) * 100
		pos.DrawRate = float64(data.Draws) / float64(pos.Total) * 100
	}

	// coverage, err := m.GetCurrentRepCoverage() // Fetch repertoire coverage
	// if err != nil {
	// 	return PositionWinrate{}, fmt.Errorf("failed to get repertoire coverage: %w", err)
	// }
	for _, m := range data.Moves {
		total := m.White + m.Black + m.Draws
		mw := MoveWinrate{
			SAN:   m.SAN,
			UCI:   m.UCI,
			Total: total,
		}
		if total > 0 {
			mw.WhiteRate = float64(m.White) / float64(total) * 100
			mw.BlackRate = float64(m.Black) / float64(total) * 100
			mw.DrawRate = float64(m.Draws) / float64(total) * 100
			mw.Chance = float64(total) / float64(pos.Total) * 100
		}
		// fmt.Println(mw.Chance, coverage)
		// if mw.Chance >= coverage {
		// 	pos.Moves = append(pos.Moves, mw)
		// }
		pos.Moves = append(pos.Moves, mw)
	}
	return pos, nil
}

func (m *RepertoireManager) PlayMoveSAN(moveSAN string) error {
	if m.selectedRep == 0 {
		return fmt.Errorf("no repertoire selected")
	}
	if m.currentFEN == "" {
		return fmt.Errorf("no current FEN set")
	}

	childFEN, err := ApplyMoveSAN(m.currentFEN, moveSAN)
	if err != nil {
		return err
	}

	m.currentFEN = childFEN
	return nil
}

// Edge represents an outgoing move (edge) from a position in a repertoire.
type Edge struct {
	RepID     int64
	ParentFEN string
	ChildFEN  string
	MoveSAN   string
}

func (m *RepertoireManager) AddEdge(moveSAN string) error {
	if m.selectedRep == 0 {
		return fmt.Errorf("no repertoire selected")
	}
	if m.currentFEN == "" {
		return fmt.Errorf("no current FEN set")
	}
	childFEN, err := ApplyMoveSAN(m.currentFEN, moveSAN)
	if err != nil {
		return err
	}

	// Insert child node into `nodes` table (ignore if already exists)
	_, err = m.db.ExecContext(context.Background(),
		`INSERT OR IGNORE INTO nodes (fen, rep_id, sr_index, due, last_review) VALUES (?, ?, 0, NULL, NULL)`,
		childFEN, m.selectedRep)
	if err != nil {
		return fmt.Errorf("failed to insert child node: %w", err)
	}

	// Insert edge into `edges` table
	_, err = m.db.ExecContext(context.Background(),
		`INSERT INTO edges (rep_id, parent_fen, child_fen, move) VALUES (?, ?, ?, ?)`,
		m.selectedRep, m.currentFEN, childFEN, moveSAN)
	if err != nil {
		return fmt.Errorf("failed to insert edge: %w", err)
	}

	// Update parent node's deadline to current time and reset sr_index to 0
	_, err = m.db.ExecContext(context.Background(),
		`UPDATE nodes SET due = CURRENT_TIMESTAMP, sr_index = 0 WHERE rep_id = ? AND fen = ?`,
		m.selectedRep, m.currentFEN)
	if err != nil {
		return fmt.Errorf("failed to update parent node: %w", err)
	}

	// Advance current position to the child (consistent with PlayMoveSAN behavior)
	if err := m.PlayMoveSAN(moveSAN); err != nil {
		return fmt.Errorf("failed to play move after adding edge: %w", err)
	}
	return nil
}

// ListEdges returns SAN moves (strings) from the current position in the selected repertoire.
func (m *RepertoireManager) ListEdges() ([]string, error) {
	if m.selectedRep == 0 {
		return nil, fmt.Errorf("no repertoire selected")
	}
	if m.currentFEN == "" {
		return nil, fmt.Errorf("no current FEN set")
	}

	rows, err := m.db.QueryContext(context.Background(),
		`SELECT move FROM edges WHERE rep_id = ? AND parent_fen = ?`,
		m.selectedRep, m.currentFEN)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var moves []string
	for rows.Next() {
		var mv string
		if err := rows.Scan(&mv); err != nil {
			return nil, err
		}
		moves = append(moves, mv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return moves, nil
}

func (m *RepertoireManager) DeleteEdge(moveSAN string) error {
	if m.selectedRep == 0 {
		return fmt.Errorf("no repertoire selected")
	}
	if m.currentFEN == "" {
		return fmt.Errorf("no current FEN set")
	}

	var childFEN string
	err := m.db.QueryRowContext(context.Background(),
		`SELECT child_fen FROM edges WHERE rep_id = ? AND parent_fen = ? AND move = ? LIMIT 1`,
		m.selectedRep, m.currentFEN, moveSAN).Scan(&childFEN)
	if err == sql.ErrNoRows {
		return fmt.Errorf("edge not found")
	}
	if err != nil {
		return err
	}

	_, err = m.db.ExecContext(context.Background(),
		`DELETE FROM edges WHERE rep_id = ? AND parent_fen = ? AND move = ?`,
		m.selectedRep, m.currentFEN, moveSAN)
	if err != nil {
		return err
	}

	// Remove orphan child node if no other edges reference it (but avoid deleting the standard start position)
	var cnt int
	err = m.db.QueryRowContext(context.Background(),
		`SELECT COUNT(1) FROM edges WHERE rep_id = ? AND child_fen = ?`,
		m.selectedRep, childFEN).Scan(&cnt)
	if err != nil {
		return err
	}
	if cnt == 0 {
		startFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		if childFEN != startFEN {
			_, err = m.db.ExecContext(context.Background(),
				`DELETE FROM nodes WHERE rep_id = ? AND fen = ?`,
				m.selectedRep, childFEN)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *RepertoireManager) GetDueFENs() ([]string, error) {
	if m.selectedRep == 0 {
		return nil, fmt.Errorf("no repertoire selected")
	}

	// Query to fetch FENs of nodes with due date <= current time
	rows, err := m.db.QueryContext(context.Background(),
		`SELECT fen FROM nodes WHERE rep_id = ? AND due <= DATETIME('now')`,
		m.selectedRep)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch due FENs: %w", err)
	}
	defer rows.Close()

	var fens []string
	for rows.Next() {
		var fen string
		if err := rows.Scan(&fen); err != nil {
			return nil, fmt.Errorf("failed to scan FEN: %w", err)
		}
		fens = append(fens, fen)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return fens, nil
}

// CountDueNodes returns the number of due nodes for a given repertoire.
func (m *RepertoireManager) CountDueNodes(repID int64) (int, error) {
	var count int
	err := m.db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM nodes WHERE rep_id = ? AND due <= DATETIME('now')`,
		repID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count due nodes: %w", err)
	}
	return count, nil
}

// TestCurrentPosition validates the SAN move against the current position and updates the Leitner box.
func (m *RepertoireManager) TestCurrentPosition(moveSAN string) error {
	if m.selectedRep == 0 {
		return fmt.Errorf("no repertoire selected")
	}
	if m.currentFEN == "" {
		return fmt.Errorf("no current FEN set")
	}

	// Check if the moveSAN is a valid edge from the current FEN
	var childFEN string
	err := m.db.QueryRowContext(context.Background(),
		`SELECT child_fen FROM edges WHERE rep_id = ? AND parent_fen = ? AND move = ?`,
		m.selectedRep, m.currentFEN, moveSAN).Scan(&childFEN)
	if err == sql.ErrNoRows {
		// Incorrect move: demote Leitner box
		_, updateErr := m.db.ExecContext(context.Background(),
			`UPDATE nodes SET sr_index = MAX(sr_index - 1, 0) WHERE rep_id = ? AND fen = ?`,
			m.selectedRep, m.currentFEN)
		if updateErr != nil {
			return fmt.Errorf("failed to demote Leitner box: %w", updateErr)
		}
		return fmt.Errorf("incorrect move, correct move is: %s", childFEN)
	} else if err != nil {
		return fmt.Errorf("failed to validate move: %w", err)
	}

	// Correct move: promote Leitner box
	_, err = m.db.ExecContext(context.Background(),
		`UPDATE nodes SET sr_index = MIN(sr_index + 1,3) WHERE rep_id = ? AND fen = ?`,
		m.selectedRep, m.currentFEN)
	if err != nil {
		return fmt.Errorf("failed to promote Leitner box: %w", err)
	}

	// Advance to the child position
	m.currentFEN = childFEN
	return nil
}

// TestCurrentPositionWithDueDate validates the SAN move against the current position, updates the Leitner box, and adjusts the due date.
func (m *RepertoireManager) TestCurrentPositionWithDueDate(moveSAN string) error {
	if m.selectedRep == 0 {
		return fmt.Errorf("no repertoire selected")
	}
	if m.currentFEN == "" {
		return fmt.Errorf("no current FEN set")
	}

	// Check if the moveSAN is a valid edge from the current FEN
	var childFEN string
	err := m.db.QueryRowContext(context.Background(),
		`SELECT child_fen FROM edges WHERE rep_id = ? AND parent_fen = ? AND move = ?`,
		m.selectedRep, m.currentFEN, moveSAN).Scan(&childFEN)
	if err == sql.ErrNoRows {
		// Incorrect move: demote Leitner box and reset due date
		_, updateErr := m.db.ExecContext(context.Background(),
			`UPDATE nodes SET sr_index = MAX(sr_index - 1, 0), due = DATETIME('now', '+1 day') WHERE rep_id = ? AND fen = ?`,
			m.selectedRep, m.currentFEN)
		if updateErr != nil {
			return fmt.Errorf("failed to demote Leitner box: %w", updateErr)
		}
		return fmt.Errorf("incorrect move")
	} else if err != nil {
		return fmt.Errorf("failed to validate move: %w", err)
	}

	// Correct move: promote Leitner box and adjust due date based on Leitner index
	_, err = m.db.ExecContext(context.Background(),
		`UPDATE nodes SET sr_index = sr_index + 1, 
		 due = CASE 
		   WHEN sr_index = 0 THEN DATETIME('now', '+1 day')
		   WHEN sr_index = 1 THEN DATETIME('now', '+3 days')
		   WHEN sr_index = 2 THEN DATETIME('now', '+1 week')
		   ELSE DATETIME('now', '+3 weeks')
		 END
		WHERE rep_id = ? AND fen = ?`,
		m.selectedRep, m.currentFEN)
	if err != nil {
		return fmt.Errorf("failed to promote Leitner box: %w", err)
	}

	// Advance to the child position
	m.currentFEN = childFEN
	return nil
}

func (m *RepertoireManager) GetCurrentRepCoverage() (float64, error) {
	if m.selectedRep == 0 {
		return 0.0, fmt.Errorf("no repertoire selected")
	}

	var coverage float64
	err := m.db.QueryRowContext(context.Background(),
		`SELECT coverage FROM repertoire WHERE id = ?`,
		m.selectedRep).Scan(&coverage)
	if err != nil {
		return 0.0, fmt.Errorf("failed to get repertoire coverage: %w", err)
	}
	return 100.0 / float64(coverage), nil
}

// Added a method to get the current repertoire's elo rating.
func (m *RepertoireManager) GetCurrentElo() (int, error) {
	if m.selectedRep == 0 {
		return 0, fmt.Errorf("no repertoire selected")
	}

	var elo int
	err := m.db.QueryRowContext(context.Background(),
		`SELECT elo FROM repertoire WHERE id = ?`,
		m.selectedRep).Scan(&elo)
	if err != nil {
		return 0, fmt.Errorf("failed to get repertoire elo: %w", err)
	}

	return elo, nil
}
