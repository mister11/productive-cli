package action

import (
	"errors"
	"github.com/mister11/productive-cli/internal/utils"

	"github.com/mister11/productive-cli/internal/client"
	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/log"
	"github.com/mister11/productive-cli/internal/stdin"
)

type LoginManger struct {
	trackingClient client.TrackingClient
	stdIn          stdin.Stdin
	config         config.ConfigManager
}

func NewLoginManager(
	trackingClient client.TrackingClient,
	stdIn stdin.Stdin,
	config config.ConfigManager,
) *LoginManger {
	return &LoginManger{
		trackingClient: trackingClient,
		stdIn:          stdIn,
		config:         config,
	}
}

func (manager *LoginManger) Init() {
	token := manager.stdIn.InputMasked("Enter Productive API token")

	log.Info("Saving API token...")
	manager.config.SaveToken(token)

	log.Info("Fetching user organizations...")
	organizationMemberships := manager.trackingClient.GetOrganizationMembership()

	if len(organizationMemberships) > 1 {
		utils.ReportError("Organization selection not yet supported :(",
			errors.New("organization_selection_not_supported"))
	}

	userID := organizationMemberships[0].User.ID

	manager.config.SaveUserID(userID)
	log.Info("User ID saved. You can now use any CLI command available.")
}
