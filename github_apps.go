package coolify

import (
	"context"
	"fmt"
	"net/http"
)

// GitHubAppsService handles communication with GitHub App integration endpoints.
type GitHubAppsService struct {
	client *Client
}

// GitHubRepositoriesResponse represents the response when listing repositories accessible by a GitHub App.
type GitHubRepositoriesResponse struct {
	Repositories []map[string]interface{} `json:"repositories"`
}

// GitHubBranchesResponse represents the response when listing branches for a GitHub Repository.
type GitHubBranchesResponse struct {
	Branches []map[string]interface{} `json:"branches"`
}

// List retrieves all GitHub App registrations inside the team.
func (s *GitHubAppsService) List(ctx context.Context) ([]GitHubApp, error) {
	req, err := s.client.newRequest(http.MethodGet, "github-apps", nil)
	if err != nil {
		return nil, err
	}

	var apps []GitHubApp
	_, err = s.client.do(ctx, req, &apps)
	return apps, err
}

// Get retrieves detailed registration information for a specific GitHub App.
func (s *GitHubAppsService) Get(ctx context.Context, id int) (*GitHubApp, error) {
	path := fmt.Sprintf("github-apps/%d", id)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	app := new(GitHubApp)
	_, err = s.client.do(ctx, req, app)
	return app, err
}

// Delete removes a GitHub App integration.
func (s *GitHubAppsService) Delete(ctx context.Context, id int) error {
	path := fmt.Sprintf("github-apps/%d", id)
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}

// LoadRepositories fetches all Git repositories accessible under the installations of the GitHub App.
func (s *GitHubAppsService) LoadRepositories(ctx context.Context, githubAppID int) (*GitHubRepositoriesResponse, error) {
	path := fmt.Sprintf("github-apps/%d/repositories", githubAppID)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	resp := new(GitHubRepositoriesResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}

// LoadBranches loads all branch tags for a specific GitHub repository under the GitHub App.
func (s *GitHubAppsService) LoadBranches(ctx context.Context, githubAppID int, owner string, repo string) (*GitHubBranchesResponse, error) {
	path := fmt.Sprintf("github-apps/%d/repositories/%s/%s/branches", githubAppID, owner, repo)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	resp := new(GitHubBranchesResponse)
	_, err = s.client.do(ctx, req, resp)
	return resp, err
}
