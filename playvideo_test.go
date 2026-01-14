package playvideo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Helper to create a test server
func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(handler)
	client := NewClient("play_test_xxx", WithBaseURL(server.URL))
	return server, client
}

// Helper to encode JSON responses in tests
func writeJSON(w http.ResponseWriter, v interface{}) {
	_ = writeJSON(w, v)
}

func TestClientInitialization(t *testing.T) {
	t.Run("creates client with API key", func(t *testing.T) {
		client := NewClient("play_test_xxx")
		if client == nil {
			t.Fatal("expected client to be created")
		}
		if client.Collections == nil {
			t.Error("expected Collections to be initialized")
		}
		if client.Videos == nil {
			t.Error("expected Videos to be initialized")
		}
		if client.Webhooks == nil {
			t.Error("expected Webhooks to be initialized")
		}
	})

	t.Run("accepts custom options", func(t *testing.T) {
		client := NewClient("play_test_xxx",
			WithBaseURL("https://custom.api.com/v1"),
			WithTimeout(60*time.Second),
		)
		if client == nil {
			t.Fatal("expected client to be created with options")
		}
	})
}

func TestCollections(t *testing.T) {
	t.Run("list collections", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" || r.URL.Path != "/collections" {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			}
			if r.Header.Get("Authorization") != "Bearer play_test_xxx" {
				t.Error("expected Authorization header")
			}

			writeJSON(w, map[string]interface{}{
				"collections": []map[string]interface{}{
					{
						"id":          "col1",
						"name":        "Test",
						"slug":        "test",
						"videoCount":  5,
						"storageUsed": 1024,
						"createdAt":   "2024-01-01T00:00:00Z",
						"updatedAt":   "2024-01-01T00:00:00Z",
					},
				},
			})
		})
		defer server.Close()

		result, err := client.Collections.List(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Collections) != 1 {
			t.Errorf("expected 1 collection, got %d", len(result.Collections))
		}
		if result.Collections[0].Slug != "test" {
			t.Errorf("expected slug 'test', got %s", result.Collections[0].Slug)
		}
	})

	t.Run("create collection", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" || r.URL.Path != "/collections" {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			}

			writeJSON(w, map[string]interface{}{
				"id":        "col1",
				"name":      "My Videos",
				"slug":      "my-videos",
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-01T00:00:00Z",
			})
		})
		defer server.Close()

		desc := "Test description"
		result, err := client.Collections.Create(context.Background(), &CreateCollectionParams{
			Name:        "My Videos",
			Description: &desc,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Slug != "my-videos" {
			t.Errorf("expected slug 'my-videos', got %s", result.Slug)
		}
	})

	t.Run("get collection", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" || r.URL.Path != "/collections/test" {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			}

			writeJSON(w, map[string]interface{}{
				"id":        "col1",
				"name":      "Test",
				"slug":      "test",
				"videos":    []interface{}{},
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-01T00:00:00Z",
			})
		})
		defer server.Close()

		result, err := client.Collections.Get(context.Background(), "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Slug != "test" {
			t.Errorf("expected slug 'test', got %s", result.Slug)
		}
	})

	t.Run("delete collection", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" || r.URL.Path != "/collections/test" {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			}

			writeJSON(w, map[string]string{"message": "Deleted"})
		})
		defer server.Close()

		err := client.Collections.Delete(context.Background(), "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestVideos(t *testing.T) {
	t.Run("list videos", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" || r.URL.Path != "/videos" {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			}

			writeJSON(w, map[string]interface{}{
				"videos": []map[string]interface{}{
					{
						"id":        "vid1",
						"filename":  "test.mp4",
						"status":    "COMPLETED",
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					},
				},
			})
		})
		defer server.Close()

		result, err := client.Videos.List(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Videos) != 1 {
			t.Errorf("expected 1 video, got %d", len(result.Videos))
		}
	})

	t.Run("list videos with params", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("collection") != "my-collection" {
				t.Error("expected collection param")
			}
			if r.URL.Query().Get("status") != "COMPLETED" {
				t.Error("expected status param")
			}

			writeJSON(w, map[string]interface{}{"videos": []interface{}{}})
		})
		defer server.Close()

		_, err := client.Videos.List(context.Background(), &VideoListParams{
			Collection: "my-collection",
			Status:     VideoStatusCompleted,
			Limit:      10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("get video", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" || r.URL.Path != "/videos/vid1" {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			}

			writeJSON(w, map[string]interface{}{
				"id":          "vid1",
				"filename":    "test.mp4",
				"status":      "COMPLETED",
				"playlistUrl": "https://cdn.example.com/vid1/playlist.m3u8",
				"createdAt":   "2024-01-01T00:00:00Z",
				"updatedAt":   "2024-01-01T00:00:00Z",
			})
		})
		defer server.Close()

		result, err := client.Videos.Get(context.Background(), "vid1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID != "vid1" {
			t.Errorf("expected id 'vid1', got %s", result.ID)
		}
		if result.PlaylistURL == nil {
			t.Error("expected playlistUrl to be set")
		}
	})

	t.Run("delete video", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" || r.URL.Path != "/videos/vid1" {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			}

			writeJSON(w, map[string]string{"message": "Deleted"})
		})
		defer server.Close()

		err := client.Videos.Delete(context.Background(), "vid1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("get embed info", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" || r.URL.Path != "/videos/vid1/embed" {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			}

			writeJSON(w, map[string]interface{}{
				"videoId":   "vid1",
				"signature": "sig123",
				"embedPath": "/embed/vid1",
			})
		})
		defer server.Close()

		result, err := client.Videos.GetEmbedInfo(context.Background(), "vid1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.VideoID != "vid1" {
			t.Errorf("expected videoId 'vid1', got %s", result.VideoID)
		}
		if result.Signature != "sig123" {
			t.Errorf("expected signature 'sig123', got %s", result.Signature)
		}
	})
}

