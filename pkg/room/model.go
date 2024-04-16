package room

type Room struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	OwnerID string   `json:"owner"`
	Sources []string `json:"sources"`
}
