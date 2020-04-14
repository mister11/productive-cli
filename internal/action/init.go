package action

import (
	"errors"
	"github.com/mister11/productive-cli/internal/utils"

	"github.com/mister11/productive-cli/internal/client"
	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/log"
	"github.com/mister11/productive-cli/internal/stdin"
)

func Init(productiveClient client.TrackingClient, promptUiStdin stdin.Stdin, configManager config.ConfigManager) {
	token := promptUiStdin.InputMasked("Enter Productive API token")

	log.Info("Saving API token...")
	configManager.SaveToken(token)

	log.Info("Fetching user organizations...")
	organizationMemberships := productiveClient.GetOrganizationMembership()

	if len(organizationMemberships) > 1 {
		utils.ReportError("Organization selection not yet supported :(",
			errors.New("organization_selection_not_supported"))
	}

	userID := organizationMemberships[0].User.ID

	configManager.SaveUserID(userID)
	log.Info("User ID saved. You can now use any CLI command available.")
}
