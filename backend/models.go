// backend/manager/models.go
package backend

type Repertoire struct {
    ID       int64   `json:"id"`
    Name     string  `json:"name"`
    Color    string  `json:"color"`    // "white" | "black"
    Elo      int     `json:"elo"`
    Coverage float64 `json:"coverage"`
}

// MoveWinrate is a simplified view of each move with totals, winrates, and chance
type MoveWinrate struct {
    SAN       string  `json:"san"`
    UCI       string  `json:"uci"`
    Total     int     `json:"total"`
    WhiteRate float64 `json:"whiteRate"`
    BlackRate float64 `json:"blackRate"`
    DrawRate  float64 `json:"drawRate"`
    Chance    float64 `json:"chance"` // % chance this move is played
}

type PositionWinrate struct {
    Total     int           `json:"total"`
    WhiteRate float64       `json:"whiteRate"`
    BlackRate float64       `json:"blackRate"`
    DrawRate  float64       `json:"drawRate"`
    Moves     []MoveWinrate `json:"moves"`
}

