package action

import (
	"github.com/mister11/productive-cli/internal/client/model"
	"github.com/mister11/productive-cli/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestInitMultipleOrgs(t *testing.T) {
	client := new(mocks.TrackingClient)
	stdIn := new(mocks.Stdin)
	configManager := new(mocks.ConfigManager)

	testToken := "test_token"
	stdIn.
		On("InputMasked", mock.Anything).
		Return(testToken).Once()

	configManager.
		On("SaveToken", testToken).
		Return().Once()

	orgMembership := model.OrganizationMembership{
		ID:   "100",
		User: &model.Person{ID: "101"},
	}
	client.
		On("GetOrganizationMembership").
		Return([]model.OrganizationMembership{orgMembership, orgMembership}).Once()

	assert.Panics(t, func() { Init(client, stdIn, configManager) })
	client.AssertExpectations(t)
	stdIn.AssertExpectations(t)
	configManager.AssertExpectations(t)
}

func TestInit(t *testing.T) {
	client := new(mocks.TrackingClient)
	stdIn := new(mocks.Stdin)
	configManager := new(mocks.ConfigManager)

	testToken := "test_token"
	stdIn.
		On("InputMasked", mock.Anything).
		Return(testToken).Once()

	configManager.
		On("SaveToken", testToken).
		Return().Once()

	orgMembership := model.OrganizationMembership{
		ID:   "100",
		User: &model.Person{ID: "101"},
	}
	client.
		On("GetOrganizationMembership").
		Return([]model.OrganizationMembership{orgMembership}).Once()

	configManager.
		On("SaveUserID", "101").
		Return().Once()

	Init(client, stdIn, configManager)

	client.AssertExpectations(t)
	stdIn.AssertExpectations(t)
	configManager.AssertExpectations(t)
}
