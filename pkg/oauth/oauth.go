package oauth

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type Config struct {
	cfg oauth2.Config
}

func NewSpotify(clientID, clientSecret, redirectURL string) *Config {
	return &Config{
		cfg: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes: []string{
				"user-read-playback-state",
				"user-modify-playback-state",
				"user-read-currently-playing",
			},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.spotify.com/authorize",
				TokenURL: "https://accounts.spotify.com/api/token",
			},
			RedirectURL: redirectURL,
		},
	}
}

func (c *Config) Client(ctx context.Context, tok *oauth2.Token) *http.Client {
	return c.cfg.Client(ctx, tok)
}

func (c *Config) GenerateURLAndVerifier(state string) (string, string) {
	verifier := oauth2.GenerateVerifier()
	return c.cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier)), verifier
}

func (c *Config) Exchange(ctx context.Context, code, state, verifier string) (*oauth2.Token, error) {
	tok, err := c.cfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func Debug(ctx context.Context) error {
	var clientID, clientSecret string
	fmt.Print("Enter client ID: ")
	if _, err := fmt.Scan(&clientID); err != nil {
		return err
	}
	fmt.Print("\nEnter client secret: ")
	if _, err := fmt.Scan(&clientSecret); err != nil {
		return err
	}
	conf := NewSpotify(clientID, clientSecret, "http://localhost:8080/callback")

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url, verifier := conf.GenerateURLAndVerifier("state")
	fmt.Printf("\nVisit the URL for the auth dialog: %s\n\nPaste the code here: ", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return err
	}
	tok, err := conf.Exchange(ctx, code, "state", verifier)
	if err != nil {
		return err
	}

	fmt.Printf("\nToken: %+v\n", tok)
	return nil
}
