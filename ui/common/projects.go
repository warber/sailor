package common

import (
	"encoding/base64"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/fileutils"
	"github.com/keptn/go-utils/pkg/common/httputils"
	"time"
)

type ProjectCreatedMsg struct{}
type ProjectCreatedErrMsg struct{ Err error }
type ShowCreateProjectMsg struct{}
type ShowProjectsReqMsg struct{}
type ProjectLoadedMsg struct{ Projects []*models.Project }
type ProjectLoadedErrMsg struct{ erro error }
type ProjectDeleteReqMsg struct{}
type ProjectDeletedMsg struct{}
type ProjectDeletedErrMsg struct{ err error }

func ShowCreateProjectView() tea.Cmd {
	return func() tea.Msg {
		return ShowCreateProjectMsg{}
	}
}

func ShowProjectsReq() tea.Cmd {
	return func() tea.Msg {
		return ShowProjectsReqMsg{}
	}
}

func LoadProjects(projectsAPI api.ProjectsV1Interface) tea.Cmd {
	return func() tea.Msg {
		allProjects, err := projectsAPI.GetAllProjects()
		if err != nil {
			return ProjectLoadedErrMsg{err}
		}
		return ProjectLoadedMsg{Projects: allProjects}
	}
}

func DeleteProject(api api.APIV1Interface, projectName string) tea.Cmd {
	return func() tea.Msg {
		_, err := api.DeleteProject(models.Project{ProjectName: projectName})
		if err != nil {
			return ProjectDeletedErrMsg{errors.New(err.GetMessage())}
		}
		return ProjectDeletedMsg{}
	}
}

func CreateProject(api api.APIV1Interface, projectName, shipyardFileLocation string) tea.Cmd {

	return func() tea.Msg {
		shipyard, err := retrieveShipyard(shipyardFileLocation)
		if err != nil {
			return ProjectCreatedErrMsg{err}
		}

		encodedShipyardContent := base64.StdEncoding.EncodeToString(shipyard)
		project := models.CreateProject{
			Name:     &projectName,
			Shipyard: &encodedShipyardContent,
		}

		_, err2 := api.CreateProject(project)
		if err2 != nil {
			return ProjectCreatedErrMsg{Err: errors.New(err2.GetMessage())}
		}
		return ProjectCreatedMsg{}
	}
}

func retrieveShipyard(location string) ([]byte, error) {
	var content []byte
	var err error
	if httputils.IsValidURL(location) {
		content, err = httputils.NewDownloader(httputils.WithTimeout(5 * time.Second)).DownloadFromURL(location)
	} else {
		content, err = fileutils.ReadFile(location)
	}
	if err != nil {
		return nil, err
	}
	return content, nil
}
