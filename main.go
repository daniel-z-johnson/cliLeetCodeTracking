package main

import (
	"fmt"
	"log"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/daniel-z-johnson/clileetcodetracking/jsondb"
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
	simpleDB        *jsondb.JsonDB
	problemView     bool
	add             bool
}

func initialModel() model {

	simpleDB, err := jsondb.NewJsonDB("data.json")
	if err != nil {
		log.Fatal(err)
	}

	return model{
		textInput:       clearTeatInput(),
		quitting:        false,
		leetCodeProblem: "",
		simpleDB:        simpleDB,
		problemView:     false,
		add:             false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.problemView {
		return m.UpdateWrite(msg)
	}
	return m.UpdateInput(msg)
}

func (m model) UpdateWrite(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.problemView = false
			m.simpleDB.Write(m.leetCodeProblem, nextSaturday())
		case "esc", "ctrl+c":
			m.quitting = true
			m.problemView = false
			return m, tea.Quit
		case "up", "down":
			m.add = !m.add
		}

	}
	return m, cmd
}

func (m model) UpdateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.leetCodeProblem = m.textInput.Value()
			m.textInput = clearTeatInput()
			m.problemView = true
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

func clearTeatInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Type something..."
	ti.Focus()
	ti.CharLimit = 4
	ti.SetWidth(20)
	return ti
}

func (m model) headerView() string {
	return fmt.Sprintf("Next Saturday: %s\nYou entered: %s", nextSaturday(), m.leetCodeProblem)
}

func (m model) footerView() string {
	return "\n(esc to quit)"
}

func nextSaturday() string {
	date := time.Now()
	for date.Weekday() != time.Saturday {
		date = date.AddDate(0, 0, 1)
	}
	return date.Format("2006-01-02")
}
