package lichess

// PuzzlesService handles communication with the puzzle related
// methods of the Lichess API.
//
// Lichess API docs: https://lichess.org/api#tag/Puzzles
type PuzzlesService service

// PuzzleRound represents a Lichess puzzle round.
type PuzzleRound struct {
	Date   int    `json:"date,omitempty"`
	Win    bool   `json:"win,omitempty"`
	Puzzle Puzzle `json:"puzzle,omitempty"`
}

// Puzzle represents a Lichess puzzle.
type Puzzle struct {
	Id         string   `json:"id,omitempty"`
	InitialPly int      `json:"initialPly,omitempty"`
	Fen        string   `json:"fen,omitempty"`
	Plays      int      `json:"plays,omitempty"`
	Rating     int      `json:"rating,omitempty"`
	Solution   []string `json:"solution,omitempty"`
	Themes     []string `json:"themes,omitempty"`
}
