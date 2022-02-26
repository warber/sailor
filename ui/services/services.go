package services

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

type ServiceModel struct {
	common   *common.CommonModel
	list     list.Model
	choice   string
	quitting bool
}

func NewServiceModel(common *common.CommonModel) ServiceModel {
	items := []list.Item{}
	l := list.New(items, itemDelegate{}, 50, 14)
	l.Title = "Select a Keptn Service"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		b := []key.Binding{
			key.NewBinding(key.WithKeys("a", "+"),
				key.WithHelp("a/+", "add service")),
			key.NewBinding(key.WithKeys("b"),
				key.WithHelp("b/←", "back"))}
		return b
	}

	return ServiceModel{
		common: common,
		list:   l,
	}
}

func (m ServiceModel) View() string {
	return "\n" + m.list.View()
}

func (m ServiceModel) Update(msg tea.Msg) (ServiceModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b", "left":
			cmds = append(cmds, common.ShowProjectsReq())
		}
	case common.ServiceLoadedMsg:
		var items []list.Item
		for _, s := range msg.Services {
			items = append(items, item(s.ServiceName))
		}
		cmds = append(cmds, m.list.SetItems(items))
	}

	newList, listCmd := m.list.Update(msg)
	m.list = newList
	cmds = append(cmds, listCmd)
	return m, tea.Batch(cmds...)
}
