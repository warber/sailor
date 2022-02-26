package ui

// inspration: https://github.com/charmbracelet/glow/blob/b36e5ad810b6ef9ee3b32869d2af5d7ce461c2fc/ui/ui.go

import (
	tea "github.com/charmbracelet/bubbletea"
	api "github.com/keptn/go-utils/pkg/api/utils"
	authenticator "github.com/warber/sailor/pkg/auth"
	"github.com/warber/sailor/pkg/authstore/static"
	"github.com/warber/sailor/ui/common"
	"github.com/warber/sailor/ui/newproject"
	"github.com/warber/sailor/ui/projects"
	"github.com/warber/sailor/ui/sendevent"
	services2 "github.com/warber/sailor/ui/services"
	"log"
)

type state int

type newKeptnAPISetMsg api.KeptnInterface

func newKeptnAPISet(keptnAPI *api.APISet) tea.Cmd {
	return func() tea.Msg {
		return newKeptnAPISetMsg(keptnAPI)
	}
}

const (
	stateShowAuth state = iota
	stateShowProjects
	stateShowServices
	stateShowCreateProject
	stateShowSendEvent
)

func (s state) String() string {
	return map[state]string{
		stateShowAuth:          "showing auth view",
		stateShowProjects:      "showing projects view",
		stateShowServices:      "showing services view",
		stateShowCreateProject: "showing create project view",
	}[s]
}

type model struct {
	common *common.CommonModel
	state  state

	// sub-models
	projects      projects.ProjectModel
	services      services2.ServiceModel
	createProject newproject.CreateProjectModel
	sendEvent     sendevent.SendEventModel
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd
	// TODO implement proper authentication
	authenticator := authenticator.New(static.New("", ""))
	apiSet, err := authenticator.ReAuth()
	if err != nil {
		log.Fatal(err)
	}

	cmds = append(cmds, newKeptnAPISet(apiSet))

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case newKeptnAPISetMsg:
		m.common.KeptnAPI = msg
		cmds = append(cmds, common.LoadProjects(m.common.KeptnAPI.ProjectsV1()))
	case common.ServiceLoadedMsg:
		m.state = stateShowServices
	case common.ProjectCreatedMsg:
		m.state = stateShowProjects
		cmds = append(cmds, common.LoadProjects(m.common.KeptnAPI.ProjectsV1()))
	case common.ProjectDeletedMsg:
		m.state = stateShowProjects
		cmds = append(cmds, common.LoadProjects(m.common.KeptnAPI.ProjectsV1()))
	case common.ShowProjectsReqMsg:
		m.state = stateShowProjects
	case common.ShowCreateProjectMsg:
		m.state = stateShowCreateProject
		m.createProject = newproject.NewCreateProjectModel(m.common)
	case common.ShowSendEventMsg:
		m.state = stateShowSendEvent
		m.sendEvent = sendevent.NewSendEventModel(m.common)
	}

	switch m.state {
	case stateShowProjects:
		newProjectModel, cmd := m.projects.Update(msg)
		m.projects = newProjectModel
		cmds = append(cmds, cmd)
	case stateShowServices:
		newServiceModel, cmd := m.services.Update(msg)
		m.services = newServiceModel
		cmds = append(cmds, cmd)
	case stateShowCreateProject:
		newCreateProjectModel, cmd := m.createProject.Update(msg)
		m.createProject = newCreateProjectModel
		cmds = append(cmds, cmd)
	case stateShowSendEvent:
		newSendEventModel, cmd := m.sendEvent.Update(msg)
		m.sendEvent = newSendEventModel
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state {
	case stateShowServices:
		return m.services.View()
	case stateShowCreateProject:
		return m.createProject.View()
	case stateShowSendEvent:
		return m.sendEvent.View()
	default:
		return m.projects.View()
	}
}

func newModel() tea.Model {
	common := common.CommonModel{}
	return model{
		common:        &common,
		state:         stateShowProjects,
		projects:      projects.NewProjectModel(&common),
		services:      services2.NewServiceModel(&common),
		createProject: newproject.NewCreateProjectModel(&common),
	}
}

func NewProgram() *tea.Program {
	return tea.NewProgram(newModel(), tea.WithAltScreen())
}
