package productive

import (
	"fmt"
	"github.com/google/jsonapi"
)

type SessionRequest struct {
	ID       string `jsonapi:"primary,sessions"`
	Otp      string `jsonapi:"attr,otp"`
	Email    string `jsonapi:"attr,email,omitempty"`
	Password string `jsonapi:"attr,password"`
	Token    string `jsonapi:"attr,token"`
	User     *User  `jsonapi:"relation,user"`
}

type SessionResponse struct {
	ID                  string `jsonapi:"primary,sessions"`
	Token               string `jsonapi:"attr,token"`
	TokenExpirationDate string `jsonapi:"attr,token_expires_at"`
	Is2FaAuthed         bool   `jsonapi:"attr,two_factor_auth"`
	User                *User  `jsonapi:"relation,user"`
}

type User struct {
	ID           string `jsonapi:"primary,users"`
	Is2FaEnabled bool   `jsonapi:"attr,two_factor_auth"`
}

type sessionService struct {
	client *Client
}

func newSessionService(client *Client) *sessionService {
	return &sessionService{
		client: client,
	}
}

func (service *sessionService) Login(
	username string,
	password string,
) (*SessionResponse, error) {
	sessionRequest := &SessionRequest{
		ID:       "0",
		Email:    username,
		Password: password,
	}

	req, err := service.client.NewRequest("POST", "sessions", sessionRequest, getDefaultHeaders())
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

func (service *sessionService) ValidateSession(
	otp string,
	password string,
	sessionResponse *SessionResponse,
) (*SessionResponse, error) {
	sessionRequest := &SessionRequest{
		ID:       "0",
		Otp:      otp,
		Password: password,
		Token:    sessionResponse.Token,
		User:     sessionResponse.User,
	}

	uri := fmt.Sprintf("sessions/%s/validate_otp", sessionResponse.ID)
	req, err := service.client.NewRequest("PUT", uri, sessionRequest, getDefaultHeaders())

	if err != nil {
		return nil, err
	}
	sessionResponseBody, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}

	var validatedSessionsResponse SessionResponse
	if err := jsonapi.UnmarshalPayload(sessionResponseBody, &validatedSessionsResponse); err != nil {
		return nil, err
	}

	return &validatedSessionsResponse, nil
}
