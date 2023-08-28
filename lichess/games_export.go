package lichess

import (
	"context"
	"fmt"
	"net/http"
)

// ExportOptions specifies parameters for GamesService.ExportById
// and GamesService.ExportByUsername methods.
type ExportOptions struct {
	Moves     *bool   `url:"moves,omitempty"`
	PgnInJson *bool   `url:"pgnInJson,omitempty"`
	Tags      *bool   `url:"tags,omitempty"`
	Clocks    *bool   `url:"clocks,omitempty"`
	Evals     *bool   `url:"evals,omitempty"`
	Accuracy  *bool   `url:"accuracy,omitempty"`
	Opening   *bool   `url:"opening,omitempty"`
	Literate  *bool   `url:"literate,omitempty"`
	Players   *string `url:"players,omitempty"`
}

func (s *GamesService) ExportById(ctx context.Context, id string, opts *ExportOptions) (*Game, *Response, error) {
	u := fmt.Sprintf("game/export/%v", id)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var game *Game
	resp, err := s.client.Do(req, &game)
	if err != nil {
		return nil, resp, err
	}

	return game, resp, nil
}

func (s *GamesService) ExportByUsername(ctx context.Context, username string, opts *ExportOptions) (*Game, *Response, error) {
	u := fmt.Sprintf("api/user/%v/current-game", username)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var game *Game
	resp, err := s.client.Do(req, &game)
	if err != nil {
		return nil, resp, err
	}

	return game, resp, nil
}
