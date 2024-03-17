package user

type ErrNotFound struct {
	UserID string `json:"user_id"`
}

func (e ErrNotFound) Error() string {
	return "user not found with id: " + e.UserID
}
