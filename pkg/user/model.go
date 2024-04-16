package user

type Public struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type Private struct {
	Public
	Email string `json:"email"`
}
