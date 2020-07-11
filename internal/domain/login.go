package domain

type LoginData struct {
	Token               string
	TokenExpirationDate string
}

type OrganizationMembershipData struct {
	PersonID string
}

type LoginManager interface {
	IsSessionValid() (bool, error)
	Login(username string, password string) error
}