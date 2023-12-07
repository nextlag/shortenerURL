package entity

type Request struct {
	UUID  string `json:"uuid"`
	Alias string `json:"alias,omitempty"`
	URL   string `json:"url" validate:"required,url"`
}
