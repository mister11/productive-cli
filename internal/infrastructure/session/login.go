package session

import (
	"errors"
	"github.com/mister11/productive-cli/internal/domain/tracking"
	"github.com/mister11/productive-cli/internal/infrastructure/client"
	"github.com/mister11/productive-cli/internal/infrastructure/log"
	"time"
)

type httpLoginManager struct {
	trackingClient tracking.TrackingClient
	sessionManager client.UserSessionManager
}

func NewProductiveLoginManager(
	trackingClient tracking.TrackingClient,
	sessionManager client.UserSessionManager,
) httpLoginManager {
	return httpLoginManager{
		trackingClient: trackingClient,
		sessionManager: sessionManager,
	}
}

func (p httpLoginManager) IsSessionValid() (bool, error) {
	userConfig, err := p.sessionManager.GetUserSession()
	if err != nil {
		return false, err
	}
	tokenExpirationDate, err := time.Parse(time.RFC3339, userConfig.TokenExpirationDate)
	if err != nil {
		return false, err
	}
	// not an error, but session data is expired
	if tokenExpirationDate.Before(time.Now()) {
		return true, nil
	}
	return false, nil
}

func (p httpLoginManager) Login(username string, password string) error {
	loginData, err := p.trackingClient.Login(username, password)
	if err != nil {
		return err
	}
	userConfig, err := p.sessionManager.GetUserSession()
	if err != nil {
		return err
	}
	if userConfig.PersonID != "" {
		return nil
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
	userConfig.PersonID = personID
	userConfig.Token = loginData.Token
	userConfig.TokenExpirationDate = loginData.TokenExpirationDate
	if err := p.sessionManager.SaveUserSession(*userConfig); err != nil {
		return err
	}
	log.Info("Assigned organization to the user. Session created. Expiration date %s", userConfig.TokenExpirationDate)
	return nil
}
