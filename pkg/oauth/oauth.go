package oauth

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2"
)

func Foo(ctx context.Context) {
	conf := &oauth2.Config{
		ClientID:     "836b26ceac5c4e9995451e3b55c6dc2b",
		ClientSecret: "7e49738fc2d2400fabb93fad0faa5a4d",
		Scopes:       []string{"user-read-playback-state", "user-modify-playback-state", "user-read-currently-playing"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		RedirectURL: "http://localhost:8080/callback",
	}

	verifier := oauth2.GenerateVerifier()

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	fmt.Printf("Visit the URL for the auth dialog: %v\n\nPaste the code here: ", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nToken: %+v\n", tok)
}
