package traefik_plugin_return_response

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

// The configuration for a response rewrite.
type Response struct {
	// The HTTP method to match.
	Method string `json:"method,omitempty"`
	// The URL path to match. Regular expressions are supported.
	UrlMatch string `json:"url_match,omitempty"`
	// The HTTP status code to return.
	StatusCode int `json:"status_code,omitempty"`
}

// The configuration for the plugin.
type Config struct {
	Response Response `json:"response,omitempty"`
}

// Creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type response struct {
	method     string
	urlMatch   *regexp.Regexp
	statusCode int
}

type returnResponse struct {
	name     string
	next     http.Handler
	response response
}

// Creates and returns a new plugin instance.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	configResponse := config.Response

	regex, err := regexp.Compile(configResponse.UrlMatch)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex %q: %w", configResponse.UrlMatch, err)
	}
	r := response{method: configResponse.Method, urlMatch: regex, statusCode: configResponse.StatusCode}

	return &returnResponse{
		name:     name,
		next:     next,
		response: r,
	}, nil
}

// Serves the HTTP request, returning the configured response if the request matches a configured method and path.
func (r *returnResponse) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if r.response.method != "" && r.response.method != req.Method {
		r.next.ServeHTTP(rw, req)
		return
	}

	if r.response.urlMatch != nil && !r.response.urlMatch.MatchString(req.URL.String()) {
		r.next.ServeHTTP(rw, req)
		return
	}

	rw.WriteHeader(r.response.statusCode)
}
