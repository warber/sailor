package newproject

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/warber/sailor/ui/common"
	"strings"
)

type createProjectViewState int

const (
	createProjectViewStateReady createProjectViewState = iota
	createProjectViewStateShowError
)

type CreateProjectModel struct {
	common     *common.CommonModel
	focusIndex int
	inputs     []textinput.Model
	viewState  createProjectViewState
	err        error
}

func NewCreateProjectModel(common *common.CommonModel) CreateProjectModel {
	m := CreateProjectModel{
		common: common,
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Keptn Project Name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Shipyard File Location"
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		m.inputs[i] = t
	}
	return m
}

func (m CreateProjectModel) View() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(list.DefaultStyles().Title.Render("Create a new Keptn project"))
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	switch m.viewState {
	case createProjectViewStateShowError:
		b.WriteString(indent("Ooops.. An error occured: "+m.err.Error(), 2))
	}
	return b.String()
}

// Lightweight version of reflow's indent function.
func indent(s string, n int) string {
	if n <= 0 || s == "" {
		return s
	}
	l := strings.Split(s, "\n")
	b := strings.Builder{}
	i := strings.Repeat(" ", n)
	for _, v := range l {
		fmt.Fprintf(&b, "%s%s\n", i, v)
	}
	return b.String()
}

func (m CreateProjectModel) Update(msg tea.Msg) (CreateProjectModel, tea.Cmd) {
	textinput.Blink() //TODO check where this shall be moved?!?!?
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "b":
			return m, common.ShowProjectsReq()
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				projectName := m.inputs[0].Value()
				shipyardLocation := m.inputs[1].Value()
				return m, common.CreateProject(m.common.KeptnAPI.APIV1(), projectName, shipyardLocation)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}
			return m, tea.Batch(cmds...)
		}
		// Handle character input and blinking
		return m, m.updateInputs(msg)
	case common.ProjectCreatedErrMsg:
		m.viewState = createProjectViewStateShowError
		m.err = msg.Err
	}

	return m, nil
}

func (m *CreateProjectModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// Update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
