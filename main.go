package main

import (
	"fmt"
	"log"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	textInput       textinput.Model
	leetCodeProblem string
	quitting        bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type something..."
	ti.Focus()
	ti.CharLimit = 4
	ti.SetWidth(20)

	return model{
		textInput:       ti,
		quitting:        false,
		leetCodeProblem: "",
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.leetCodeProblem = m.textInput.Value()
			ti := textinput.New()
			ti.Placeholder = "Type something..."
			ti.Focus()
			ti.CharLimit = 4
			ti.SetWidth(20)
			m.textInput = ti
			return m, cmd
		case "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	var c *tea.Cursor
	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
		c.Y += lipgloss.Height(m.headerView())
	}

	str := lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.textInput.View(), m.footerView())
	if m.quitting {
		str += "\n"
	}

	v := tea.NewView(str)
	v.Cursor = c
	return v
}

func (m model) headerView() string {
	return fmt.Sprintf("You entered: %s", m.leetCodeProblem)
}

func (m model) footerView() string {
	return "\n(esc to quit)"
}
