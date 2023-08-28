package lichess

import (
	"context"
	"fmt"
	"net/http"
)

// ExportOptions specifies parameters for GamesService.ExportById
// GamesService.ExportCurrent, GamesService.ExportByUsername methods.
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

// ExportByUsernameOptions specifies parameters for
// GamesService.ExportByUsername method.
type ExportByUsernameOptions struct {
	ExportOptions
	Since    *int    `url:"since,omitempty"`
	Until    *int    `url:"until,omitempty"`
	Max      *int    `url:"max,omitempty"`
	VS       *string `url:"vs,omitempty"`
	Rated    *bool   `url:"rated,omitempty"`
	PerfType *string `url:"perfType,omitempty"`
	Color    *string `url:"color,omitempty"`
	Analysed *string `url:"analysed,omitempty"`
	Ongoing  *bool   `url:"ongoing,omitempty"`
	Finished *bool   `url:"finished,omitempty"`
	LastFen  *bool   `url:"lastFen,omitempty"`
	Sort     *string `url:"sort,omitempty"` // Either "dateAsc" or "dateDesc"
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

func (s *GamesService) ExportCurrent(ctx context.Context, username string, opts *ExportOptions) (*Game, *Response, error) {
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

func (s *GamesService) ExportByUsername(ctx context.Context, username string, opts *ExportByUsernameOptions) ([]*Game, *Response, error) {
	u := fmt.Sprintf("api/games/user/%v", username)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	var games []*Game
	resp, err := s.client.Do(req, &games)
	if err != nil {
		return nil, resp, err
	}

	return games, resp, nil
}
