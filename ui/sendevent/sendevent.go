package sendevent

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/warber/sailor/ui/common"
	"strings"
)

type SendEventModel struct {
	common           *common.CommonModel
	focusIndex       int
	inputs           []textinput.Model
	lastSentEventCtx string
}

func NewSendEventModel(common *common.CommonModel) SendEventModel {
	m := SendEventModel{
		common: common,
		inputs: make([]textinput.Model, 1),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Event File Location"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		m.inputs[i] = t
	}
	return m
}

func (m SendEventModel) View() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(list.DefaultStyles().Title.Render("Send a Keptn Event"))
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
	if m.lastSentEventCtx != "" {
		b.WriteString("Sent event with context " + m.lastSentEventCtx)
	}

	return b.String()
}

func (m SendEventModel) Update(msg tea.Msg) (SendEventModel, tea.Cmd) {
	textinput.Blink()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "b", "left":
			return m, common.ShowProjectsReq()
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				eventFileLocation := m.inputs[0].Value()
				return m, common.SendEvent(m.common.KeptnAPI.APIV1(), eventFileLocation)
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
	case common.EventSentMsg:
		m.lastSentEventCtx = msg.KeptnCtx
	}

	return m, nil
}

func (m *SendEventModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// Update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
