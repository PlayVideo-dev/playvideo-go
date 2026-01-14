package playvideo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   map[string]interface{}
		checkError func(t *testing.T, err error)
	}{
		{
			name:       "401 returns AuthenticationError",
			statusCode: 401,
			response:   map[string]interface{}{"error": "Unauthorized", "message": "Invalid API key"},
			checkError: func(t *testing.T, err error) {
				if !IsAuthenticationError(err) {
					t.Errorf("expected AuthenticationError, got %T", err)
				}
			},
		},
		{
			name:       "403 returns AuthorizationError",
			statusCode: 403,
			response:   map[string]interface{}{"error": "Forbidden", "message": "Insufficient permissions"},
			checkError: func(t *testing.T, err error) {
				if !IsAuthorizationError(err) {
					t.Errorf("expected AuthorizationError, got %T", err)
				}
			},
		},
		{
			name:       "404 returns NotFoundError",
			statusCode: 404,
			response:   map[string]interface{}{"error": "Not Found", "message": "Resource not found"},
			checkError: func(t *testing.T, err error) {
				if !IsNotFoundError(err) {
					t.Errorf("expected NotFoundError, got %T", err)
				}
			},
		},
		{
			name:       "400 returns ValidationError",
			statusCode: 400,
			response:   map[string]interface{}{"error": "Bad Request", "message": "Invalid input"},
			checkError: func(t *testing.T, err error) {
				if !IsValidationError(err) {
					t.Errorf("expected ValidationError, got %T", err)
				}
			},
		},
		{
			name:       "422 returns ValidationError",
			statusCode: 422,
			response:   map[string]interface{}{"error": "Unprocessable Entity", "message": "Invalid input"},
			checkError: func(t *testing.T, err error) {
				if !IsValidationError(err) {
					t.Errorf("expected ValidationError, got %T", err)
				}
			},
		},
		{
			name:       "409 returns ConflictError",
			statusCode: 409,
			response:   map[string]interface{}{"error": "Conflict", "message": "Resource conflict"},
			checkError: func(t *testing.T, err error) {
				if !IsConflictError(err) {
					t.Errorf("expected ConflictError, got %T", err)
				}
			},
		},
		{
			name:       "429 returns RateLimitError",
			statusCode: 429,
			response:   map[string]interface{}{"error": "Too Many Requests", "message": "Rate limit exceeded"},
			checkError: func(t *testing.T, err error) {
				if !IsRateLimitError(err) {
					t.Errorf("expected RateLimitError, got %T", err)
				}
			},
		},
		{
			name:       "500 returns ServerError",
			statusCode: 500,
			response:   map[string]interface{}{"error": "Internal Server Error", "message": "Something went wrong"},
			checkError: func(t *testing.T, err error) {
				if !IsServerError(err) {
					t.Errorf("expected ServerError, got %T", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				writeJSON(w, tt.response)
			}))
			defer server.Close()

			client := NewClient("play_test_xxx", WithBaseURL(server.URL))
			_, err := client.Collections.List(context.Background())

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			tt.checkError(t, err)
		})
	}
}

func TestErrorMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", "req-abc-123")
		w.WriteHeader(404)
		writeJSON(w, map[string]interface{}{
			"error":   "Not Found",
			"message": "Video not found",
		})
	}))
	defer server.Close()

	client := NewClient("play_test_xxx", WithBaseURL(server.URL))
	_, err := client.Videos.Get(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("expected error")
	}

	notFoundErr, ok := err.(*NotFoundError)
	if !ok {
		t.Fatalf("expected NotFoundError, got %T", err)
	}

	if notFoundErr.RequestID != "req-abc-123" {
		t.Errorf("expected requestId 'req-abc-123', got %s", notFoundErr.RequestID)
	}

	if notFoundErr.Message != "Video not found" {
		t.Errorf("expected message 'Video not found', got %s", notFoundErr.Message)
	}
}
