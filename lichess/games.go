package lichess

// GamesService handles communication with the game related
// methods of the Lichess API.
//
// Lichess API docs: https://lichess.org/api#tag/Games
type GamesService service

// Game represents a Lichess game.
type Game struct {
	Id         string      `json:"id,omitempty"`
	Rated      bool        `json:"rated,omitempty"`
	Variant    GameVariant `json:"variant,omitempty"`
	Speed      GameSpeed   `json:"speed,omitempty"`
	Perf       string      `json:"perf,omitempty"`
	CreatedAt  int64       `json:"createdAt,omitempty"`
	LastMoveAt int64       `json:"lastMoveAt,omitempty"`
	Status     GameStatus  `json:"status,omitempty"`
	Players    struct {
		White GameUser `json:"white,omitempty"`
		Black GameUser `json:"black,omitempty"`
	} `json:"players,omitempty"`
	InitialFen  *string         `json:"initialFen,omitempty"`
	Winner      *string         `json:"winner,omitempty"`
	Opening     *GameOpening    `json:"opening,omitempty"`
	Moves       *string         `json:"moves,omitempty"`
	Pgn         *string         `json:"pgn,omitempty"`
	DaysPerTurn *int            `json:"daysPerTurn,omitempty"`
	Analysis    []*GameAnalysis `json:"analysis,omitempty"`
	Tournament  *string         `json:"tournament,omitempty"`
	Swiss       *string         `json:"swiss,omitempty"`
	Clock       *GameClock      `json:"clock,omitempty"`
}

// GameVariant represents a Lichess game variant.
type GameVariant string

const (
	Standard      GameVariant = "standard"
	Chess960      GameVariant = "chess960"
	Crazyhouse    GameVariant = "crazyhouse"
	Antichess     GameVariant = "antichess"
	Atomic        GameVariant = "atomic"
	Horde         GameVariant = "horde"
	KingOfTheHill GameVariant = "kingOfTheHill"
	RacingKings   GameVariant = "racingKings"
	ThreeCheck    GameVariant = "threeCheck"
	FromPosition  GameVariant = "fromPosition"
)

// GameSpeed represents a Lichess game speed.
type GameSpeed string

const (
	UltraBullet    GameSpeed = "ultraBullet"
	Bullet         GameSpeed = "bullet"
	Blitz          GameSpeed = "blitz"
	Rapid          GameSpeed = "rapid"
	Classical      GameSpeed = "classical"
	Correspondence GameSpeed = "correspondence"
)

// GameStatus represents a Lichess game status.
type GameStatus string

const (
	Created       GameStatus = "created"
	Started       GameStatus = "started"
	Aborted       GameStatus = "aborted"
	Mate          GameStatus = "mate"
	Resign        GameStatus = "resign"
	Stalemate     GameStatus = "stalemate"
	Timeout       GameStatus = "timeout"
	Draw          GameStatus = "draw"
	OutOfTime     GameStatus = "outoftime"
	Cheat         GameStatus = "cheat"
	NoStart       GameStatus = "noStart"
	UnknownFinish GameStatus = "unknownFinish"
	VariantEnd    GameStatus = "variantEnd"
)

// GameUser represents a Lichess game user.
type GameUser struct {
	User        *LightUser        `json:"user,omitempty"`
	Rating      *int              `json:"rating,omitempty"`
	RatingDiff  *int              `json:"ratingDiff,omitempty"`
	Name        *string           `json:"name,omitempty"`
	Provisional *bool             `json:"provisional,omitempty"`
	AILevel     *int              `json:"aiLevel,omitempty"`
	Analysis    *GameUserAnalysis `json:"analysis,omitempty"`
	Team        *string           `json:"team,omitempty"`
}

// GameUserAnalysis represents a Lichess game user analysis.
type GameUserAnalysis struct {
	Inaccuracy int `json:"inaccuracy,omitempty"`
	Mistake    int `json:"mistake,omitempty"`
	Blunder    int `json:"blunder,omitempty"`
	Acpl       int `json:"acpl,omitempty"`
}

// GameOpening represents a Lichess opening analysis.
type GameOpening struct {
	Eco  *string `json:"eco,omitempty"`
	Name *string `json:"name,omitempty"`
	Ply  *int    `json:"ply,omitempty"`
}

// GameAnalysis represents a Lichess game user analysis.
type GameAnalysis struct {
	Eval      int                   `json:"eval,omitempty"`
	Best      *string               `json:"best,omitempty"`
	Variation *string               `json:"variation,omitempty"`
	Judgement *GameAnalysisJudgment `json:"judgement,omitempty"`
}

// GameAnalysisJudgment represents a Lichess game user analysis judgment.
type GameAnalysisJudgment struct {
	Name    *string `json:"name,omitempty"`
	Comment *string `json:"comment,omitempty"`
}

// GameClock represents a Lichess game clock.
type GameClock struct {
	Initial   *int `json:"initial,omitempty"`
	Increment *int `json:"increment,omitempty"`
	TotalTime *int `json:"totalTime,omitempty"`
	Limit     *int `json:"limit,omitempty"`
}
