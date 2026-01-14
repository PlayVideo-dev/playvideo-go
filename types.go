package playvideo

import "time"

// VideoStatus represents the processing status of a video
type VideoStatus string

const (
	VideoStatusPending    VideoStatus = "PENDING"
	VideoStatusProcessing VideoStatus = "PROCESSING"
	VideoStatusCompleted  VideoStatus = "COMPLETED"
	VideoStatusFailed     VideoStatus = "FAILED"
)

// Plan represents a subscription plan
type Plan string

const (
	PlanFree     Plan = "FREE"
	PlanPro      Plan = "PRO"
	PlanBusiness Plan = "BUSINESS"
)

// LogoPosition represents the position of a logo in the video player
type LogoPosition string

const (
	LogoPositionTopLeft     LogoPosition = "top-left"
	LogoPositionTopRight    LogoPosition = "top-right"
	LogoPositionBottomLeft  LogoPosition = "bottom-left"
	LogoPositionBottomRight LogoPosition = "bottom-right"
)

// WebhookEventType represents the type of webhook event
type WebhookEventType string

const (
	WebhookEventVideoUploaded     WebhookEventType = "video.uploaded"
	WebhookEventVideoProcessing   WebhookEventType = "video.processing"
	WebhookEventVideoCompleted    WebhookEventType = "video.completed"
	WebhookEventVideoFailed       WebhookEventType = "video.failed"
	WebhookEventCollectionCreated WebhookEventType = "collection.created"
	WebhookEventCollectionDeleted WebhookEventType = "collection.deleted"
)

// ProgressStage represents the stage in video processing progress
type ProgressStage string

const (
	ProgressStagePending    ProgressStage = "pending"
	ProgressStageProcessing ProgressStage = "processing"
	ProgressStageCompleted  ProgressStage = "completed"
	ProgressStageFailed     ProgressStage = "failed"
	ProgressStageTimeout    ProgressStage = "timeout"
)

// ===========================================
// Collection Types
// ===========================================

// Collection represents a video collection
type Collection struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	VideoCount  int       `json:"videoCount"`
	StorageUsed int64     `json:"storageUsed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CollectionWithVideos is a collection with its videos
type CollectionWithVideos struct {
	Collection
	Videos []Video `json:"videos"`
}

// CollectionListResponse is the response for listing collections
type CollectionListResponse struct {
	Collections []Collection `json:"collections"`
}

// CreateCollectionParams are the parameters for creating a collection
type CreateCollectionParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// ===========================================
// Video Types
// ===========================================

// VideoCollection is the collection info embedded in a video
type VideoCollection struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// Video represents a video
type Video struct {
	ID            string           `json:"id"`
	Filename      string           `json:"filename"`
	Status        VideoStatus      `json:"status"`
	Duration      *float64         `json:"duration"`
	OriginalSize  int64            `json:"originalSize"`
	ProcessedSize *int64           `json:"processedSize"`
	PlaylistURL   *string          `json:"playlistUrl"`
	ThumbnailURL  *string          `json:"thumbnailUrl"`
	PreviewURL    *string          `json:"previewUrl"`
	Resolutions   []string         `json:"resolutions"`
	ErrorMessage  *string          `json:"errorMessage"`
	Collection    *VideoCollection `json:"collection,omitempty"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

// VideoListResponse is the response for listing videos
type VideoListResponse struct {
	Videos []Video `json:"videos"`
}

// VideoListParams are the parameters for listing videos
type VideoListParams struct {
	Collection string      `url:"collection,omitempty"`
	Status     VideoStatus `url:"status,omitempty"`
	Limit      int         `url:"limit,omitempty"`
	Offset     int         `url:"offset,omitempty"`
}

// UploadVideo is the video info returned after upload
type UploadVideo struct {
	ID         string      `json:"id"`
	Filename   string      `json:"filename"`
	Status     VideoStatus `json:"status"`
	Collection string      `json:"collection"`
}

// UploadResponse is the response from uploading a video
type UploadResponse struct {
	Message string      `json:"message"`
	Video   UploadVideo `json:"video"`
}

// UploadProgress represents upload progress
type UploadProgress struct {
	Loaded  int64 `json:"loaded"`
	Total   int64 `json:"total"`
	Percent int   `json:"percent"`
}

// UploadProgressFunc is a callback for upload progress
type UploadProgressFunc func(progress UploadProgress)

// VideoEmbedInfo contains embed information for a video
type VideoEmbedInfo struct {
	VideoID   string `json:"videoId"`
	Signature string `json:"signature"`
	EmbedPath string `json:"embedPath"`
}

