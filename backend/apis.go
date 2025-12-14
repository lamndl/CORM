package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Structs matching Explorer API response
type Move struct {
	UCI   string `json:"uci"`
	SAN   string `json:"san"`
	White int    `json:"white"`
	Black int    `json:"black"`
	Draws int    `json:"draws"`
}

type Opening struct {
	ECO  string `json:"eco"`
	Name string `json:"name"`
}

type ExplorerResponse struct {
	White   int     `json:"white"`
	Black   int     `json:"black"`
	Draws   int     `json:"draws"`
	Moves   []Move  `json:"moves"`
	Opening Opening `json:"opening"`
}

// Fetch data from Lichess Explorer
func FetchExplorerData(fen string, elo int) (ExplorerResponse, error) {
	url := fmt.Sprintf("https://explorer.lichess.ovh/lichess?variant=standard&speeds=rapid&ratings=%d&fen=%s", elo, url.QueryEscape(fen))

	resp, err := http.Get(url)
	if err != nil {
		return ExplorerResponse{}, err
	}
	defer resp.Body.Close()

	var data ExplorerResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ExplorerResponse{}, err
	}
	return data, nil
}
