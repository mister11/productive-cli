package config

import "github.com/mister11/productive-cli/internal/infrastructure/client/model"

type Config struct {
	UserID    string    `json:"user_id,omitempty"`
	UserToken string    `json:"user_token,omitempty"`
	Projects  []Project `json:"projects,omitempty"`
}

type Project struct {
	DealID      string
	DealName    string
	ServiceID   string
	ServiceName string
}

func NewProject(deal model.Deal, service model.Service) Project {
	return Project{
		DealID:      deal.ID,
		DealName:    deal.Name,
		ServiceID:   service.ID,
		ServiceName: service.Name,
	}
}

type Manager interface {
	GetToken() string
	GetUserID() string
	SaveToken(token string)
	SaveUserID(userID string)
	SaveProject(project Project)
	GetSavedProjects() []Project
	RemoveExistingProject(project Project)
}