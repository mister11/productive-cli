package productive

import (
	"github.com/google/jsonapi"
	"reflect"
)

type OrganizationMembershipResponse struct {
	PersonID string
}

type organizationMembershipService struct {
	client *client
}

func newOrganizationMembershipService(client *client) *organizationMembershipService {
	return &organizationMembershipService{
		client: client,
	}
}

func (service *organizationMembershipService) FetchAll(
	headers map[string]string,
) ([]OrganizationMembershipResponse, error) {

	req, err := service.client.NewRequest("GET", "organization_memberships", nil, headers)

	if err != nil {
		return nil, err
	}

	organizationMembershipsBody, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}

	organizationMembershipInterfaces, err := jsonapi.UnmarshalManyPayload(
		organizationMembershipsBody,
		reflect.TypeOf(new(OrganizationMembershipResponse)),
	)
	if err != nil {
		return nil, err
	}

	var organizationMemberships []OrganizationMembershipResponse
	for _, organizationMembershipInterface := range organizationMembershipInterfaces {
		organizationMemberships = append(organizationMemberships, organizationMembershipInterface.(OrganizationMembershipResponse))
	}
	return organizationMemberships, nil
}
