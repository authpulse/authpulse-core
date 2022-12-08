package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/arthureichelberger/authpulse/pkg/oauth"
	"github.com/arthureichelberger/authpulse/pkg/router"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sign
		log.Debug().Msg("Shutting down")
		cancel()
	}()

	googleOAuthConfig := oauth.NewGoogleConfig(env("GOOGLE_CLIENT_ID", "undefined"), env("GOOGLE_CLIENT_SECRET", "undefined"), env("APP_URL", "http://localhost:8080"))
	googleOAuthenticator := oauth.NewOAuthenticator(googleOAuthConfig)

	githubOAuthConfig := oauth.NewGithubConfig(env("GITHUB_CLIENT_ID", "undefined"), env("GITHUB_CLIENT_SECRET", "undefined"), env("APP_URL", "http://localhost:8080"))
	githubOAuthenticator := oauth.NewOAuthenticator(githubOAuthConfig)

	r := gin.New()

	registerAuthGrp := r.Group("/register/auth")
	registerAuthGrp.GET("/google", oAuthRedirect(googleOAuthenticator))
	registerAuthGrp.GET("/google/callback", oAuthCallback(googleOAuthenticator, oauth.GetGoogleUser(http.DefaultClient)))
	registerAuthGrp.GET("/github", oAuthRedirect(githubOAuthenticator))
	registerAuthGrp.GET("/github/callback", oAuthCallback(githubOAuthenticator, oauth.GetGithubUser(http.DefaultClient)))

	log.Debug().Msg("Starting server")
	if err := router.Run(ctx, ":8080", r); err != nil {
		log.Fatal().Err(err).Msg("Server crashed")
	}
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func oAuthRedirect(authenticator oauth.OAuthenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		externalID := uuid.NewString()
		url, err := authenticator.GetConnectionURL(externalID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func oAuthCallback(authenticator oauth.OAuthenticator, getUser oauth.GetUserFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.Query("code")
		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
			return
		}

		token, err := authenticator.Exchange(ctx.Request.Context(), code)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user, err := getUser(ctx.Request.Context(), token)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"user": gin.H{"id": user.GetID(), "email": user.GetEmail()}})
	}
}
