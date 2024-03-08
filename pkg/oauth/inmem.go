package oauth

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// InMemoryTokenStore is an insecure store for temporary local development.
type InMemoryTokenStore struct {
	toks map[string]*oauth2.Token
}

func (s *InMemoryTokenStore) GetToken(ctx *gin.Context, ID string) (*oauth2.Token, error) {
	if tok, ok := s.toks[ID]; ok {
		return tok, nil
	} else {
		return nil, errors.New("token not found")
	}
}

func (s *InMemoryTokenStore) SaveToken(ctx *gin.Context, ID string, token *oauth2.Token) error {
	if s.toks == nil {
		s.toks = make(map[string]*oauth2.Token)
	}
	s.toks[ID] = token
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
