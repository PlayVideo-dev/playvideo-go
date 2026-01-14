# PlayVideo Go SDK

[![CI](https://github.com/PlayVideo-dev/playvideo-go/actions/workflows/ci.yml/badge.svg)](https://github.com/PlayVideo-dev/playvideo-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/PlayVideo-dev/playvideo-go.svg)](https://pkg.go.dev/github.com/PlayVideo-dev/playvideo-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/PlayVideo-dev/playvideo-go)](https://goreportcard.com/report/github.com/PlayVideo-dev/playvideo-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go SDK for the [PlayVideo](https://playvideo.dev) API - Video hosting for developers.

## Installation

```bash
go get github.com/PlayVideo-dev/playvideo-go
```

## Requirements

- Go 1.21 or later

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/PlayVideo-dev/playvideo-go"
)

func main() {
    // Create a client with your API key
    client := playvideo.NewClient("play_live_xxx")

    ctx := context.Background()

    // List collections
    collections, err := client.Collections.List(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d collections\n", len(collections.Collections))

    // Upload a video
    upload, err := client.Videos.UploadFile(ctx, "./video.mp4", "my-collection", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Uploaded video: %s\n", upload.Video.ID)
}
```

## Configuration

```go
// With custom options
client := playvideo.NewClient("play_live_xxx",
    playvideo.WithBaseURL("https://api.playvideo.dev/api/v1"),
    playvideo.WithTimeout(60 * time.Second),
    playvideo.WithHTTPClient(&http.Client{
        Transport: &http.Transport{
            MaxIdleConns: 10,
        },
    }),
)
```

## Usage Examples

### Collections

```go
ctx := context.Background()

// List all collections
collections, err := client.Collections.List(ctx)

// Create a collection
collection, err := client.Collections.Create(ctx, &playvideo.CreateCollectionParams{
    Name:        "My Collection",
    Description: playvideo.String("Optional description"),
})

// Get a collection with videos
collection, err := client.Collections.Get(ctx, "my-collection")

// Delete a collection
err := client.Collections.Delete(ctx, "my-collection")
```

### Videos

```go
ctx := context.Background()

// List videos with filters
videos, err := client.Videos.List(ctx, &playvideo.VideoListParams{
    Collection: "my-collection",
    Status:     playvideo.VideoStatusCompleted,
    Limit:      10,
})

// Upload from file path
upload, err := client.Videos.UploadFile(ctx, "./video.mp4", "my-collection", func(p playvideo.UploadProgress) {
    fmt.Printf("Upload progress: %d%%\n", p.Percent)
})

// Upload from io.Reader
file, _ := os.Open("./video.mp4")
defer file.Close()
upload, err := client.Videos.Upload(ctx, file, "video.mp4", "my-collection", nil)

// Get a video
video, err := client.Videos.Get(ctx, "video-id")

// Delete a video
err := client.Videos.Delete(ctx, "video-id")

// Get embed info
embedInfo, err := client.Videos.GetEmbedInfo(ctx, "video-id")
```

### Watch Processing Progress (SSE)

```go
ctx := context.Background()

// Watch video processing progress
events, errs := client.Videos.WatchProgress(ctx, "video-id")

for {
    select {
    case event, ok := <-events:
        if !ok {
            return // Channel closed
        }
        fmt.Printf("Stage: %s, Message: %s\n", event.Stage, event.Message)
        
        if event.Stage == playvideo.ProgressStageCompleted {
            fmt.Printf("Processing complete! Playlist: %s\n", event.PlaylistURL)
        }
    case err := <-errs:
        if err != nil {
            log.Fatal(err)
        }
    }
}
```

### Webhooks

```go
ctx := context.Background()

// List webhooks
webhooks, err := client.Webhooks.List(ctx)

// Create a webhook
webhook, err := client.Webhooks.Create(ctx, &playvideo.CreateWebhookParams{
    URL: "https://example.com/webhook",
    Events: []playvideo.WebhookEventType{
        playvideo.WebhookEventVideoCompleted,
        playvideo.WebhookEventVideoFailed,
    },
})
fmt.Printf("Webhook secret: %s\n", webhook.Webhook.Secret)

// Update a webhook
updated, err := client.Webhooks.Update(ctx, "webhook-id", &playvideo.UpdateWebhookParams{
    IsActive: playvideo.Bool(false),
})

// Test a webhook
result, err := client.Webhooks.Test(ctx, "webhook-id")

// Delete a webhook
err := client.Webhooks.Delete(ctx, "webhook-id")
```

### Webhook Signature Verification

```go
import "github.com/PlayVideo-dev/playvideo-go"

func handleWebhook(w http.ResponseWriter, r *http.Request) {
    payload, _ := io.ReadAll(r.Body)
    signature := r.Header.Get("X-PlayVideo-Signature")
    timestamp := r.Header.Get("X-PlayVideo-Timestamp")
    secret := "whsec_xxx" // Your webhook secret

    // Verify and parse the event
    event, err := playvideo.ConstructEvent(payload, signature, timestamp, secret)
    if err != nil {
        if playvideo.IsWebhookSignatureError(err) {
            http.Error(w, "Invalid signature", http.StatusBadRequest)
            return
        }
        http.Error(w, "Error", http.StatusInternalServerError)
        return
    }

    // Handle the event
    switch event.Event {
    case playvideo.WebhookEventVideoCompleted:
        fmt.Printf("Video completed: %v\n", event.Data)
    case playvideo.WebhookEventVideoFailed:
        fmt.Printf("Video failed: %v\n", event.Data)
    }

    w.WriteHeader(http.StatusOK)
}
```

### Embed Settings

```go
ctx := context.Background()

// Get embed settings
settings, err := client.Embed.GetSettings(ctx)

// Update embed settings
updated, err := client.Embed.UpdateSettings(ctx, &playvideo.UpdateEmbedSettingsParams{
    PrimaryColor:   playvideo.String("#FF0000"),
    Autoplay:       playvideo.Bool(true),
    AllowLocalhost: playvideo.Bool(true),
})

// Sign an embed URL
signed, err := client.Embed.Sign(ctx, &playvideo.SignEmbedParams{
    VideoID: "video-id",
})
fmt.Printf("Embed URL: %s\n", signed.EmbedURL)
```

### API Keys

```go
ctx := context.Background()

// List API keys
keys, err := client.APIKeys.List(ctx)

// Create an API key
key, err := client.APIKeys.Create(ctx, &playvideo.CreateAPIKeyParams{
    Name: "My API Key",
})
fmt.Printf("New API key: %s\n", key.APIKey.Key) // Only shown once!

// Delete an API key
err := client.APIKeys.Delete(ctx, "key-id")
```

### Account

```go
ctx := context.Background()

// Get account info
account, err := client.Account.Get(ctx)
fmt.Printf("Plan: %s\n", account.Plan)

// Update account
updated, err := client.Account.Update(ctx, &playvideo.UpdateAccountParams{
    AllowedDomains: []string{"example.com", "*.example.com"},
    AllowLocalhost: playvideo.Bool(true),
})
```

### Usage

```go
ctx := context.Background()

// Get usage info
usage, err := client.Usage.Get(ctx)
fmt.Printf("Videos this month: %d\n", usage.Usage.VideosThisMonth)
fmt.Printf("Storage used: %s GB\n", usage.Usage.StorageUsedGB)
```

## Error Handling

The SDK provides typed errors for easy handling:

```go
video, err := client.Videos.Get(ctx, "non-existent-id")
if err != nil {
    if playvideo.IsNotFoundError(err) {
        fmt.Println("Video not found")
    } else if playvideo.IsAuthenticationError(err) {
        fmt.Println("Invalid API key")
    } else if playvideo.IsRateLimitError(err) {
        fmt.Println("Rate limited, try again later")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

Available error types:
- `AuthenticationError` - Invalid or missing API key (401)
- `AuthorizationError` - Insufficient permissions (403)
- `NotFoundError` - Resource not found (404)
- `ValidationError` - Invalid request parameters (400/422)
- `ConflictError` - Resource conflict (409)
- `RateLimitError` - Too many requests (429)
- `ServerError` - Server error (5xx)
- `NetworkError` - Connection issues
- `TimeoutError` - Request timeout
- `WebhookSignatureError` - Invalid webhook signature

## Helper Functions

```go
// Create pointer to string
s := playvideo.String("value")

// Create pointer to bool  
b := playvideo.Bool(true)

// Create pointer to int
i := playvideo.Int(42)

// Create pointer to float64
f := playvideo.Float64(3.14)
```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Links

- [Documentation](https://playvideo.dev/docs)
- [API Reference](https://playvideo.dev/docs/api)
- [Dashboard](https://playvideo.dev/dashboard)
- [GitHub](https://github.com/PlayVideo-dev/playvideo-go)
