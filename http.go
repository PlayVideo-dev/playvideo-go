package playvideo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

// httpClient is the internal HTTP client wrapper
type httpClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// newHTTPClient creates a new HTTP client wrapper
func newHTTPClient(apiKey, baseURL string, client *http.Client) *httpClient {
	return &httpClient{
		apiKey:     apiKey,
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: client,
	}
}

// get performs a GET request
func (c *httpClient) get(ctx context.Context, endpoint string, result interface{}) error {
	return c.request(ctx, http.MethodGet, endpoint, nil, result)
}

// getWithParams performs a GET request with query parameters
func (c *httpClient) getWithParams(ctx context.Context, endpoint string, params interface{}, result interface{}) error {
	if params != nil {
		v, err := query.Values(params)
		if err != nil {
			return &Error{Message: fmt.Sprintf("failed to encode query params: %v", err)}
		}
		if encoded := v.Encode(); encoded != "" {
			endpoint = endpoint + "?" + encoded
		}
	}
	return c.request(ctx, http.MethodGet, endpoint, nil, result)
}

// post performs a POST request with JSON body
func (c *httpClient) post(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.request(ctx, http.MethodPost, endpoint, body, result)
}

// patch performs a PATCH request with JSON body
func (c *httpClient) patch(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.request(ctx, http.MethodPatch, endpoint, body, result)
}

// delete performs a DELETE request
func (c *httpClient) delete(ctx context.Context, endpoint string, result interface{}) error {
	return c.request(ctx, http.MethodDelete, endpoint, nil, result)
}

// request performs an HTTP request
func (c *httpClient) request(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	reqURL := c.baseURL + endpoint

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return &Error{Message: fmt.Sprintf("failed to marshal request body: %v", err)}
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return &NetworkError{Message: fmt.Sprintf("failed to create request: %v", err), Err: err}
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return &TimeoutError{Message: "request timed out", Err: err}
		}
		if ctx.Err() == context.Canceled {
			return &NetworkError{Message: "request canceled", Err: err}
		}
		return &NetworkError{Message: fmt.Sprintf("request failed: %v", err), Err: err}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &NetworkError{Message: fmt.Sprintf("failed to read response: %v", err), Err: err}
	}

	requestID := resp.Header.Get("X-Request-ID")

	if resp.StatusCode >= 400 {
		var errResp apiErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			errResp.Error = string(respBody)
		}
		return parseAPIError(resp.StatusCode, &errResp, requestID)
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return &Error{Message: fmt.Sprintf("failed to parse response: %v", err)}
		}
	}

	return nil
}

// upload performs a multipart file upload
func (c *httpClient) upload(ctx context.Context, endpoint string, file io.Reader, filename, collection string, progressFn UploadProgressFunc) (*UploadResponse, error) {
	reqURL := c.baseURL + endpoint

	// Create a pipe for streaming the multipart form
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Channel to capture errors from the goroutine
	errChan := make(chan error, 1)

	go func() {
		defer pw.Close()
		defer writer.Close()

		// Add collection field
		if err := writer.WriteField("collection", collection); err != nil {
			errChan <- err
			return
		}

		// Create file field
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			errChan <- err
			return
		}

		// Copy file data
		if _, err := io.Copy(part, file); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, pr)
	if err != nil {
		return nil, &NetworkError{Message: fmt.Sprintf("failed to create request: %v", err), Err: err}
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, &TimeoutError{Message: "upload timed out", Err: err}
		}
		return nil, &NetworkError{Message: fmt.Sprintf("upload failed: %v", err), Err: err}
	}
	defer resp.Body.Close()

	// Wait for the goroutine to finish
	if writeErr := <-errChan; writeErr != nil {
		return nil, &Error{Message: fmt.Sprintf("failed to write multipart form: %v", writeErr)}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &NetworkError{Message: fmt.Sprintf("failed to read response: %v", err), Err: err}
	}

	requestID := resp.Header.Get("X-Request-ID")

	if resp.StatusCode >= 400 {
		var errResp apiErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			errResp.Error = string(respBody)
		}
		return nil, parseAPIError(resp.StatusCode, &errResp, requestID)
	}

	var result UploadResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, &Error{Message: fmt.Sprintf("failed to parse response: %v", err)}
	}

	// Call progress callback with 100% on success
	if progressFn != nil {
		progressFn(UploadProgress{Loaded: 1, Total: 1, Percent: 100})
	}

	return &result, nil
}

// sseRequest performs an SSE request and returns events via a channel
func (c *httpClient) sseRequest(ctx context.Context, endpoint string) (<-chan *ProgressEvent, <-chan error) {
	events := make(chan *ProgressEvent)
	errChan := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(errChan)

		reqURL := c.baseURL + endpoint

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			errChan <- &NetworkError{Message: fmt.Sprintf("failed to create request: %v", err), Err: err}
			return
		}

		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "text/event-stream")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				errChan <- &TimeoutError{Message: "request timed out", Err: err}
				return
			}
			errChan <- &NetworkError{Message: fmt.Sprintf("request failed: %v", err), Err: err}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(resp.Body)
			var errResp apiErrorResponse
			if err := json.Unmarshal(body, &errResp); err != nil {
				errResp.Error = string(body)
			}
			errChan <- parseAPIError(resp.StatusCode, &errResp, resp.Header.Get("X-Request-ID"))
			return
		}

		// Parse SSE stream
		decoder := newSSEDecoder(resp.Body)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				event, err := decoder.Next()
				if err == io.EOF {
					return
				}
				if err != nil {
					errChan <- &NetworkError{Message: fmt.Sprintf("failed to read SSE stream: %v", err), Err: err}
					return
				}

				if event != nil {
					events <- event
					// Stop on terminal states
					if event.Stage == ProgressStageCompleted || event.Stage == ProgressStageFailed || event.Stage == ProgressStageTimeout {
						return
					}
				}
			}
		}
	}()

	return events, errChan
}

// sseDecoder decodes Server-Sent Events
type sseDecoder struct {
	reader    *bytes.Buffer
	rawReader io.Reader
}

func newSSEDecoder(r io.Reader) *sseDecoder {
	return &sseDecoder{
		reader:    bytes.NewBuffer(nil),
		rawReader: r,
	}
}

func (d *sseDecoder) Next() (*ProgressEvent, error) {
	// Read more data if needed
	buf := make([]byte, 4096)
	for {
		// Check if we have a complete event in buffer
		data := d.reader.String()
		if idx := strings.Index(data, "\n\n"); idx >= 0 {
			eventData := data[:idx]
			d.reader.Reset()
			d.reader.WriteString(data[idx+2:])

			// Parse the event
			for _, line := range strings.Split(eventData, "\n") {
				if strings.HasPrefix(line, "data: ") {
					jsonData := strings.TrimPrefix(line, "data: ")
					var event ProgressEvent
					if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
						continue // Skip malformed events
					}
					return &event, nil
				}
			}
			continue // No data line found, try next event
		}

		// Read more data
		n, err := d.rawReader.Read(buf)
		if n > 0 {
			d.reader.Write(buf[:n])
		}
		if err != nil {
			return nil, err
		}
	}
}

// encodeURLPath encodes a string for use in a URL path
func encodeURLPath(s string) string {
	return url.PathEscape(s)
}
