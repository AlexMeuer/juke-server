package server

import (
	"fmt"
	"net/http"

	"github.com/alexmeuer/juke/pkg/oauth"
	"github.com/alexmeuer/juke/pkg/openapi"
	"github.com/gin-gonic/gin"
)

func Serve(host string, port uint16) error {
	r := openapi.NewRouter(openapi.ApiHandleFunctions{
		RoomsAPI: &RoomsAPI{},
	})

	tokenStore := &oauth.InMemoryTokenStore{}
	stateManager := &oauth.InMemoryStateManager{}

	var clientID, clientSecret string
	fmt.Print("Enter client ID: ")
	if _, err := fmt.Scan(&clientID); err != nil {
		return err
	}
	fmt.Print("Enter client secret: ")
	if _, err := fmt.Scan(&clientSecret); err != nil {
		return err
	}

	cfg := oauth.NewSpotify(clientID, clientSecret, fmt.Sprintf("http://%s:%d/spotify/callback", host, port))
	oauth.RegisterRoutes(r.Group("/spotify"), cfg, tokenStore, stateManager, stateManager)

	r.Group("/foo").Use(oauth.NewTokenClientMiddleware(cfg, tokenStore)).GET("/bar", func(ctx *gin.Context) {
		client, err := oauth.GetClient(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp, err := client.Get("https://api.spotify.com/v1/me")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		ctx.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	return r.Run(fmt.Sprintf("%s:%d", host, port))
}
