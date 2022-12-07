package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (g googleUser) GetID() string {
	return g.ID
}

func (g googleUser) GetEmail() string {
	return g.Email
}

func NewGoogleConfig(clientID, clientSecret, host string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
		RedirectURL:  fmt.Sprintf("%s/register/auth/google/callback", host),
	}
}

func GetGoogleUser(client *http.Client) GetUserFunc {
	return func(ctx context.Context, token *oauth2.Token) (OAuthUser, error) {
		endpoint, _ := url.Parse("https://www.googleapis.com/oauth2/v2/userinfo")
		params := url.Values{}
		params.Add("access_token", token.AccessToken)
		endpoint.RawQuery = params.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
		if err != nil {
			return googleUser{}, err
		}

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return googleUser{}, fmt.Errorf("failed to get user info: %w", err)
		}
		defer resp.Body.Close()

		var gu googleUser
		if err := json.NewDecoder(resp.Body).Decode(&gu); err != nil || gu == (googleUser{}) {
			return googleUser{}, fmt.Errorf("failed to decode user info: %w", err)
		}

		return gu, nil
	}
}
