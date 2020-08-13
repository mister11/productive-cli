package service

import (
	"errors"
	"github.com/mister11/productive-cli/internal/log"
	"github.com/mister11/productive-cli/internal/productive"
	"time"
)

type SessionService struct {
	client         *productive.Client
	sessionManager UserSessionManager
}

func NewSessionService(
	client *productive.Client,
	sessionManager UserSessionManager,
) *SessionService {
	return &SessionService{
		client:         client,
		sessionManager: sessionManager,
	}
}

func (service *SessionService) GetUserId() (string, error) {
	userSessionData, err := service.sessionManager.GetUserSession()
	if err != nil {
		return "", err
	}
	return userSessionData.PersonID, nil
}

func (service *SessionService) IsSessionValid() (bool, error) {
	userSession, err := service.sessionManager.GetUserSession()
	if err != nil {
		return false, err
	}
	if userSession == nil {
		return false, nil
	}
	tokenExpirationDate, err := time.Parse(time.RFC3339, userSession.TokenExpirationDate)
	if err != nil {
		return false, err
	}
	if tokenExpirationDate.Before(time.Now()) {
		return false, nil
	}
	return true, nil
}

func (service *SessionService) Login(username string, password string) error {
	sessionResponse, err := service.client.SessionService.Login(username, password)
	if err != nil {
		return err
	}
	currentUserSession, err := service.sessionManager.GetUserSession()
	if err != nil {
		return err
	}
	var personID string
	if isFirstLogin(currentUserSession) {
		personID, err = service.handleFirstLogin(*sessionResponse)
		if err != nil {
			return err
		}
	} else {
		personID = currentUserSession.PersonID
	}
	userSession := &UserSessionData{
		Token:               sessionResponse.Token,
		TokenExpirationDate: sessionResponse.TokenExpirationDate,
		PersonID:            personID,
	}
	if err := service.sessionManager.SaveUserSession(*userSession); err != nil {
		return err
	}
	return nil
}

func isFirstLogin(currentUserSession *UserSessionData) bool {
	return currentUserSession == nil || currentUserSession.PersonID == ""
}

func (service *SessionService) handleFirstLogin(response productive.SessionResponse) (string, error) {
	log.Info("First login. Setting up necessary tracking data.")
	userSession := &UserSessionData{
		Token: response.Token,
		TokenExpirationDate: response.TokenExpirationDate,
	}
	err := service.sessionManager.SaveUserSession(*userSession)
	if err != nil {
		return "", err
	}
	organizationMemberships, err := service.client.OrganizationMembershipService.FetchAll(response.Token)
	if len(organizationMemberships) == 0 {
		return "", errors.New("no organization memberships found")
	}
	if len(organizationMemberships) > 1 {
		log.Info("Multiple organizations found. This is not currently supported. Taking the first one.")
	}
	return organizationMemberships[0].User.ID, nil
}
