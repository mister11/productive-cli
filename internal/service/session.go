package service

import (
	"errors"
	"time"

	"github.com/mister11/productive-cli/internal/log"
	"github.com/mister11/productive-cli/internal/productive"
)

type SessionService struct {
	productiveService *ProductiveService
	prompt            Prompt
	sessionManager    UserSessionManager
}

func NewSessionService(
	productiveService *ProductiveService,
	prompt Prompt,
	sessionManager UserSessionManager,
) *SessionService {
	return &SessionService{
		productiveService: productiveService,
		prompt:            prompt,
		sessionManager:    sessionManager,
	}
}

func (s *SessionService) ObtainUserSession() (*UserSessionData, error) {
	verifiedSession, err := s.verifySession()
	if err != nil {
		return nil, err
	}
	if verifiedSession != nil {
		return verifiedSession, nil
	}
	username, err := s.prompt.Input("E-mail")
	if err != nil {
		return nil, err
	}
	password, err := s.prompt.InputMasked("Password")
	if err != nil {
		return nil, err
	}
	return s.login(username, password)
}

func (s *SessionService) verifySession() (*UserSessionData, error) {
	userSession, err := s.sessionManager.GetUserSession()
	if err != nil {
		return nil, err
	}
	if userSession == nil {
		return nil, nil
	}
	tokenExpirationDate, err := time.Parse(time.RFC3339, userSession.TokenExpirationDate)
	if err != nil {
		return nil, err
	}
	if tokenExpirationDate.Before(time.Now()) {
		return nil, nil
	}
	return userSession, nil
}

func (s *SessionService) login(username string, password string) (*UserSessionData, error) {
	sessionResponse, err := s.productiveService.Login(username, password)
	if err != nil {
		return nil, err
	}
	if sessionResponse.User.Is2FaEnabled {
		otp, err := s.prompt.Input("Enter OTP from authentication application")
		if err != nil {
			return nil, err
		}
		sessionResponse, err = s.productiveService.ValidateSession(otp, password, sessionResponse)
		if err != nil {
			return nil, err
		}
		if !sessionResponse.Is2FaAuthed {
			return nil, errors.New("error with 2FA code")
		}
	}
	currentUserSession, err := s.sessionManager.GetUserSession()
	if err != nil {
		return nil, err
	}
	var personID string
	if isFirstLogin(currentUserSession) {
		personID, err = s.handleFirstLogin(*sessionResponse)
		if err != nil {
			return nil, err
		}
	} else {
		personID = currentUserSession.PersonID
	}
	userSession := &UserSessionData{
		Token:               sessionResponse.Token,
		TokenExpirationDate: sessionResponse.TokenExpirationDate,
		PersonID:            personID,
	}
	if err := s.sessionManager.SaveUserSession(*userSession); err != nil {
		return nil, err
	}
	return userSession, nil
}

func isFirstLogin(currentUserSession *UserSessionData) bool {
	return currentUserSession == nil || currentUserSession.PersonID == ""
}

func (s *SessionService) handleFirstLogin(response productive.SessionResponse) (string, error) {
	log.Info("First login. Setting up necessary tracking data.")
	userSession := &UserSessionData{
		Token:               response.Token,
		TokenExpirationDate: response.TokenExpirationDate,
	}
	err := s.sessionManager.SaveUserSession(*userSession)
	if err != nil {
		return "", err
	}
	organizationMemberships, err := s.productiveService.GetOrganizationMemberships(userSession)
	if err != nil {
		return "", err
	}
	if len(organizationMemberships) == 0 {
		return "", errors.New("no organization memberships found")
	}
	if len(organizationMemberships) > 1 {
		log.Info("Multiple organizations found. This is not currently supported. Taking the first one.")
	}
	return organizationMemberships[0].User.ID, nil
}
