package oauth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/arthureichelberger/authpulse/pkg/oauth"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func TestNewGithubConfig(t *testing.T) {
	cfg := oauth.NewGithubConfig("test", "test", "test")
	assert.Equal(t, "test", cfg.ClientID)
	assert.Equal(t, "test", cfg.ClientSecret)
	assert.Equal(t, "test/register/auth/github/callback", cfg.RedirectURL)
	assert.Equal(t, github.Endpoint, cfg.Endpoint)
	assert.Equal(t, []string{"user:email", "read:user"}, cfg.Scopes)
}

func TestGetGithubUser(t *testing.T) {
	type testCase struct {
		name              string
		ctx               context.Context
		interceptor       *http.Client
		shouldExpectError bool
	}

	testCases := []testCase{
		{
			name:              "no valid context",
			ctx:               nil,
			shouldExpectError: true,
		},
		{
			name:              "Wrong response status code from interceptor",
			ctx:               context.Background(),
			interceptor:       intercept(http.StatusBadRequest, json.RawMessage(`{}`)),
			shouldExpectError: true,
		},
		{
			name:              "Wrong response payload from interceptor",
			ctx:               context.Background(),
			interceptor:       intercept(http.StatusOK, json.RawMessage(`{"0": 1}`)),
			shouldExpectError: true,
		},
		{
			name:        "Valid response from interceptor",
			ctx:         context.Background(),
			interceptor: intercept(http.StatusOK, json.RawMessage(`{"id": 1, "email": "john.doe@gmail.com"}`)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := oauth.GetGithubUser(tc.interceptor)
			user, err := f(tc.ctx, &oauth2.Token{})
			if tc.shouldExpectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Implements(t, (*oauth.OAuthUser)(nil), user)
			assert.Equal(t, "1", user.GetID())
			assert.Equal(t, "john.doe@gmail.com", user.GetEmail())
		})
	}
}
