import React from "react";
import { ChessBoard } from "react-fen-chess-board";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";

export default function ChessboardComponent({ fen }) {
  // fallback to starting position if no FEN is provided
  const startingFEN =
    "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";

  return (
    <DndProvider backend={HTML5Backend}>
      <div style={{ width: "100%", height: "100%" }}>
        <ChessBoard
          key={fen || startingFEN} // force remount when FEN changes
          fen={fen || startingFEN} // ✅ board updates when fen prop changes
          rotated={false} // show from White’s perspective
          onMove={() => {}} // noop handler to avoid passing undefined
        />
      </div>
    </DndProvider>
  );
}