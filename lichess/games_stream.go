package lichess

import (
	"bufio"
	"bytes"
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

	req, err := s.client.NewRequest(ctx, http.MethodGet, u)
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

// StreamUserGames streams [Game] played by the given username.
// It closes the channel of [Game] and the [Response] body when the context is done.
// So, please use the [context.Context] argument to control the lifetime of the stream.
// Equivalent to [GamesService.ExportByUsername] but streams the response.
// Find more details at https://lichess.org/api#tag/Games/operation/apiGamesUser.
func (s *GamesService) StreamUserGames(
	ctx context.Context,
	username string,
	opts *ExportByUsernameOptions,
) (chan *Game, *Response, error) {
	u := fmt.Sprintf("api/games/user/%v", username)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, u)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	ch := make(chan *Game)

	go s.streamUserGames(ctx, ch, resp)

	return ch, resp, nil
}

func (s *GamesService) streamUserGames(ctx context.Context, ch chan *Game, resp *Response) {
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

		var game Game
		if err := json.Unmarshal([]byte(scanner.Text()), &game); err != nil {
			// We might want to revisit the error handling here
			continue
		}

		select {
		case <-ctx.Done():
			return
		case ch <- &game:
		}
	}

	/*
		if scanner.Err() != nil {
			// We might want to revisit the
			// use of scanner.Err().
		}
	*/
}

// StreamGamesOfUsersOptions specifies parameters for
// GamesService.StreamGamesOfUsers method.
type StreamGamesOfUsersOptions struct {
	CurrentGames *bool `url:"withCurrentGames,omitempty"`
}

// GameStream is a Lichess [Game] representation for streaming methods.
// For instance: [GamesService.StreamGamesOfUsers].
type GameStream struct {
	Id         string      `json:"id"`
	Rated      bool        `json:"rated"`
	Variant    GameVariant `json:"variant"`
	Speed      GameSpeed   `json:"speed"`
	Perf       string      `json:"perf"`
	CreatedAt  int64       `json:"createdAt"`
	Status     int64       `json:"status"`
	StatusName GameStatus  `json:"statusName"`
	Clock      GameClock   `json:"clock"`
	Players    struct {
		White GameUser `json:"white"`
		Black GameUser `json:"black"`
	} `json:"players"`
}

// StreamGamesOfUsers streams [Game] played among the given usernames.
// It closes the channel of [Game] and the [Response] body when the context is done.
// So, please use the [context.Context] argument to control the lifetime of the stream.
// Find more details at https://lichess.org/api#tag/Games/operation/gamesByUsers.
func (s *GamesService) StreamGamesOfUsers(
	ctx context.Context,
	usernames []string,
	opts *StreamGamesOfUsersOptions,
) (chan *GameStream, *Response, error) {
	u, err := addOptions("api/stream/games-by-users", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequestWithBody(ctx, http.MethodPost, u, RequestBody{
		Bytes: bytes.NewReader([]byte(strings.Join(usernames, ","))),
		Type:  "text/plain",
	})
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, resp, err
	}

	ch := make(chan *GameStream)

	go s.streamGamesOfUsers(ctx, ch, resp)

	return ch, resp, nil
}

func (s *GamesService) streamGamesOfUsers(ctx context.Context, ch chan *GameStream, resp *Response) {
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

		var gameStream GameStream
		if err := json.Unmarshal([]byte(scanner.Text()), &gameStream); err != nil {
			// We might want to revisit the error handling here
			continue
		}

		select {
		case <-ctx.Done():
			return
		case ch <- &gameStream:
		}
	}

	/*
		if scanner.Err() != nil {
			// We might want to revisit the
			// use of scanner.Err().
		}
	*/
}
