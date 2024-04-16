package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexmeuer/juke/pkg/oauth"
	"github.com/alexmeuer/juke/pkg/spotify"
	"github.com/gin-contrib/cors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

func Serve(host string, port uint16) error {
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options(sessions.Options{
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/",
	})

	cfg := oauth.NewSpotify(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"), fmt.Sprintf("https://%s:%d/spotify/callback", host, port))

	tokenStore, err := oauth.NewJSONDB("api/tokens.db")
	if err != nil {
		return err
	}

	r := gin.Default()

	r.Use(CORSMiddleware()).
		Use(sessions.Sessions("juke-session", store)).
		Use(func(ctx *gin.Context) {
			ctx.Set("spotifyConfig", cfg)
			ctx.Next()
		})

	if err = oauth.RegisterRoutes(r.Group("/spotify"), cfg, tokenStore, &oauth.GinStateGenerator{}, tokenStore, tokenStore); err != nil {
		return err
	}

	foo := r.Group("/foo").Use(requireSessionMiddleware).Use(spotifyClientMiddleware(tokenStore))

	foo.GET("/bar", func(ctx *gin.Context) {
		client := ctx.MustGet("spotifyClient").(*spotify.Client)
		me, err := client.Me(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, me)
	})

	foo.GET("/devices", func(ctx *gin.Context) {
		client := ctx.MustGet("spotifyClient").(*spotify.Client)
		me, err := client.AvailableDevices(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, me)
	})

	if os.Getenv("GIN_MODE") == "release" {
		return r.Run(fmt.Sprintf("%s:%d", host, port))
	}
	return r.RunTLS(fmt.Sprintf("%s:%d", host, port), "api/cert.pem", "api/key.pem")
}

func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	return cors.New(config)
}

func requireSessionMiddleware(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get(oauth.USER_ID_SESSION_KEY)
	if userID == nil {
		// No valid session found, return 401 Unauthorized status code
		ctx.Header("Location", "/spotify/login")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set("userID", userID)
	// Valid session found, continue with the request
	ctx.Next()
}

func spotifyClientMiddleware(r oauth.TokenReader) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tok, err := r.ReadToken(ctx, ctx.GetString(oauth.USER_ID_SESSION_KEY))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cfg := ctx.MustGet("spotifyConfig").(*oauth.Config)
		client := spotify.New(cfg.Client(ctx, tok))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.Set("spotifyClient", client)
	}
}
