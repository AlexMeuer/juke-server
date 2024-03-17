package room

type ErrNotFound struct {
	RoomID string `json:"room_id"`
}

func (e ErrNotFound) Error() string {
	return "room not found with ID: " + e.RoomID
}
