package oauth

import (
	"context"
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/buntdb"
	"golang.org/x/oauth2"
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

func (j *JSONDB) NewJSONDB(path string) (*JSONDB, error) {
	if path == "" {
		path = ":memory:"
	}
	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}
	return &JSONDB{db: db}, nil
}

func (j *JSONDB) WriteToken(ctx context.Context, ID string, token *oauth2.Token) error {
	if ID == "" {
		return errors.New("ID is required to write token")
	}

	if token == nil {
		return errors.New("token is required to write token")
	}

	data, err := jsoniter.MarshalToString(token)
	if err != nil {
		return fmt.Errorf("error marshalling token: %w", err)
	}

	err = j.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(ID, data, nil)
		return err
	})
	if err != nil {
		return fmt.Errorf("error writing token: %w", err)
	}
	return nil
}

func (j *JSONDB) ReadToken(ctx context.Context, ID string) (*oauth2.Token, error) {
	var token oauth2.Token
	err := j.db.View(func(tx *buntdb.Tx) error {
		data, err := tx.Get(ID)
		if err != nil {
			return err
		}
		return jsoniter.UnmarshalFromString(data, &token)
	})
	if err != nil {
		return nil, fmt.Errorf("error reading token: %w", err)
	}
	return &token, nil
}

func (j *JSONDB) WriteState(ctx context.Context, ID string, state string) error {
	if ID == "" {
		return errors.New("ID is required to write state")
	}

	if state == "" {
		return errors.New("state is required to write state")
	}

	err := j.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(ID, state, nil)
		return err
	})
	if err != nil {
		return fmt.Errorf("error writing state: %w", err)
	}
	return nil
}

func (j *JSONDB) VerifyState(ctx context.Context, ID, state string) error {
	var actualState string
	err := j.db.View(func(tx *buntdb.Tx) error {
		data, err := tx.Get(ID)
		if err != nil {
			return err
		}
		actualState = data
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading state: %w", err)
	}
	if state != actualState {
		return errors.New("state mismatch")
	}
	return nil
}
