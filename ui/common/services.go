package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
)

type ServiceLoadedErrMsg struct{ err error }
type ServiceLoadedMsg struct{ Services []*models.Service }

func LoadServices(servicesAPI api.ServicesV1Interface, project *models.Project) tea.Cmd {
	return func() tea.Msg {
		var svcs []*models.Service
		for _, stage := range project.Stages {
			services, err := servicesAPI.GetAllServices(project.ProjectName, stage.StageName)
			if err != nil {
				return ServiceLoadedErrMsg{err}
			}
			for _, svc := range services {
				svcs = append(svcs, svc)
			}
		}
		return ServiceLoadedMsg{Services: svcs}
	}
}
