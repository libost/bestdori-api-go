package bestdoriapigo

import (
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
