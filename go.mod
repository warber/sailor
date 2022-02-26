module github.com/warber/sailor

go 1.16

require (
	github.com/charmbracelet/bubbles v0.10.2
	github.com/charmbracelet/bubbletea v0.20.0
	github.com/charmbracelet/charm v0.10.1
	github.com/charmbracelet/lipgloss v0.5.0
	github.com/keptn/go-utils v0.13.0
	github.com/muesli/termenv v0.11.1-0.20220212125758-44cd13922739
)

replace (
   github.com/keptn/go-utils => ../go-utils
)