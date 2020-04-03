package config

import (
	"encoding/json"
	"gitlab.com/mister11/productive-cli/internal/client/model"
	"os"
	"path/filepath"

	"gitlab.com/mister11/productive-cli/internal/utils"
)

const configFolder = ".productive"
const configFilename = "config"

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

func GetToken() string {
	return loadConfig().UserToken
}

func GetUserID() string {
	return loadConfig().UserID
}

func SaveToken(token string) {
	config := loadConfig()

	config.UserToken = token

	saveConfig(config)
}

func SaveUserID(userID string) {
	config := loadConfig()

	config.UserID = userID

	saveConfig(config)
}

func SaveProjectToConfig(project Project) {
	config := loadConfig()

	config.Projects = append(config.Projects, project)

	saveConfig(config)
}

func GetSavedProjects() []Project {
	return loadConfig().Projects
}

func RemoveExistingProject(project Project) {
	config := loadConfig()

	var projects []Project
	for _, savedProject := range config.Projects {
		if !(savedProject.DealID == project.DealID && savedProject.ServiceID == project.ServiceID) {
			projects = append(projects, savedProject)
		}
	}
	config.Projects = projects
	saveConfig(config)
}

func loadConfig() Config {
	configPath := getConfigPath()
	configJSON, err := utils.ReadFile(configPath)

	if err != nil {
		return Config{}
	}

	var config Config
	if err := json.Unmarshal(configJSON, &config); err != nil {
		utils.ReportError("Error parsing config JSON", err)
	}
	return config
}

func saveConfig(config Config) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		utils.ReportError("Cannot convert config "+string(configJSON)+" to JSON", err)
	}
	utils.WriteFile(getConfigPath(), configJSON)
}

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		utils.ReportError("Error retrieving home directory", err)
	}
	return homeDir + getSeparator() + configFolder + getSeparator() + configFilename
}

func getSeparator() string {
	return string(filepath.Separator)
}
