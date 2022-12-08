package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type githubUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (g githubUser) GetID() string {
	return fmt.Sprintf("%d", g.ID)
}

func (g githubUser) GetEmail() string {
	return g.Email
}

func NewGithubConfig(clientID, clientSecret, host string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
		RedirectURL:  fmt.Sprintf("%s/register/auth/github/callback", host),
	}
}

func GetGithubUser(client *http.Client) GetUserFunc {
	return func(ctx context.Context, token *oauth2.Token) (OAuthUser, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
		if err != nil {
			return githubUser{}, err
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return githubUser{}, fmt.Errorf("failed to get user info: %w", err)
		}
		defer resp.Body.Close()

		var gu githubUser
		if err := json.NewDecoder(resp.Body).Decode(&gu); err != nil || gu == (githubUser{}) {
			return githubUser{}, fmt.Errorf("failed to decode user info: %w", err)
		}

		return gu, nil
	}
}
