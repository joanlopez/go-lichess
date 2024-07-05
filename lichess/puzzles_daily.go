package lichess

import (
	"context"
	"net/http"
)

// DailyPuzzle represents a Lichess puzzle.
type DailyPuzzle struct {
	Game   *DailyPuzzleGame `json:"game,omitempty"`
	Puzzle *Puzzle          `json:"puzzle,omitempty"`
}

// DailyPuzzleGame represents a Lichess daily puzzle game.
type DailyPuzzleGame struct {
	Clock   string                    `json:"clock,omitempty"`
	Id      string                    `json:"id,omitempty"`
	Perf    *DailyPuzzleGamePerf      `json:"perf,omitempty"`
	Pgn     string                    `json:"pgn,omitempty"`
	Players [2]*DailyPuzzleGamePlayer `json:"players,omitempty"`
	Rated   bool                      `json:"rated,omitempty"`
}

// DailyPuzzleGamePerf represents a Lichess daily puzzle game perf.
type DailyPuzzleGamePerf struct {
	// Key is one of GameSpeed or GameVariant.
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}

// DailyPuzzleGamePlayer represents a Lichess daily puzzle game player.
type DailyPuzzleGamePlayer struct {
	Color  string  `json:"color,omitempty"`
	Flair  *string `json:"flair,omitempty"`
	Id     string  `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Patron *bool   `json:"patron,omitempty"`
	Rating int     `json:"rating,omitempty"`
	Title  *string `json:"title,omitempty"` // TODO: Define enum.
}

func (s *PuzzlesService) GetDailyPuzzle(
	ctx context.Context,
) (*DailyPuzzle, *Response, error) {
	u := "api/puzzle/daily"

	req, err := s.client.NewRequest(ctx, http.MethodGet, u)
	if err != nil {
		return nil, nil, err
	}

	var puzzle *DailyPuzzle
	resp, err := s.client.Do(req, &puzzle)
	if err != nil {
		return nil, resp, err
	}

	return puzzle, resp, nil
}
