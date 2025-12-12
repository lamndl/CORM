package backend
import (
    "fmt"
    "github.com/notnil/chess"
)

// ApplyMoveSAN takes a FEN and a SAN move string, returns the new FEN.
func ApplyMoveSAN(fen, san string) (string, error) {
    pos, err := chess.FEN(fen)
    if err != nil {
        return "", fmt.Errorf("invalid FEN: %w", err)
    }

    game := chess.NewGame(pos, chess.UseNotation(chess.AlgebraicNotation{}))

    move := game.MoveStr(san)
    if move != nil {
        return "", fmt.Errorf("invalid SAN move: %s", san)
    }
    return game.Position().String(), nil
}
