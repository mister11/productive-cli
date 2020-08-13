package productive

import (
	"github.com/google/jsonapi"
	"reflect"
)

type OrganizationMembership struct {
	ID   string  `jsonapi:"primary,organization_memberships"`
	User *Person `jsonapi:"relation,person"`
}

type organizationMembershipService struct {
	client *Client
}

func newOrganizationMembershipService(client *Client) *organizationMembershipService {
	return &organizationMembershipService{
		client: client,
	}
}

func (service *organizationMembershipService) FetchAll(
	token string,
) ([]OrganizationMembership, error) {

	req, err := service.client.NewRequest("GET", "organization_memberships", nil, getHeaders(token))

	if err != nil {
		return nil, err
	}

	organizationMembershipsBody, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}

	organizationMembershipInterfaces, err := jsonapi.UnmarshalManyPayload(
		organizationMembershipsBody,
		reflect.TypeOf(new(OrganizationMembership)),
	)
	if err != nil {
		return nil, err
	}

	var organizationMemberships []OrganizationMembership
	for _, organizationMembershipInterface := range organizationMembershipInterfaces {
		orgMembership := organizationMembershipInterface.(*OrganizationMembership)
		organizationMemberships = append(organizationMemberships, *orgMembership)
	}
	return organizationMemberships, nil
}
