package coolify

import (
	"context"
	"fmt"
	"net/http"
)

// DeploymentsService handles communication with deployment-related endpoints of the Coolify API.
type DeploymentsService struct {
	client *Client
}

// DeployRequest represents the parameters for deploying an application by its UUID or tag.
type DeployRequest struct {
	ApplicationUUID string  `json:"uuid"`
	Tag             *string `json:"tag,omitempty"`
}

// DeployResponse represents the response when a deployment is triggered.
type DeployResponse struct {
	Message string `json:"message"`
}

// List retrieves a list of all active or queued deployments across all applications.
func (s *DeploymentsService) List(ctx context.Context) ([]ApplicationDeploymentQueue, error) {
	req, err := s.client.newRequest(http.MethodGet, "deployments", nil)
	if err != nil {
		return nil, err
	}

	var queues []ApplicationDeploymentQueue
	_, err = s.client.do(ctx, req, &queues)
	return queues, err
}

// Get retrieves detailed status and live logs for a specific deployment run.
func (s *DeploymentsService) Get(ctx context.Context, uuid string) (*ApplicationDeploymentQueue, error) {
	path := fmt.Sprintf("deployments/%s", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	queue := new(ApplicationDeploymentQueue)
	_, err = s.client.do(ctx, req, queue)
	return queue, err
}

// Cancel terminates a running or queued deployment.
func (s *DeploymentsService) Cancel(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("deployments/%s/cancel", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

// Deploy triggers a new manual deployment for an application.
// If force is true, the deployment bypasses cache and forces rebuild (query param `force=true`).
func (s *DeploymentsService) Deploy(ctx context.Context, force bool, applicationUUID string, tag *string) (*DeployResponse, error) {
	path := "deploy"
	if force {
		path = "deploy?force=true"
	}

	reqBody := DeployRequest{
		ApplicationUUID: applicationUUID,
		Tag:             tag,
	}

	req, err := s.client.newRequest(http.MethodPost, path, reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(DeployResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// ListApplicationDeployments retrieves all past deployment runs specifically for one application.
func (s *DeploymentsService) ListApplicationDeployments(ctx context.Context, applicationUUID string) ([]ApplicationDeploymentQueue, error) {
	path := fmt.Sprintf("deployments/applications/%s", applicationUUID)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var queues []ApplicationDeploymentQueue
	_, err = s.client.do(ctx, req, &queues)
	return queues, err
}
