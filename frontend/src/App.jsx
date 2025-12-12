import React, { useState } from "react";
import ChessboardComponent from "./Chessboard";
import Repertoires from "./Repertoires";
import MoveList from "./MoveList";

export default function App() {
    const [selectedRepID, setSelectedRepID] = useState(null);
    const [fen, setFen] = useState("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");       // ✅ track current FEN
    const [moves, setMoves] = useState([]);   // ✅ track move list
    const [refreshRepertoires, setRefreshRepertoires] = useState(() => () => { }); // Callback to refresh repertoires
    const [refreshMoveList, setRefreshMoveList] = useState(() => () => {}); // Callback to refresh MoveList
    const [isPracticeMode, setIsPracticeMode] = useState(false); // Track practice mode

    function resetFen() {
        setFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"); // Reset to starting position
    }

    function handlePracticeModeChange(active) {
        setIsPracticeMode(active);
        console.log("Practice mode changed:", active);
        if (!active) {
            console.log("Refreshing MoveList...");
            refreshMoveList(); // Refresh MoveList when practice ends
        }
    }

    return (
        <div style={{ display: "flex", height: "100vh" }}>
            {/* Left: Chessboard */}
            <div style={{ flex: 2, borderRight: "1px solid #ccc" }}>
                <ChessboardComponent fen={fen} />
            </div>

            {/* Right: Sidebar */}
            <div style={{ flex: 1, display: "flex", flexDirection: "column" }}>
                {/* Top of sidebar: Repertoire list */}
                <div style={{ borderBottom: "1px solid #ccc", padding: "1rem" }}>
                    <Repertoires
                        onSelect={setSelectedRepID}
                        onResetFen={resetFen}
                        setRefreshRepertoires={setRefreshRepertoires} // Pass callback setter
                        setFen={setFen}
                        onPracticeModeChange={handlePracticeModeChange}
                    />
                </div>
                {/* Conditionally render MoveList based on practice mode */}
                {!isPracticeMode && (
                    <MoveList
                        setFen={setFen}
                        selectedRepID={selectedRepID}
                        refreshRepertoires={refreshRepertoires} // Pass refresh callback
                        setRefreshMoveList={setRefreshMoveList} // Pass MoveList refresh callback
                    />
                )}
            </div>
        </div>
    );
}