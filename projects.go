package coolify

import (
	"context"
	"fmt"
	"net/http"
)

// ProjectsService handles communication with the project-related endpoints of the Coolify API.
type ProjectsService struct {
	client *Client
}

// CreateProjectRequest represents the request body for creating a new project.
type CreateProjectRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// UpdateProjectRequest represents the request body for updating an existing project.
type UpdateProjectRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// CreateEnvironmentRequest represents the request body for creating a new environment in a project.
type CreateEnvironmentRequest struct {
	Name string `json:"name"`
}

// List retrieves all projects.
func (s *ProjectsService) List(ctx context.Context) ([]Project, error) {
	req, err := s.client.newRequest(http.MethodGet, "projects", nil)
	if err != nil {
		return nil, err
	}

	var projects []Project
	_, err = s.client.do(ctx, req, &projects)
	return projects, err
}

// Get retrieves a specific project by its UUID.
func (s *ProjectsService) Get(ctx context.Context, uuid string) (*Project, error) {
	path := fmt.Sprintf("projects/%s", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	project := new(Project)
	_, err = s.client.do(ctx, req, project)
	return project, err
}

// Create registers a new project.
func (s *ProjectsService) Create(ctx context.Context, reqBody CreateProjectRequest) (*Project, error) {
	req, err := s.client.newRequest(http.MethodPost, "projects", reqBody)
	if err != nil {
		return nil, err
	}

	project := new(Project)
	_, err = s.client.do(ctx, req, project)
	return project, err
}

// Update modifies configurations of an existing project.
func (s *ProjectsService) Update(ctx context.Context, uuid string, reqBody UpdateProjectRequest) (*Project, error) {
	path := fmt.Sprintf("projects/%s", uuid)
	req, err := s.client.newRequest(http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, err
	}

	project := new(Project)
	_, err = s.client.do(ctx, req, project)
	return project, err
}

// Delete removes a project by its UUID.
func (s *ProjectsService) Delete(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("projects/%s", uuid)
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

// ListEnvironments retrieves all environments for a specific project.
func (s *ProjectsService) ListEnvironments(ctx context.Context, projectUUID string) ([]Environment, error) {
	path := fmt.Sprintf("projects/%s/environments", projectUUID)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var environments []Environment
	_, err = s.client.do(ctx, req, &environments)
	return environments, err
}

// GetEnvironment retrieves a specific environment configuration by name or UUID inside a project.
func (s *ProjectsService) GetEnvironment(ctx context.Context, projectUUID string, envNameOrUUID string) (*Environment, error) {
	path := fmt.Sprintf("projects/%s/%s", projectUUID, envNameOrUUID)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	env := new(Environment)
	_, err = s.client.do(ctx, req, env)
	return env, err
}

// CreateEnvironment registers a new environment within a project.
func (s *ProjectsService) CreateEnvironment(ctx context.Context, projectUUID string, reqBody CreateEnvironmentRequest) (*Environment, error) {
	path := fmt.Sprintf("projects/%s/environments", projectUUID)
	req, err := s.client.newRequest(http.MethodPost, path, reqBody)
	if err != nil {
		return nil, err
	}

	env := new(Environment)
	_, err = s.client.do(ctx, req, env)
	return env, err
}

// DeleteEnvironment deletes a specific environment within a project.
func (s *ProjectsService) DeleteEnvironment(ctx context.Context, projectUUID string, envNameOrUUID string) error {
	path := fmt.Sprintf("projects/%s/environments/%s", projectUUID, envNameOrUUID)
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}
