package lichess

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GameStreamEventType represents a Lichess game stream event type.
type GameStreamEventType string

const (
	GameDescriptionEventType GameStreamEventType = "description"
	GameMoveEventType        GameStreamEventType = "move"
	GameStreamErrorEventType GameStreamEventType = "error"
)

// GameStreamEvent represents a Lichess game stream event.
type GameStreamEvent interface {
	GameStreamEventType() GameStreamEventType
}

// GameDescription the description of a Lichess game
// sent at the beginning and end of stream.
type GameDescription struct {
	Id      string `json:"id"`
	Variant struct {
		Key GameVariant `json:"key"`
	} `json:"variant"`
	Speed         GameSpeed `json:"speed"`
	Perf          string    `json:"perf"`
	Rated         bool      `json:"rated"`
	InitialFen    string    `json:"initialFen"`
	Fen           string    `json:"fen"`
	Player        string    `json:"player"`
	Turns         int       `json:"turns"`
	StartedAtTurn int       `json:"startedAtTurn"`
	Source        string    `json:"source"`
	Status        struct {
		Name GameStatus `json:"name"`
	} `json:"status"`
	CreatedAt int64  `json:"createdAt"`
	LastMove  string `json:"lastMove"`
	Players   struct {
		White GameUser `json:"white"`
		Black GameUser `json:"black"`
	} `json:"players"`
}

func (e GameDescription) GameStreamEventType() GameStreamEventType {
	return GameDescriptionEventType
}

type GameMove struct {
	Fen string `json:"fen"`
	LM  string `json:"lm"`
	WC  int    `json:"wc"`
	BC  int    `json:"bc"`
}

func (e GameMove) GameStreamEventType() GameStreamEventType {
	return GameMoveEventType
}

// GameStreamEventError represents a Lichess game stream event error.
type GameStreamEventError struct {
	error
}

func (e GameStreamEventError) GameStreamEventType() GameStreamEventType {
	return GameStreamErrorEventType
}

// StreamGameMoves streams [GameStreamEvent] happening at [Game] identified by id.
// It closes the channel of [GameStreamEvent] and the [Response] body when the context is done.
// So, please use the [context.Context] argument to control the lifetime of the stream.
// Find more details at https://lichess.org/api#tag/Games/operation/streamGame.
func (s *GamesService) StreamGameMoves(ctx context.Context, id string) (chan GameStreamEvent, *Response, error) {
	u := fmt.Sprintf("api/stream/game/%v", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	ch := make(chan GameStreamEvent)

	go s.streamGameMoves(ctx, ch, resp)

	return ch, resp, nil
}

func (s *GamesService) streamGameMoves(ctx context.Context, ch chan GameStreamEvent, resp *Response) {
	defer func() {
		// Explicit ignore error.
		// We might want to revisit this later.
		_ = resp.Body.Close()
		close(ch)
	}()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case ch <- s.parseGameStreamEvent(scanner.Text()):
		}
	}

	if scanner.Err() != nil {
		select {
		case <-ctx.Done():
		case ch <- GameStreamEventError{error: scanner.Err()}:
		}
	}
}

func (s *GamesService) parseGameStreamEvent(event string) GameStreamEvent {
	switch {
	case strings.Contains(event, `"wc":`):
		var gameMove GameMove
		if err := json.Unmarshal([]byte(event), &gameMove); err != nil {
			return GameStreamEventError{error: err}
		}
		return gameMove

	case strings.Contains(event, `"variant":`):
		var gameDescription GameDescription
		if err := json.Unmarshal([]byte(event), &gameDescription); err != nil {
			return GameStreamEventError{error: err}
		}
		return gameDescription

	default:
		return GameStreamEventError{error: fmt.Errorf(event)}
	}
}
