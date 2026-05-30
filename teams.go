package coolify

import (
	"context"
	"fmt"
	"net/http"
)

// TeamsService handles communication with team and user membership endpoints.
type TeamsService struct {
	client *Client
}

// List retrieves all teams.
func (s *TeamsService) List(ctx context.Context) ([]Team, error) {
	req, err := s.client.newRequest(http.MethodGet, "teams", nil)
	if err != nil {
		return nil, err
	}

	var teams []Team
	_, err = s.client.do(ctx, req, &teams)
	return teams, err
}

// Get retrieves detailed configurations of a team by its ID.
func (s *TeamsService) Get(ctx context.Context, id int) (*Team, error) {
	path := fmt.Sprintf("teams/%d", id)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	team := new(Team)
	_, err = s.client.do(ctx, req, team)
	return team, err
}

// Members lists all user members belonging to the specified team.
func (s *TeamsService) Members(ctx context.Context, id int) ([]User, error) {
	path := fmt.Sprintf("teams/%d/members", id)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var users []User
	_, err = s.client.do(ctx, req, &users)
	return users, err
}

// Current retrieves details of the currently authenticated team.
func (s *TeamsService) Current(ctx context.Context) (*Team, error) {
	req, err := s.client.newRequest(http.MethodGet, "teams/current", nil)
	if err != nil {
		return nil, err
	}

	team := new(Team)
	_, err = s.client.do(ctx, req, team)
	return team, err
}

// CurrentMembers retrieves members list of the currently authenticated team.
func (s *TeamsService) CurrentMembers(ctx context.Context) ([]User, error) {
	req, err := s.client.newRequest(http.MethodGet, "teams/current/members", nil)
	if err != nil {
		return nil, err
	}

	var users []User
	_, err = s.client.do(ctx, req, &users)
	return users, err
}
