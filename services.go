package coolify

import (
	"context"
	"fmt"
	"net/http"
)

// ServicesService handles communication with one-click service endpoints of the Coolify API.
type ServicesService struct {
	client *Client
}

// List retrieves all one-click services.
func (s *ServicesService) List(ctx context.Context) ([]Service, error) {
	req, err := s.client.newRequest(http.MethodGet, "services", nil)
	if err != nil {
		return nil, err
	}

	var svcs []Service
	_, err = s.client.do(ctx, req, &svcs)
	return svcs, err
}

// Get retrieves a specific one-click service configuration.
func (s *ServicesService) Get(ctx context.Context, uuid string) (*Service, error) {
	path := fmt.Sprintf("services/%s", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	svc := new(Service)
	_, err = s.client.do(ctx, req, svc)
	return svc, err
}

// Start triggers container build and startup for the service.
func (s *ServicesService) Start(ctx context.Context, uuid string) (string, error) {
	path := fmt.Sprintf("services/%s/start", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}

	var message string
	_, err = s.client.do(ctx, req, &message)
	return message, err
}

// Stop stops the service containers.
func (s *ServicesService) Stop(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("services/%s/stop", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

// Restart restarts the service containers.
func (s *ServicesService) Restart(ctx context.Context, uuid string) (string, error) {
	path := fmt.Sprintf("services/%s/restart", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}

	var message string
	_, err = s.client.do(ctx, req, &message)
	return message, err
}

// ServiceUrlMapping represents service name to URL mapping.
type ServiceUrlMapping struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// CreateServiceRequest represents the request body for creating a service.
type CreateServiceRequest struct {
	Type                          *string             `json:"type,omitempty"`
	Name                          *string             `json:"name,omitempty"`
	Description                   *string             `json:"description,omitempty"`
	ProjectUUID                   string              `json:"project_uuid"`
	EnvironmentName               *string             `json:"environment_name,omitempty"`
	EnvironmentUUID               string              `json:"environment_uuid,omitempty"`
	ServerUUID                    string              `json:"server_uuid"`
	DestinationUUID               *string             `json:"destination_uuid,omitempty"`
	InstantDeploy                 *bool               `json:"instant_deploy,omitempty"`
	DockerComposeRaw              *string             `json:"docker_compose_raw,omitempty"`
	Urls                          []ServiceUrlMapping `json:"urls,omitempty"`
	ForceDomainOverride           *bool               `json:"force_domain_override,omitempty"`
	IsContainerLabelEscapeEnabled *bool               `json:"is_container_label_escape_enabled,omitempty"`
}

// CreateServiceResponse represents the response when a service is successfully created.
type CreateServiceResponse struct {
	UUID    string   `json:"uuid"`
	Domains []string `json:"domains,omitempty"`
}

// Create creates a new service (one-click or custom docker-compose).
func (s *ServicesService) Create(ctx context.Context, reqBody CreateServiceRequest) (*CreateServiceResponse, error) {
	req, err := s.client.newRequest(http.MethodPost, "services", reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(CreateServiceResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// Delete deletes a specific service by UUID.
func (s *ServicesService) Delete(ctx context.Context, uuid string, deleteConfigs, deleteVolumes, dockerCleanup, deleteNetworks bool) (string, error) {
	path := fmt.Sprintf("services/%s?delete_configurations=%t&delete_volumes=%t&docker_cleanup=%t&delete_connected_networks=%t", uuid, deleteConfigs, deleteVolumes, dockerCleanup, deleteNetworks)
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return "", err
	}

	var resp struct {
		Message string `json:"message"`
	}
	_, err = s.client.do(ctx, req, &resp)
	return resp.Message, err
}

// CreateServiceEnvRequest represents parameters to create a new environment variable for a service.
type CreateServiceEnvRequest struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	IsPreview   bool   `json:"is_preview"`
	IsLiteral   bool   `json:"is_literal"`
	IsMultiline bool   `json:"is_multiline"`
	IsShownOnce bool   `json:"is_shown_once"`
}

// CreateServiceEnvResponse represents the response when a service env variable is created.
type CreateServiceEnvResponse struct {
	UUID string `json:"uuid"`
}

// CreateEnv creates a new environment variable for a service.
func (s *ServicesService) CreateEnv(ctx context.Context, uuid string, reqBody CreateServiceEnvRequest) (*CreateServiceEnvResponse, error) {
	path := fmt.Sprintf("services/%s/envs", uuid)
	req, err := s.client.newRequest(http.MethodPost, path, reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(CreateServiceEnvResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// ListEnvs retrieves all environment variables for the service.
func (s *ServicesService) ListEnvs(ctx context.Context, uuid string) ([]EnvironmentVariable, error) {
	path := fmt.Sprintf("services/%s/envs", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var envs []EnvironmentVariable
	_, err = s.client.do(ctx, req, &envs)
	return envs, err
}

// UpdateServiceEnvRequest represents parameters to update an existing environment variable for a service.
type UpdateServiceEnvRequest struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	IsPreview   bool   `json:"is_preview"`
	IsLiteral   bool   `json:"is_literal"`
	IsMultiline bool   `json:"is_multiline"`
	IsShownOnce bool   `json:"is_shown_once"`
}

// UpdateEnv updates an existing environment variable for a service.
// It uses PATCH /services/{uuid}/envs with key+value body.
func (s *ServicesService) UpdateEnv(ctx context.Context, uuid string, reqBody UpdateServiceEnvRequest) (*EnvironmentVariable, error) {
	path := fmt.Sprintf("services/%s/envs", uuid)
	req, err := s.client.newRequest(http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, err
	}

	env := new(EnvironmentVariable)
	_, err = s.client.do(ctx, req, env)
	return env, err
}
