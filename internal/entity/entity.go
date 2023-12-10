package entity

type Request struct {
	UserID int    `json:"user_id,omitempty"`
	UUID   string `json:"uuid,omitempty"`
	Alias  string `json:"short_url,omitempty"`
	URL    string `json:"original_url,omitempty" validate:"required,url"`
}
