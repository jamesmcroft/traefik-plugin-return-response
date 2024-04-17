package traefik_plugin_return_response

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	target := "https://127.0.0.1"

	testCases := []struct {
		description           string
		config                Response
		requestPath           string
		requestMethod         string
		matchConfigStatusCode bool
	}{
		{
			description:           "should return status code for a RegEx path and method",
			requestPath:           target,
			requestMethod:         "OPTIONS",
			matchConfigStatusCode: true,
			config: Response{
				Method:     "OPTIONS",
				UrlMatch:   "^https://(.+)$",
				StatusCode: 200,
			},
		},
		{
			description:           "should return 404 status code for a specific path and method",
			requestPath:           target + "/test",
			requestMethod:         "GET",
			matchConfigStatusCode: true,
			config: Response{
				Method:     "GET",
				UrlMatch:   target + "/test",
				StatusCode: 404,
			},
		},
		{
			description:           "should bypass the plugin without a URL match",
			requestPath:           target,
			requestMethod:         "POST",
			matchConfigStatusCode: false,
			config: Response{
				Method:     "POST",
				UrlMatch:   "https://example.com",
				StatusCode: 500,
			},
		},
		{
			description:           "should bypass the plugin without a method match",
			requestPath:           target,
			requestMethod:         "GET",
			matchConfigStatusCode: false,
			config: Response{
				Method:     "POST",
				UrlMatch:   target,
				StatusCode: 500,
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.description, func(t *testing.T) {
			config := &Config{Response: testCase.config}

			ctx := context.Background()
			next := func(responseWriter http.ResponseWriter, request *http.Request) {}

			returnResponse, _ := New(ctx, http.HandlerFunc(next), config, "returnResponse")

			recorder := httptest.NewRecorder()
			request, err := http.NewRequestWithContext(ctx, testCase.requestMethod, testCase.requestPath, nil)

			if err != nil {
				t.Fatalf("error creating request: %v", err)
			}

			returnResponse.ServeHTTP(recorder, request)

			if testCase.matchConfigStatusCode {
				if recorder.Code != testCase.config.StatusCode {
					t.Fatalf("expected status code %d, got %d", testCase.config.StatusCode, recorder.Code)
				}
			} else {
				if recorder.Code != 200 {
					t.Fatalf("expected status code 200, got %d", recorder.Code)
				}
			}
		})
	}
}
