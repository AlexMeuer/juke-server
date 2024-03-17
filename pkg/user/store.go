package user

import "context"

type Writer interface {
	Write(ctx context.Context, user Public) (Public, error)
}

type Reader interface {
	Read(ctx context.Context, ID string) (Public, error)
}
