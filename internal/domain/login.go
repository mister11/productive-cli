package domain

type LoginData struct {
	Token               string
	TokenExpirationDate string
}

type OrganizationMembershipData struct {
	PersonID string
}

type SessionData struct {
	PersonID string
	Token    string
}
