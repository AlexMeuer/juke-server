package room

import "context"

type Writer interface {
	Write(ctx context.Context, room Room) (Room, error)
}

type Reader interface {
	Read(ctx context.Context, ID string) (Room, error)
}