func TestWebhooks(t *testing.T) {
	t.Run("list webhooks", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, map[string]interface{}{
				"webhooks": []map[string]interface{}{
					{
						"id":        "wh1",
						"url":       "https://example.com/webhook",
						"events":    []string{"video.completed"},
						"isActive":  true,
						"createdAt": "2024-01-01T00:00:00Z",
						"updatedAt": "2024-01-01T00:00:00Z",
					},
				},
				"availableEvents": []string{"video.uploaded", "video.completed"},
			})
		})
		defer server.Close()

		result, err := client.Webhooks.List(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Webhooks) != 1 {
			t.Errorf("expected 1 webhook, got %d", len(result.Webhooks))
		}
	})

	t.Run("create webhook", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, map[string]interface{}{
				"message": "Webhook created",
				"webhook": map[string]interface{}{
					"id":        "wh1",
					"url":       "https://example.com/webhook",
					"events":    []string{"video.completed"},
					"secret":    "whsec_test123",
					"isActive":  true,
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-01T00:00:00Z",
				},
			})
		})
		defer server.Close()

		result, err := client.Webhooks.Create(context.Background(), &CreateWebhookParams{
			URL:    "https://example.com/webhook",
			Events: []WebhookEventType{WebhookEventVideoCompleted},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Webhook.Secret != "whsec_test123" {
			t.Errorf("expected secret 'whsec_test123', got %s", result.Webhook.Secret)
		}
	})
}

func TestEmbed(t *testing.T) {
	t.Run("get settings", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, map[string]interface{}{
				"allowedDomains": []string{"example.com"},
				"allowLocalhost": true,
				"primaryColor":   "#FF0000",
			})
		})
		defer server.Close()

		result, err := client.Embed.GetSettings(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.PrimaryColor != "#FF0000" {
			t.Errorf("expected primaryColor '#FF0000', got %s", result.PrimaryColor)
		}
	})

	t.Run("sign embed", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, map[string]interface{}{
				"videoId":   "vid1",
				"signature": "sig123",
				"embedUrl":  "https://embed.playvideo.dev/vid1",
				"embedCode": map[string]string{
					"responsive": "<div></div>",
					"fixed":      "<iframe></iframe>",
				},
			})
		})
		defer server.Close()

		result, err := client.Embed.Sign(context.Background(), &SignEmbedParams{VideoID: "vid1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Signature != "sig123" {
			t.Errorf("expected signature 'sig123', got %s", result.Signature)
		}
	})
}

func TestAPIKeys(t *testing.T) {
	t.Run("list api keys", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, map[string]interface{}{
				"apiKeys": []map[string]interface{}{
					{
						"id":        "key1",
						"name":      "Production",
						"keyPrefix": "play_live_abc",
						"createdAt": "2024-01-01T00:00:00Z",
					},
				},
			})
		})
		defer server.Close()

		result, err := client.APIKeys.List(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.APIKeys) != 1 {
			t.Errorf("expected 1 API key, got %d", len(result.APIKeys))
		}
	})

	t.Run("create api key", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, map[string]interface{}{
				"message": "API key created",
				"apiKey": map[string]interface{}{
					"id":        "key1",
					"name":      "New Key",
					"keyPrefix": "play_live_xyz",
					"key":       "play_live_xyz123456789",
					"createdAt": "2024-01-01T00:00:00Z",
				},
			})
		})
		defer server.Close()

		result, err := client.APIKeys.Create(context.Background(), &CreateAPIKeyParams{Name: "New Key"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.APIKey.Key != "play_live_xyz123456789" {
			t.Errorf("expected key 'play_live_xyz123456789', got %s", result.APIKey.Key)
		}
	})
}

func TestAccount(t *testing.T) {
	t.Run("get account", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, map[string]interface{}{
				"id":        "acc1",
				"email":     "user@example.com",
				"plan":      "PRO",
				"createdAt": "2024-01-01T00:00:00Z",
			})
		})
		defer server.Close()

		result, err := client.Account.Get(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Email != "user@example.com" {
			t.Errorf("expected email 'user@example.com', got %s", result.Email)
		}
		if result.Plan != "PRO" {
			t.Errorf("expected plan 'PRO', got %s", result.Plan)
		}
	})
}

func TestUsage(t *testing.T) {
	t.Run("get usage", func(t *testing.T) {
		server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, map[string]interface{}{
				"plan": "PRO",
				"usage": map[string]interface{}{
					"videosThisMonth": 50,
					"videosLimit":     500,
				},
				"limits": map[string]interface{}{
					"maxFileSizeMB": 500,
					"apiAccess":     true,
				},
			})
		})
		defer server.Close()

		result, err := client.Usage.Get(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Plan != "PRO" {
			t.Errorf("expected plan 'PRO', got %s", result.Plan)
		}
		if result.Usage.VideosThisMonth != 50 {
			t.Errorf("expected videosThisMonth 50, got %d", result.Usage.VideosThisMonth)
		}
	})
}
