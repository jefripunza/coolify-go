package coolify

import (
	"context"
	"fmt"
	"net/http"
)

// ServersService handles communication with server-related endpoints of the Coolify API.
type ServersService struct {
	client *Client
}

// CreateHetznerServerRequest represents parameters to create a new server on Hetzner Cloud.
type CreateHetznerServerRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	IP             string `json:"ip"`
	Port           int    `json:"port"`
	User           string `json:"user"`
	PrivateKeyUUID string `json:"private_key_uuid"`
	ProjectUUID    string `json:"project_uuid"`
}

// CreateHetzner VPS configures a new server on Hetzner Cloud and registers it in Coolify.
func (s *ServersService) CreateHetzner(ctx context.Context, reqBody CreateHetznerServerRequest) (*Server, error) {
	req, err := s.client.newRequest(http.MethodPost, "servers/hetzner", reqBody)
	if err != nil {
		return nil, err
	}

	server := new(Server)
	_, err = s.client.do(ctx, req, server)
	return server, err
}

// List retrieves all servers registered under the team.
func (s *ServersService) List(ctx context.Context) ([]Server, error) {
	req, err := s.client.newRequest(http.MethodGet, "servers", nil)
	if err != nil {
		return nil, err
	}

	var servers []Server
	_, err = s.client.do(ctx, req, &servers)
	return servers, err
}

// Get retrieves detailed server configurations.
func (s *ServersService) Get(ctx context.Context, uuid string) (*Server, error) {
	path := fmt.Sprintf("servers/%s", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	server := new(Server)
	_, err = s.client.do(ctx, req, server)
	return server, err
}

// Delete removes a server from Coolify.
func (s *ServersService) Delete(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("servers/%s", uuid)
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

// ServerValidationResponse represents validation feedback of a server.
type ServerValidationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Validate tests SSH and Docker connectivity for the server.
func (s *ServersService) Validate(ctx context.Context, uuid string) (*ServerValidationResponse, error) {
	path := fmt.Sprintf("servers/%s/validate", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	valResp := new(ServerValidationResponse)
	_, err = s.client.do(ctx, req, valResp)
	return valResp, err
}

// ServerResource represent resources allocated on the server.
type ServerResource struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Type string `json:"type"` // e.g. application, database, service
}

// Resources lists all active containers and resources deployed on a server.
func (s *ServersService) Resources(ctx context.Context, uuid string) ([]ServerResource, error) {
	path := fmt.Sprintf("servers/%s/resources", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resources []ServerResource
	_, err = s.client.do(ctx, req, &resources)
	return resources, err
}

// ServerDomain represents domain configurations mapped to a server.
type ServerDomain struct {
	FQDN         string `json:"fqdn"`
	ResourceName string `json:"resource_name"`
	ResourceUUID string `json:"resource_uuid"`
	ResourceType string `json:"resource_type"`
}

// Domains retrieves list of domains routing through Traefik/Caddy on a server.
func (s *ServersService) Domains(ctx context.Context, uuid string) ([]ServerDomain, error) {
	path := fmt.Sprintf("servers/%s/domains", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var domains []ServerDomain
	_, err = s.client.do(ctx, req, &domains)
	return domains, err
}
