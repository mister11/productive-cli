package client

import (
	"encoding/json"
	"github.com/mister11/productive-cli/internal/utils"
	"os"
	"path/filepath"
)

const dataFolder = ".productive"
const userSessionFile = "user"

type UserSessionData struct {
	PersonID            string `json:"person_id,omitempty"`
	Token               string `json:"session_token,omitempty"`
	TokenExpirationDate string `json:"token_expiration_date,omitempty"`
}

type UserSessionManager interface {
	GetUserSession() (*UserSessionData, error)
	SaveUserSession(session UserSessionData) error
}

type FileUserSessionManager struct{}

func NewFileUserSessionManager() *FileUserSessionManager {
	return &FileUserSessionManager{}
}

func (f FileUserSessionManager) GetUserSession() (*UserSessionData, error) {
	sessionPath, err := getUserSessionPath()
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(*sessionPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	sessionJSON, err := utils.ReadFile(*sessionPath)
	if err != nil {
		return nil, err
	}

	var userSession UserSessionData
	if err := json.Unmarshal(sessionJSON, &userSession); err != nil {
		return nil, err
	}
	return &userSession, nil
}

func (f FileUserSessionManager) SaveUserSession(session UserSessionData) error {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}
	sessionPath, err := getUserSessionPath()
	if err != nil {
		return err
	}
	return utils.WriteFile(*sessionPath, sessionJSON)
}

func getUserSessionPath() (*string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	sessionPath := homeDir + getSeparator() + dataFolder + getSeparator() + userSessionFile
	return &sessionPath, nil
}

func getSeparator() string {
	return string(filepath.Separator)
}
