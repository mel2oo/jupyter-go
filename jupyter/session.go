package jupyter

import (
	"context"
	"net/http"
)

type SessionService service

// Get returns a session by ID.
func (s *SessionService) Get(ctx context.Context, id string) (*Session, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/api/sessions/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp *Session
	if err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// List returns a list of all sessions.
func (s *SessionService) List(ctx context.Context) ([]Session, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/api/sessions", nil)
	if err != nil {
		return nil, err
	}

	var resp []Session
	if err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Create creates a new session.
func (s *SessionService) Create(ctx context.Context, body *Session) (*Session, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/api/sessions", body)
	if err != nil {
		return nil, err
	}

	var resp *Session
	if err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates a session by ID.
func (s *SessionService) Update(ctx context.Context, id string, body *Session) (*Session, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPatch, "/api/sessions/"+id, body)
	if err != nil {
		return nil, err
	}

	var resp *Session
	if err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a session by ID.
func (s *SessionService) Delete(ctx context.Context, id string) error {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, "/api/sessions/"+id, nil)
	if err != nil {
		return err
	}

	return s.client.Do(req, nil)
}
