package domain

import (
	"encoding/json"
	"github.com/mister11/productive-cli/internal/utils"
	"os"
	"path/filepath"
)

const dataFolder = ".productive"
const trackedProjectsFile = "projects"

type TrackedProjects struct {
	Projects []TrackedProject `json:"projects,omitempty"`
}

type TrackedProject struct {
	DealID      string
	DealName    string
	ServiceID   string
	ServiceName string
}

type TrackedProjectManager interface {
	SaveTrackedProject(project TrackedProject) error
	GetTrackedProjects() ([]TrackedProject, error)
	RemoveTrackedProject(project TrackedProject) error
}

type FileTrackedProjectsManager struct{}

func NewFileTrackedProjectsManager() *FileTrackedProjectsManager {
	return &FileTrackedProjectsManager{}
}

func (f FileTrackedProjectsManager) SaveTrackedProject(project TrackedProject) error {
	projects, err := f.GetTrackedProjects()
	if err != nil {
		return err
	}
	projects = append(projects, project)
	projectsConfigJSON, err := json.Marshal(TrackedProjects{
		Projects: projects,
	})
	if err != nil {
		return err
	}
	projectsConfigPath, err := getProjectConfigPath()
	if err != nil {
		return err
	}
	return utils.WriteFile(*projectsConfigPath, projectsConfigJSON)
}

func (f FileTrackedProjectsManager) GetTrackedProjects() ([]TrackedProject, error) {
	configPath, err := getProjectConfigPath()
	if err != nil {
		return nil, err
	}
	projectsJSON, err := utils.ReadFile(*configPath)

	if err != nil {
		return nil, err
	}

	var projectsConfig TrackedProjects
	if err := json.Unmarshal(projectsJSON, &projectsConfig); err != nil {
		return nil, err
	}
	return projectsConfig.Projects, nil
}

func (f FileTrackedProjectsManager) RemoveTrackedProject(project TrackedProject) error {
	savedProjects, err := f.GetTrackedProjects()
	if err != nil {
		return err
	}

	var newProjects []TrackedProject
	for _, savedProject := range savedProjects {
		if !(savedProject.DealID == project.DealID && savedProject.ServiceID == project.ServiceID) {
			newProjects = append(newProjects, savedProject)
		}
	}
	projectsConfigJSON, err := json.Marshal(TrackedProjects{
		Projects: newProjects,
	})
	if err != nil {
		return err
	}
	projectsConfigPath, err := getProjectConfigPath()
	if err != nil {
		return err
	}
	return utils.WriteFile(*projectsConfigPath, projectsConfigJSON)
}

func getProjectConfigPath() (*string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := homeDir + getSeparator() + dataFolder + getSeparator() + trackedProjectsFile
	return &configPath, nil
}

func getSeparator() string {
	return string(filepath.Separator)
}