// ProgressEvent represents a video processing progress event
type ProgressEvent struct {
	Stage         ProgressStage `json:"stage"`
	Message       string        `json:"message,omitempty"`
	Error         string        `json:"error,omitempty"`
	PlaylistURL   string        `json:"playlistUrl,omitempty"`
	ThumbnailURL  *string       `json:"thumbnailUrl,omitempty"`
	PreviewURL    *string       `json:"previewUrl,omitempty"`
	Duration      float64       `json:"duration,omitempty"`
	ProcessedSize *int64        `json:"processedSize,omitempty"`
	Resolutions   []string      `json:"resolutions,omitempty"`
}

// ===========================================
// Webhook Types
// ===========================================

// Webhook represents a webhook configuration
type Webhook struct {
	ID        string             `json:"id"`
	URL       string             `json:"url"`
	Events    []WebhookEventType `json:"events"`
	IsActive  bool               `json:"isActive"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

// WebhookWithSecret is a webhook with its secret (returned on creation)
type WebhookWithSecret struct {
	Webhook
	Secret string `json:"secret"`
}

// WebhookListResponse is the response for listing webhooks
type WebhookListResponse struct {
	Webhooks        []Webhook          `json:"webhooks"`
	AvailableEvents []WebhookEventType `json:"availableEvents"`
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID           string           `json:"id"`
	Event        WebhookEventType `json:"event"`
	StatusCode   *int             `json:"statusCode"`
	Error        *string          `json:"error"`
	AttemptCount int              `json:"attemptCount"`
	DeliveredAt  *time.Time       `json:"deliveredAt"`
	CreatedAt    time.Time        `json:"createdAt"`
}

// WebhookWithDeliveries is a webhook with recent delivery attempts
type WebhookWithDeliveries struct {
	Webhook
	RecentDeliveries []WebhookDelivery `json:"recentDeliveries"`
}

// CreateWebhookParams are the parameters for creating a webhook
type CreateWebhookParams struct {
	URL    string             `json:"url"`
	Events []WebhookEventType `json:"events"`
}

// UpdateWebhookParams are the parameters for updating a webhook
type UpdateWebhookParams struct {
	URL      *string            `json:"url,omitempty"`
	Events   []WebhookEventType `json:"events,omitempty"`
	IsActive *bool              `json:"isActive,omitempty"`
}

// CreateWebhookResponse is the response from creating a webhook
type CreateWebhookResponse struct {
	Message string            `json:"message"`
	Webhook WebhookWithSecret `json:"webhook"`
}

// TestWebhookResponse is the response from testing a webhook
type TestWebhookResponse struct {
	Message    string  `json:"message"`
	StatusCode *int    `json:"statusCode,omitempty"`
	Error      *string `json:"error,omitempty"`
}

// ===========================================
// Embed Types
// ===========================================

// EmbedSettings represents player embed settings
type EmbedSettings struct {
	AllowedDomains      []string     `json:"allowedDomains"`
	AllowLocalhost      bool         `json:"allowLocalhost"`
	PrimaryColor        string       `json:"primaryColor"`
	AccentColor         string       `json:"accentColor"`
	LogoURL             *string      `json:"logoUrl"`
	LogoPosition        LogoPosition `json:"logoPosition"`
	LogoOpacity         float64      `json:"logoOpacity"`
	ShowPlaybackSpeed   bool         `json:"showPlaybackSpeed"`
	ShowQualitySelector bool         `json:"showQualitySelector"`
	ShowFullscreen      bool         `json:"showFullscreen"`
	ShowVolume          bool         `json:"showVolume"`
	ShowProgress        bool         `json:"showProgress"`
	ShowTime            bool         `json:"showTime"`
	ShowKeyboardHints   bool         `json:"showKeyboardHints"`
	Autoplay            bool         `json:"autoplay"`
	Muted               bool         `json:"muted"`
	Loop                bool         `json:"loop"`
}

// UpdateEmbedSettingsParams are the parameters for updating embed settings
type UpdateEmbedSettingsParams struct {
	AllowedDomains      []string      `json:"allowedDomains,omitempty"`
	AllowLocalhost      *bool         `json:"allowLocalhost,omitempty"`
	PrimaryColor        *string       `json:"primaryColor,omitempty"`
	AccentColor         *string       `json:"accentColor,omitempty"`
	LogoURL             *string       `json:"logoUrl,omitempty"`
	LogoPosition        *LogoPosition `json:"logoPosition,omitempty"`
	LogoOpacity         *float64      `json:"logoOpacity,omitempty"`
	ShowPlaybackSpeed   *bool         `json:"showPlaybackSpeed,omitempty"`
	ShowQualitySelector *bool         `json:"showQualitySelector,omitempty"`
	ShowFullscreen      *bool         `json:"showFullscreen,omitempty"`
	ShowVolume          *bool         `json:"showVolume,omitempty"`
	ShowProgress        *bool         `json:"showProgress,omitempty"`
	ShowTime            *bool         `json:"showTime,omitempty"`
	ShowKeyboardHints   *bool         `json:"showKeyboardHints,omitempty"`
	Autoplay            *bool         `json:"autoplay,omitempty"`
	Muted               *bool         `json:"muted,omitempty"`
	Loop                *bool         `json:"loop,omitempty"`
}

// UpdateEmbedSettingsResponse is the response from updating embed settings
type UpdateEmbedSettingsResponse struct {
	Message  string        `json:"message"`
	Settings EmbedSettings `json:"settings"`
}

// SignEmbedParams are the parameters for signing an embed URL
type SignEmbedParams struct {
	VideoID string  `json:"videoId"`
	BaseURL *string `json:"baseUrl,omitempty"`
}

// EmbedCode contains HTML embed code snippets
type EmbedCode struct {
	Responsive string `json:"responsive"`
	Fixed      string `json:"fixed"`
}

// SignEmbedResponse is the response from signing an embed URL
type SignEmbedResponse struct {
	VideoID   string    `json:"videoId"`
	Signature string    `json:"signature"`
	EmbedURL  string    `json:"embedUrl"`
	EmbedCode EmbedCode `json:"embedCode"`
}

// ===========================================
// API Key Types
// ===========================================

// APIKey represents an API key
type APIKey struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"keyPrefix"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}

// APIKeyWithKey is an API key with the full key value (returned on creation)
type APIKeyWithKey struct {
	APIKey
	Key string `json:"key"`
}

// APIKeyListResponse is the response for listing API keys
type APIKeyListResponse struct {
	APIKeys []APIKey `json:"apiKeys"`
}

// CreateAPIKeyParams are the parameters for creating an API key
type CreateAPIKeyParams struct {
	Name string `json:"name"`
}

// CreateAPIKeyResponse is the response from creating an API key
type CreateAPIKeyResponse struct {
	Message string        `json:"message"`
	APIKey  APIKeyWithKey `json:"apiKey"`
}

// ===========================================
// Account Types
// ===========================================

// Account represents a user account
type Account struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	Name           *string   `json:"name"`
	Plan           Plan      `json:"plan"`
	AllowedDomains []string  `json:"allowedDomains"`
	AllowLocalhost bool      `json:"allowLocalhost"`
	R2BucketName   *string   `json:"r2BucketName"`
	R2BucketRegion *string   `json:"r2BucketRegion"`
	CreatedAt      time.Time `json:"createdAt"`
}

// UpdateAccountParams are the parameters for updating an account
type UpdateAccountParams struct {
	AllowedDomains []string `json:"allowedDomains,omitempty"`
	AllowLocalhost *bool    `json:"allowLocalhost,omitempty"`
}

// UpdateAccountResponse is the response from updating an account
type UpdateAccountResponse struct {
	Message string  `json:"message"`
	Account Account `json:"account"`
}

// ===========================================
// Usage Types
// ===========================================

// UsageInfo contains usage information
type UsageInfo struct {
	VideosThisMonth  int         `json:"videosThisMonth"`
	VideosLimit      interface{} `json:"videosLimit"` // int or "unlimited"
	StorageUsedBytes int64       `json:"storageUsedBytes"`
	StorageUsedGB    string      `json:"storageUsedGB"`
	StorageLimitGB   int         `json:"storageLimitGB"`
}

// UsageLimits contains plan limits
type UsageLimits struct {
	MaxFileSizeMB      int      `json:"maxFileSizeMB"`
	MaxDurationMinutes int      `json:"maxDurationMinutes"`
	Resolutions        []string `json:"resolutions"`
	APIAccess          bool     `json:"apiAccess"`
	Webhooks           bool     `json:"webhooks"`
	DeliveryGB         int      `json:"deliveryGB"`
}

// Usage contains account usage and limits
type Usage struct {
	Plan   Plan        `json:"plan"`
	Usage  UsageInfo   `json:"usage"`
	Limits UsageLimits `json:"limits"`
}

// ===========================================
// Webhook Payload Types
// ===========================================

// WebhookPayload is the structure of a webhook payload
type WebhookPayload struct {
	Event     WebhookEventType       `json:"event"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}
