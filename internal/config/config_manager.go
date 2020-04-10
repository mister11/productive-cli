package config

type ConfigManager interface {
	GetToken() string
	GetUserID() string
	SaveToken(token string)
	SaveUserID(userID string)
	SaveProject(project Project)
	GetSavedProjects() []Project
	RemoveExistingProject(project Project)
}