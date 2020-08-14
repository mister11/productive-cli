package service

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

type ProjectStorage interface {
	UpsertTrackedProject(project TrackedProject) error
	GetTrackedProjects() ([]TrackedProject, error)
}

type FileProjectStorage struct{}

func NewFileProjectStorage() *FileProjectStorage {
	return &FileProjectStorage{}
}

func (f FileProjectStorage) UpsertTrackedProject(project TrackedProject) error {
	projects, err := f.GetTrackedProjects()
	if err != nil {
		return err
	}
	if projectExists(projects, project) {
		return nil
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
	if err := ensureConfigFolderExists(); err != nil {
		return err
	}
	return utils.WriteFile(*projectsConfigPath, projectsConfigJSON)
}

func (f FileProjectStorage) GetTrackedProjects() ([]TrackedProject, error) {
	configPath, err := getProjectConfigPath()
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(*configPath)
	if os.IsNotExist(err) {
		return nil, nil
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

func projectExists(
	projects []TrackedProject,
	project TrackedProject,
) bool {
	for _, savedProject := range projects {
		if savedProject.DealID == project.DealID && savedProject.ServiceID == project.ServiceID {
			return true
		}
	}
	return false
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