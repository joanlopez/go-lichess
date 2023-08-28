package lichess

type LightUser struct {
	Id     string  `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Title  *string `json:"title,omitempty"`
	Patron bool    `json:"patron,omitempty"`
}
