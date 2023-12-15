package lichess

import (
	"context"
	"net/http"
)

// GetPuzzleActivityOptions specifies parameters for
// PuzzlesService.GetPuzzleActivity method.
type GetPuzzleActivityOptions struct {
	// Max specifies how many entries to download.
	// Leave empty to download all activity.
	Max *int `url:"max,omitempty"` // >= 1
	// Before is used to download entries before this timestamp.
	// Defaults to now. Use before and max for pagination.
	Before *int `url:"before,omitempty"` // >= 1356998400070
}

func (s *PuzzlesService) GetPuzzleActivity(
	ctx context.Context,
	opts *GetPuzzleActivityOptions,
) ([]*PuzzleRound, *Response, error) {
	u := "api/puzzle/activity"
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u)
	if err != nil {
		return nil, nil, err
	}

	var rounds []*PuzzleRound
	resp, err := s.client.Do(req, &rounds)
	if err != nil {
		return nil, resp, err
	}

	return rounds, resp, nil
}
