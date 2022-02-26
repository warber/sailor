package common

import (
	"encoding/json"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/fileutils"
)

type EventSentMsg struct{ KeptnCtx string }
type EventSentErrMsg struct{ err error }
type ShowSendEventMsg struct{}

func SendEvent(api api.APIV1Interface, eventFileLocation string) tea.Cmd {
	return func() tea.Msg {
		eventString, err := fileutils.ReadFile(eventFileLocation)
		if err != nil {
			return EventSentErrMsg{err}
		}
		apiEvent := models.KeptnContextExtendedCE{}
		err = json.Unmarshal(eventString, &apiEvent)
		if err != nil {
			return EventSentErrMsg{err}
		}

		context, err2 := api.SendEvent(apiEvent)
		if err2 != nil {
			return EventSentErrMsg{errors.New(err2.GetMessage())}
		}

		return EventSentMsg{*context.KeptnContext}
	}
}

func ShowSendEventView() tea.Cmd {
	return func() tea.Msg {
		return ShowSendEventMsg{}
	}
}
