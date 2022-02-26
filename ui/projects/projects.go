package projects

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/warber/sailor/ui/common"
	"io"
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}
	fmt.Fprintf(w, fn(str))
}

type ProjectModel struct {
	common   *common.CommonModel
	projects []*models.Project
	list     list.Model
	choice   string
	quitting bool
}

func NewProjectModel(common *common.CommonModel) ProjectModel {
	items := []list.Item{}
	l := list.New(items, itemDelegate{}, 50, 14)
	l.Title = "Select a Keptn Project"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		b := []key.Binding{
			key.NewBinding(key.WithKeys("a", "+"),
				key.WithHelp("a/+", "add project")),
			key.NewBinding(key.WithKeys("d", "-"),
				key.WithHelp("d", "delete project")),
			key.NewBinding(key.WithKeys("e"),
				key.WithHelp("e", "send event"))}
		return b
	}

	return ProjectModel{
		common:   common,
		list:     l,
		projects: []*models.Project{},
	}
}

func (m ProjectModel) View() string {
	return "\n" + m.list.View()
}

func (m ProjectModel) Update(msg tea.Msg) (ProjectModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			cmds = append(cmds, common.LoadServices(m.common.KeptnAPI.ServicesV1(), m.projects[m.list.Index()]))
		case "d", "-":
			projectToDelete := m.projects[m.list.Index()]
			cmds = append(cmds, common.DeleteProject(m.common.KeptnAPI.APIV1(), projectToDelete.ProjectName))
		case "a", "+":
			cmds = append(cmds, common.ShowCreateProjectView())
		case "e":
			cmds = append(cmds, common.ShowSendEventView())
		}

	case common.ProjectLoadedMsg:
		m.projects = msg.Projects
		var items []list.Item
		for _, p := range msg.Projects {
			items = append(items, item(p.ProjectName))
		}
		cmds = append(cmds, m.list.SetItems(items))
	}

	newList, listCmd := m.list.Update(msg)
	m.list = newList
	cmds = append(cmds, listCmd)
	return m, tea.Batch(cmds...)
}
