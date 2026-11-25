package bestdoriapigo

import (
	"context"
	"net/http"

	"github.com/libost/bestdori-api-go/api"
)

// API is the main Bestdori API client.
//
// The client type is kept as an alias so code importing the module root and
// code importing the api subpackage can use the same API.
type API = api.Api

// NewClient creates a Bestdori API client.
func NewClient(httpClient *http.Client) *API {
	return api.NewApi(httpClient)
}

// NewClientWithLogin creates a Bestdori API client and logs in with the provided username and password.
//
// If the username or password is empty, it returns a api.ErrNoUsernameOrPassword error. If the login request fails, it returns an error with the status code and response body.
//
// At the moment, login won't give you any additional access to the API, but it may be useful for future features.
func NewClientWithLogin(ctx context.Context, httpClient *http.Client, username, password string) (*API, error) {
	return api.NewApiWithLogin(ctx, httpClient, username, password)
}
