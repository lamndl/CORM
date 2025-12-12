import React, { useEffect, useState } from "react";
import * as RepMgr from "../wailsjs/go/backend/RepertoireManager";

export default function MoveList({ setFen, selectedRepID, refreshRepertoires, setRefreshMoveList }) {
  const [moves, setMoves] = useState([]);
  const [total, setTotal] = useState(0);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [existingEdges, setExistingEdges] = useState([]);

  async function refreshMoves() {
    setIsRefreshing(true);
    try {
      const pos = await RepMgr.GetCurrentWinrates();
      console.log(pos);
      setMoves(pos.moves);
      setTotal(pos.total);

      // Fetch existing edges
      const edges = await RepMgr.ListEdges();
      setExistingEdges(edges || []); // Handle empty edges by defaulting to an empty array
    } catch (e) {
      console.error("Failed to fetch winrates or edges:", e);
    } finally {
      setIsRefreshing(false);
    }
  }

  useEffect(() => {
    if (selectedRepID) {
      refreshMoves();
    }
  }, [selectedRepID]);

  useEffect(() => {
    console.log("MoveList refresh triggered");
    setRefreshMoveList(() => refreshMoves); // Provide refresh function to parent
  }, [setRefreshMoveList]);

  async function addEdge(san) {
    if (isRefreshing) return; // Prevent adding edges while refreshing
    try {
      await RepMgr.AddEdge(san); // ✅ backend insert
      const newFEN = await RepMgr.GetCurrentFEN(); // ✅ fetch updated FEN
      setFen(newFEN); // ✅ update board
      await refreshMoves(); // ✅ reload UI
      refreshRepertoires(); // ✅ refresh repertoire list to update due count
    } catch (err) {
      console.error("Failed to add edge:", err);
    }
  }

  async function deleteEdge(san) {
    if (isRefreshing) return; // Prevent deleting edges while refreshing
    try {
      await RepMgr.DeleteEdge(san); // ✅ backend delete
      await refreshMoves(); // ✅ reload UI
    } catch (err) {
      console.error("Failed to delete edge:", err);
    }
  }

  async function playMove(san) {
    if (isRefreshing) return; // Prevent playing moves while refreshing
    try {
      await RepMgr.PlayMoveSAN(san); // ✅ backend play move
      const newFEN = await RepMgr.GetCurrentFEN(); // ✅ fetch updated FEN
      setFen(newFEN); // ✅ update board
      await refreshMoves(); // ✅ reload UI
    } catch (err) {
      console.error("Failed to play move:", err);
    }
  }

  async function resetToStart() {
    try {
      const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"; // Default starting position
      await RepMgr.SetCurrentFEN(startFEN); // ✅ Update backend state
      setFen(startFEN); // ✅ Update frontend state
      await refreshMoves(); // ✅ Reload moves
    } catch (err) {
      console.error("Failed to reset to starting position:", err);
    }
  }

  if (!selectedRepID) {
    return <p style={{ padding: "1rem" }}>Please select a repertoire to view moves.</p>;
  }

  return (
    <div style={{ padding: "0.5rem", overflowY: "auto", height: "100%" }}>
      <h4>Moves from current position</h4>
      <button onClick={resetToStart} style={{ marginBottom: "1rem" }} disabled={isRefreshing}>
        Back to Starting Position
      </button>
      <table style={{ width: "100%", borderCollapse: "collapse" }}>
        <thead>
          <tr>
            <th style={{ textAlign: "left" }}>Move</th>
            <th>Chance</th>
            <th>White%</th>
            <th>Black%</th>
            <th>Draw%</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {moves.map((m, i) => (
            <tr key={i}>
              <td>{m.san}</td>
              <td>{m.chance.toFixed(1)}%</td>
              <td>{m.whiteRate.toFixed(1)}%</td>
              <td>{m.blackRate.toFixed(1)}%</td>
              <td>{m.drawRate.toFixed(1)}%</td>
              <td>
                {existingEdges.includes(m.san) ? (
                  <>
                    <button
                      onClick={() => playMove(m.san)}
                      disabled={isRefreshing}
                    >
                      Move
                    </button>
                    <button
                      onClick={() => deleteEdge(m.san)}
                      disabled={isRefreshing}
                      style={{ marginLeft: "0.5rem" }}
                    >
                      Delete
                    </button>
                  </>
                ) : (
                  <button
                    onClick={() => addEdge(m.san)}
                    disabled={isRefreshing}
                  >
                    Add
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}