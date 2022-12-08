package oauth

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

type OAuthUser interface {
	GetID() string
	GetEmail() string
}

type OAuthenticator struct {
	cfg *oauth2.Config
}

func NewOAuthenticator(cfg *oauth2.Config) OAuthenticator {
	return OAuthenticator{cfg: cfg}
}

func (oa OAuthenticator) GetConnectionURL(externalID string) (string, error) {
	authURL, err := url.ParseRequestURI(oa.cfg.Endpoint.AuthURL)
	if err != nil {
		return "", fmt.Errorf("invalid auth url: %w", err)
	}

	parameters := url.Values{
		"client_id":     {oa.cfg.ClientID},
		"scope":         {strings.Join(oa.cfg.Scopes, " ")},
		"redirect_uri":  {oa.cfg.RedirectURL},
		"response_type": {"code"},
		"state":         {externalID},
	}
	authURL.RawQuery = parameters.Encode()

	return authURL.String(), nil
}

func (oa OAuthenticator) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := oa.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return token, nil
}

type GetUserFunc func(ctx context.Context, token *oauth2.Token) (OAuthUser, error)
