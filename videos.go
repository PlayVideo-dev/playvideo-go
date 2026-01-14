package playvideo

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// VideosService handles video-related API operations
type VideosService struct {
	client *httpClient
}

// List retrieves videos with optional filters
func (s *VideosService) List(ctx context.Context, params *VideoListParams) (*VideoListResponse, error) {
	var result VideoListResponse
	err := s.client.getWithParams(ctx, "/videos", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a video by ID
func (s *VideosService) Get(ctx context.Context, id string) (*Video, error) {
	var result Video
	err := s.client.get(ctx, fmt.Sprintf("/videos/%s", encodeURLPath(id)), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Upload uploads a video from an io.Reader
func (s *VideosService) Upload(ctx context.Context, file io.Reader, filename, collection string, progressFn UploadProgressFunc) (*UploadResponse, error) {
	return s.client.upload(ctx, "/videos", file, filename, collection, progressFn)
}

// UploadFile uploads a video from a file path
func (s *VideosService) UploadFile(ctx context.Context, filePath, collection string, progressFn UploadProgressFunc) (*UploadResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("failed to open file: %v", err)}
	}
	defer file.Close()

	filename := filepath.Base(filePath)
	return s.Upload(ctx, file, filename, collection, progressFn)
}

// Delete deletes a video by ID
func (s *VideosService) Delete(ctx context.Context, id string) error {
	var result map[string]string
	err := s.client.delete(ctx, fmt.Sprintf("/videos/%s", encodeURLPath(id)), &result)
	if err != nil {
		return err
	}
	return nil
}

// GetEmbedInfo retrieves embed information for a video
func (s *VideosService) GetEmbedInfo(ctx context.Context, id string) (*VideoEmbedInfo, error) {
	var result VideoEmbedInfo
	err := s.client.get(ctx, fmt.Sprintf("/videos/%s/embed", encodeURLPath(id)), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// WatchProgress watches video processing progress via SSE
// Returns a channel that receives progress events
// The channel is closed when processing completes, fails, or times out
func (s *VideosService) WatchProgress(ctx context.Context, id string) (<-chan *ProgressEvent, <-chan error) {
	return s.client.sseRequest(ctx, fmt.Sprintf("/videos/%s/progress", encodeURLPath(id)))
}
