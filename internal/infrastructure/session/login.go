package session

import (
	"errors"
	"github.com/mister11/productive-cli/internal/domain/tracking"
	"github.com/mister11/productive-cli/internal/infrastructure/client"
	"github.com/mister11/productive-cli/internal/infrastructure/log"
	"time"
)

type HttpLoginManager struct {
	trackingClient tracking.TrackingClient
	sessionManager client.UserSessionManager
}

func NewProductiveLoginManager(
	trackingClient tracking.TrackingClient,
	sessionManager client.UserSessionManager,
) HttpLoginManager {
	return HttpLoginManager{
		trackingClient: trackingClient,
		sessionManager: sessionManager,
	}
}

func (p HttpLoginManager) IsSessionValid() (bool, error) {
	userSession, err := p.sessionManager.GetUserSession()
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
	if tokenExpirationDate.After(time.Now()) {
		return true, nil
	}
	return false, nil
}

func (p HttpLoginManager) Login(username string, password string) error {
	loginData, err := p.trackingClient.Login(username, password)
	if err != nil {
		return err
	}
	userSession, err := p.sessionManager.GetUserSession()
	if err != nil {
		return err
	}
	if userSession != nil {
		return nil
	}
	userSession = &client.UserSessionData{
		Token: loginData.Token,
		TokenExpirationDate: loginData.TokenExpirationDate,
	}
	if err := p.sessionManager.SaveUserSession(*userSession); err != nil {
		return err
	}
	log.Info("First login. Setting up necessary tracking data.")
	organizationMemberships, err := p.trackingClient.GetOrganizationMemberships()
	if len(organizationMemberships) == 0 {
		return errors.New("no organization memberships found")
	}
	if len(organizationMemberships) > 1 {
		log.Info("Multiple organizations found. This is not currently supported. Taking the first one.")
	}
	personID := organizationMemberships[0].PersonID
	userSession = &client.UserSessionData{
		Token: loginData.Token,
		TokenExpirationDate: loginData.TokenExpirationDate,
		PersonID: personID,
	}
	if err := p.sessionManager.SaveUserSession(*userSession); err != nil {
		return err
	}
	return nil
}
