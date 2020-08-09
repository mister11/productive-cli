package productive

import (
	"github.com/google/jsonapi"
)

type SessionRequest struct {
	ID       string `jsonapi:"primary,sessions"`
	Email    string `jsonapi:"attr,email"`
	Password string `jsonapi:"attr,password"`
}

type SessionResponse struct {
	ID                  string `jsonapi:"primary,sessions"`
	Token               string `jsonapi:"attr,token"`
	TokenExpirationDate string `jsonapi:"attr,token_expires_at"`
}

type sessionService struct {
	client *client
}

func newSessionService(client *client) *sessionService {
	return &sessionService{
		client: client,
	}
}

func (service *sessionService) Login(
	username string,
	password string,
	headers map[string]string,
) (*SessionResponse, error) {
	sessionRequest := &SessionRequest{
		ID:       "0",
		Email:    username,
		Password: password,
	}

	req, err := service.client.NewRequest("POST", "sessions", sessionRequest, headers)
	if err != nil {
		return nil, err
	}
	sessionResponseBody, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}
	var sessionResponse SessionResponse
	if err := jsonapi.UnmarshalPayload(sessionResponseBody, &sessionResponse); err != nil {
		return nil, err
	}
	return &sessionResponse, nil
}
