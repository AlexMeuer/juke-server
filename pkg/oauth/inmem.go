package oauth

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// InMemoryTokenStore is an insecure store
// that stores a single token in memory.
type InMemoryTokenStore struct {
	tok *oauth2.Token
}

func (s *InMemoryTokenStore) GetToken(ctx *gin.Context) (*oauth2.Token, error) {
	if s.tok == nil {
		return nil, errors.New("token not stored")
	}
	return s.tok, nil
}

func (s *InMemoryTokenStore) SetToken(ctx *gin.Context, token *oauth2.Token) error {
	s.tok = token
	return nil
}

type InMemoryStateManager struct {
	state string
}

func (m *InMemoryStateManager) GenerateState(ctx *gin.Context) string {
	m.state = time.Now().Format(time.Kitchen)
	return m.state
}

func (m *InMemoryStateManager) VerifyState(ctx *gin.Context, state string) error {
	if state != m.state {
		return errors.New("state mismatch")
	}
	return nil
}
