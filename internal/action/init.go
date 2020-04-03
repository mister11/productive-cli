package action

import (
	"gitlab.com/mister11/productive-cli/internal/utils"

	"gitlab.com/mister11/productive-cli/internal/client"
	"gitlab.com/mister11/productive-cli/internal/config"
	"gitlab.com/mister11/productive-cli/internal/log"
	"gitlab.com/mister11/productive-cli/internal/prompt"
)

func Init(client client.ProductiveClient) {
	token := prompt.InputMasked("Enter Productive API token")

	log.Info("Saving API token...")
	config.SaveToken(token)

	log.Info("Fetching user organizations...")
	organizationMemberships := client.GetOrganizationMembership()

	if len(organizationMemberships) > 1 {
		utils.ReportError("Organization selection not yet supported :(", nil)
	}

	log.Debug(organizationMemberships)
	userID := organizationMemberships[0].User.ID

	config.SaveUserID(userID)
	log.Info("User ID saved. You can now use any CLI command available.")
}
