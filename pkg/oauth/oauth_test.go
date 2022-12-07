package oauth_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/arthureichelberger/authpulse/pkg/oauth"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestGetConnectionURL(t *testing.T) {
	type testCase struct {
		name              string
		config            oauth2.Config
		externalID        string
		shouldExpectError bool
		expectedURL       string
	}

	testCases := []testCase{
		{
			name: "invalid auth url",
			config: oauth2.Config{
				Endpoint: oauth2.Endpoint{
					AuthURL: "test",
				},
			},
			shouldExpectError: true,
		},
		{
			name: "valid configuration",
			config: oauth2.Config{
				ClientID:     "test",
				ClientSecret: "test",
				Scopes:       []string{"test"},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://test.com",
					TokenURL: "https://test.com",
				},
				RedirectURL: "https://test.com",
			},
			externalID:        uuid.New().String(),
			shouldExpectError: false,
			expectedURL:       "https://test.com?client_id=test&redirect_uri=https%3A%2F%2Ftest.com&response_type=code&scope=test&state=",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oa := oauth.NewOAuthenticator(&tc.config)
			connectionURL, err := oa.GetConnectionURL(tc.externalID)
			if tc.shouldExpectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			url, err := url.Parse(connectionURL)
			assert.NoError(t, err)
			assert.Equal(t, tc.config.ClientID, url.Query().Get("client_id"))
			assert.Equal(t, strings.Join(tc.config.Scopes, ","), url.Query().Get("scope"))
			assert.Equal(t, tc.config.RedirectURL, url.Query().Get("redirect_uri"))
			assert.Equal(t, "code", url.Query().Get("response_type"))
			assert.Equal(t, tc.externalID, url.Query().Get("state"))
		})
	}
}

func TestExchange(t *testing.T) {
	type testCase struct {
		name              string
		client            *http.Client
		config            oauth2.Config
		code              string
		shouldExpectError bool
	}

	testCases := []testCase{
		{
			name:   "invalid code",
			client: intercept(http.StatusBadRequest, json.RawMessage(`{"access_token": "test"}`)),
			config: oauth2.Config{
				ClientID:     "test",
				ClientSecret: "test",
				Scopes:       []string{"test"},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://test.com",
					TokenURL: "https://test.com",
				},
				RedirectURL: "https://test.com",
			},
			shouldExpectError: true,
		},
		{
			name:   "valid code",
			client: intercept(http.StatusOK, json.RawMessage(`{"access_token": "test"}`)),
			config: oauth2.Config{
				ClientID:     "test",
				ClientSecret: "test",
				Scopes:       []string{"test"},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://test.com",
					TokenURL: "https://test.com",
				},
				RedirectURL: "https://test.com",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, tc.client)
			oa := oauth.NewOAuthenticator(&tc.config)
			token, err := oa.Exchange(ctx, tc.code)
			if tc.shouldExpectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, token)
			assert.Equal(t, "test", token.AccessToken)
		})
	}
}
