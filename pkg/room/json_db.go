package room

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

func (j *JSONDB) Write(ctx context.Context, room Room) (Room, error) {
	if room.ID == "" {
		room.ID = ksuid.New().String()
	}

	data, err := jsoniter.MarshalToString(room)
	if err != nil {
		return room, err
	}

	err = j.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(room.ID, data, nil)
		return err
	})
	return room, err
}

func (j *JSONDB) Read(ctx context.Context, ID string) (Room, error) {
	var room Room
	err := j.db.View(func(tx *buntdb.Tx) error {
		data, err := tx.Get(ID)
		if err != nil {
			return err
		}
		return jsoniter.UnmarshalFromString(data, &room)
	})
	return room, err
}
