package oauth

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const (
	COOKIE_VERIFIER      = "oauth_verifier"
	TOKEN_CLIENT_CTX_KEY = "token_client"
)

type TokenGetter interface {
	GetToken(ctx *gin.Context) (*oauth2.Token, error)
}

type TokenSetter interface {
	SetToken(ctx *gin.Context, token *oauth2.Token) error
}

type StateGenerator interface {
	GenerateState(ctx *gin.Context) string
}

type StateVerifier interface {
	VerifyState(ctx *gin.Context, state string) error
}

func GetClient(ctx *gin.Context) (*http.Client, error) {
	value, ok := ctx.Get(TOKEN_CLIENT_CTX_KEY)
	if !ok {
		return nil, errors.New("token client not found")
	}
	client, ok := value.(*http.Client)
	if !ok {
		return nil, errors.New("token client is not of type *http.Client")
	}
	return client, nil
}

// NewTokenClientMiddleware returns a middleware that sets the token client in the context.
// You can use the constant `TOKEN_CLIENT_CTX_KEY` to retrieve the `*http.client` from the context.
func NewTokenClientMiddleware(config *Config, tokenGetter TokenGetter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tok, err := tokenGetter.GetToken(ctx)
		if err != nil {
			// TODO: redirect to auth flow if we can.
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set(TOKEN_CLIENT_CTX_KEY, config.Client(ctx, tok))
		ctx.Next()
	}
}

// RegisterRoutes registers the routes for the oauth flow.
// It will register the following routes:
// - GET /auth: Redirects the user to the oauth provider.
// - GET /callback: Handles the callback from the oauth provider.
func RegisterRoutes(r gin.IRouter, config *Config, tokenSetter TokenSetter, stateGenerator StateGenerator, stateVerifier StateVerifier) error {
	if tokenSetter == nil {
		return errors.New("tokenSetter is required")
	}
	if stateGenerator == nil {
		return errors.New("stateGenerator is required")
	}

	// We all know that `GET` handlers should be idempotent and not mutate state,
	// however, we need to use `GET` here for the oauth flow to work.

	r.GET("/auth", func(ctx *gin.Context) {
		state := url.QueryEscape(stateGenerator.GenerateState(ctx))
		url, verifier := config.GenerateURLAndVerifier(state)
		ctx.SetCookie(COOKIE_VERIFIER, verifier, 600, "/spotify", "", true, true)
		ctx.Redirect(http.StatusFound, url)
	})

	r.GET("/callback", func(ctx *gin.Context) {
		if err := ctx.Query("error"); err != "" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		}

		code := ctx.Query("code")
		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
			return
		}
		state := ctx.Query("state")

		if state == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "state is required"})
			return
		} else {
			parsedState, err := url.QueryUnescape(state)
			if err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
				return
			}
			if err := stateVerifier.VerifyState(ctx, parsedState); err != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
				return
			}
		}

		verifier, err := ctx.Cookie(COOKIE_VERIFIER)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "cookie is required"})
			return
		}

		tok, err := config.Exchange(ctx, code, state, verifier)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		tokenSetter.SetToken(ctx, tok)

		ctx.JSON(http.StatusOK, gin.H{})
	})

	return nil
}
