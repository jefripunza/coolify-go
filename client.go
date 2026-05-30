package coolify

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://app.coolify.io/api/v1"
	userAgent      = "go-coolify/v0.1.0"
)

// Client coordinates communication with the Coolify REST API.
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	token      string
	userAgent  string

	// Services for different endpoints
	Applications *ApplicationsService
	Servers      *ServersService
	Projects     *ProjectsService
	Databases    *DatabasesService
	Deployments  *DeploymentsService
	PrivateKeys  *PrivateKeysService
	GitHubApps   *GitHubAppsService
	Services     *ServicesService
	Teams        *TeamsService
	System       *SystemService
}

// NewClient initializes a new Coolify API client with a base URL and Bearer token.
func NewClient(baseURL, token string) *Client {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	// Automatically append /api/v1 if it is missing from the base URL
	if !strings.Contains(baseURL, "/api/v1") {
		baseURL = strings.TrimSuffix(baseURL, "/") + "/api/v1"
	}
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	parsedBaseURL, _ := url.Parse(baseURL)
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    parsedBaseURL,
		token:      strings.TrimSpace(token),
		userAgent:  userAgent,
	}

	// Initialize services
	c.Applications = &ApplicationsService{client: c}
	c.Servers = &ServersService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.Databases = &DatabasesService{client: c}
	c.Deployments = &DeploymentsService{client: c}
	c.PrivateKeys = &PrivateKeysService{client: c}
	c.GitHubApps = &GitHubAppsService{client: c}
	c.Services = &ServicesService{client: c}
	c.Teams = &TeamsService{client: c}
	c.System = &SystemService{client: c}

	return c
}

// newRequest constructs an HTTP Request with proper URL, Content-Type, User-Agent and body.
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	// Trim leading slash from path to resolve relative reference correctly
	path = strings.TrimPrefix(path, "/")
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return req, nil
}

// do sends the API request and unmarshals the response body to v if v is not nil.
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Handle context cancellation / deadline errors
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return resp, parseResponseError(resp.StatusCode, bodyBytes)
	}

	if v != nil {
		var err error
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else if strPtr, ok := v.(*string); ok {
			bodyBytes, err2 := io.ReadAll(resp.Body)
			if err2 != nil {
				return resp, err2
			}
			rawStr := strings.TrimSpace(string(bodyBytes))
			if strings.HasPrefix(rawStr, `"`) && strings.HasSuffix(rawStr, `"`) && len(rawStr) >= 2 {
				var parsedStr string
				if json.Unmarshal(bodyBytes, &parsedStr) == nil {
					rawStr = parsedStr
				} else {
					rawStr = rawStr[1 : len(rawStr)-1]
				}
			}
			*strPtr = rawStr
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
		if err != nil && err != io.EOF {
			return resp, err
		}
	}

	return resp, nil
}

// String is a helper helper to allocate a new string value
// for field pointers in Go structs.
func String(v string) *string {
	return &v
}

// Int is a helper helper to allocate a new int value
// for field pointers in Go structs.
func Int(v int) *int {
	return &v
}

// Bool is a helper helper to allocate a new bool value
// for field pointers in Go structs.
func Bool(v bool) *bool {
	return &v
}
