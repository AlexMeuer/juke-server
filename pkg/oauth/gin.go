package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/alexmeuer/juke/pkg/spotify"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const (
	COOKIE_VERIFIER     = "oauth_verifier"
	USER_ID_SESSION_KEY = "userID"
)

type TokenSaver interface {
	SaveToken(ctx *gin.Context, ID string, token *oauth2.Token) error
}

type StateGenerator interface {
	GenerateState(ctx *gin.Context) string
}

type StateVerifier interface {
	VerifyState(ctx *gin.Context, state string) error
}

// RegisterRoutes registers the routes for the oauth flow.
// It will register the following routes:
// - GET /login: Redirects the user to the oauth provider.
// - GET /callback: Handles the callback from the oauth provider.
// TODO: Support remembering the redirect URL that sent the user to /login and send them back there after the oauth flow is complete.
func RegisterRoutes(r gin.IRouter, config *Config, tokenSaver TokenSaver, stateGenerator StateGenerator, stateVerifier StateVerifier) error {
	if tokenSaver == nil {
		return errors.New("tokenSetter is required")
	}
	if stateGenerator == nil {
		return errors.New("stateGenerator is required")
	}

	// We all know that `GET` handlers should be idempotent and not mutate state,
	// however, we need to use `GET` here for the oauth flow to work.

	r.GET("/login", func(ctx *gin.Context) {
		state := url.QueryEscape(stateGenerator.GenerateState(ctx))
		url, verifier := config.GenerateURLAndVerifier(state)
		ctx.SetCookie(COOKIE_VERIFIER, verifier, 600, "/spotify", "", true, true)
		ctx.Redirect(http.StatusFound, url)
	})

	r.GET("/callback", func(ctx *gin.Context) {
		handleOAuthCallback(ctx, config, tokenSaver, stateVerifier)
	})

	return nil
}

func handleOAuthCallback(ctx *gin.Context, config *Config, tokenSaver TokenSaver, stateVerifier StateVerifier) {
	if err := ctx.Query("error"); err != "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	// Fail if the query parameters contain 'error' or if 'code' or 'state' are missing.
	code, state, err := validateOAuthCallbackParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Fail if the state does not match.
	if err := stateVerifier.VerifyState(ctx, state); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	verifier, err := ctx.Cookie(COOKIE_VERIFIER)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "verifier cookie is required"})
		return
	}

	tok, err := config.Exchange(ctx, code, state, verifier)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := fetchUserInfoAndSave(ctx, config, tokenSaver, tok); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: Redirect to the URL that sent the user to /login.
	ctx.Header("Location", "/foo/bar")
	ctx.Status(http.StatusFound)
}

func validateOAuthCallbackParams(ctx *gin.Context) (string, string, error) {
	code := ctx.Query("code")
	if code == "" {
		return "", "", errors.New("code is required")
	}

	state := ctx.Query("state")
	if state == "" {
		return "", "", errors.New("state is required")
	}

	state, err := url.QueryUnescape(state)
	if err != nil {
		return "", "", fmt.Errorf("failed to unescape state parameter: %w", err)
	}

	return code, state, nil
}

func exchangeToken(ctx *gin.Context, config *Config, code, state, verifier string) (*oauth2.Token, error) {
	verifier, err := ctx.Cookie(COOKIE_VERIFIER)
	if err != nil {
		return nil, errors.New("cookie is required")
	}

	tok, err := config.Exchange(ctx, code, state, verifier)
	if err != nil {
		return nil, err
	}

	return tok, nil
}

// fetchUserInfoAndSave fetches the user's info from the spotify API, it then saves the user's ID in the session and saves the token.
func fetchUserInfoAndSave(ctx *gin.Context, config *Config, tokenSaver TokenSaver, tok *oauth2.Token) error {
	client := spotify.New(config.Client(ctx, tok))
	me, err := client.Me(ctx)
	if err != nil {
		return fmt.Errorf("login succeeded but failed to get user info: %w", err)
	}

	s := sessions.Default(ctx)
	s.Set(USER_ID_SESSION_KEY, me.ID)
	if err := s.Save(); err != nil {
		return fmt.Errorf("login succeeded but failed to save session: %w", err)
	}
	fmt.Printf("Saved session: %+v\n", s)

	if err := tokenSaver.SaveToken(ctx, me.ID, tok); err != nil {
		return fmt.Errorf("login succeeded but failed to save token: %w", err)
	}

	return nil
}
