package user

import (
	"context"

	"github.com/segmentio/ksuid"
	"github.com/tidwall/buntdb"

	jsoniter "github.com/json-iterator/go"
)

type JSONDB struct {
	db *buntdb.DB
}

func NewJSONDB(path string) (*JSONDB, error) {
	if path == "" {
		path = ":memory:"
	}
	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}
	return &JSONDB{db: db}, nil
}

func (j *JSONDB) Write(ctx context.Context, user Public) (Public, error) {
	if user.ID == "" {
		user.ID = ksuid.New().String()
	}

	data, err := jsoniter.MarshalToString(user)
	if err != nil {
		return user, err
	}

	err = j.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(user.ID, data, nil)
		return err
	})
	return user, err
}

func (j *JSONDB) Read(ctx context.Context, ID string) (Public, error) {
	var user Public
	err := j.db.View(func(tx *buntdb.Tx) error {
		data, err := tx.Get(ID)
		if err != nil {
			return err
		}
		return jsoniter.UnmarshalFromString(data, &user)
	})
	return user, err
}
