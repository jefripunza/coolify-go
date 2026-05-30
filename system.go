package coolify

import (
	"context"
	"net/http"
)

// SystemService handles communication with instance-level utility and configuration endpoints.
type SystemService struct {
	client *Client
}

// VersionResponse represents the API version structure.
type VersionResponse struct {
	Version string `json:"version"`
}

// Version retrieves the active Coolify instance version.
func (s *SystemService) Version(ctx context.Context) (string, error) {
	req, err := s.client.newRequest(http.MethodGet, "version", nil)
	if err != nil {
		return "", err
	}

	var version string
	_, err = s.client.do(ctx, req, &version)
	return version, err
}

// Health checks the health and responsiveness of the Coolify instance.
func (s *SystemService) Health(ctx context.Context) (string, error) {
	req, err := s.client.newRequest(http.MethodGet, "health", nil)
	if err != nil {
		return "", err
	}

	var health string
	_, err = s.client.do(ctx, req, &health)
	return health, err
}

// EnableAPI enables external API access on the instance.
func (s *SystemService) EnableAPI(ctx context.Context) (string, error) {
	req, err := s.client.newRequest(http.MethodGet, "enable", nil)
	if err != nil {
		return "", err
	}

	var message string
	_, err = s.client.do(ctx, req, &message)
	return message, err
}

// DisableAPI disables external API access on the instance.
func (s *SystemService) DisableAPI(ctx context.Context) (string, error) {
	req, err := s.client.newRequest(http.MethodGet, "disable", nil)
	if err != nil {
		return "", err
	}

	var message string
	_, err = s.client.do(ctx, req, &message)
	return message, err
}

// EnableMCPServer enables the Model Context Protocol (MCP) server for local AI assistant integrations.
func (s *SystemService) EnableMCPServer(ctx context.Context) (string, error) {
	req, err := s.client.newRequest(http.MethodPost, "mcp/enable", nil)
	if err != nil {
		return "", err
	}

	var message string
	_, err = s.client.do(ctx, req, &message)
	return message, err
}

// DisableMCPServer disables the Model Context Protocol (MCP) server on the instance.
func (s *SystemService) DisableMCPServer(ctx context.Context) (string, error) {
	req, err := s.client.newRequest(http.MethodPost, "mcp/disable", nil)
	if err != nil {
		return "", err
	}

	var message string
	_, err = s.client.do(ctx, req, &message)
	return message, err
}
