package oauth

import (
	"net/http"

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

// NewTokenClientMiddleware returns a middleware that sets the token client in the context.
// You can use the constant TOKEN_CLIENT_CTX_KEY to retrieve the http.client from the context.
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

func RegisterRoutes(r gin.IRouter, config *Config, tokenSetter TokenSetter) {
	r.POST("/spotify/auth/:userId", func(ctx *gin.Context) {
		url, verifier := config.GenerateURLAndVerifier("state")
		ctx.SetCookie(COOKIE_VERIFIER, verifier, 600, "/spotify", "", true, true)
		ctx.Redirect(http.StatusFound, url)
	})

	r.POST("/spotify/callback", func(ctx *gin.Context) {
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

		tokenSetter.SetToken(ctx, state, tok)

		ctx.JSON(http.StatusOK, gin.H{})
	})
}
