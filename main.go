package main

import (
	"fmt"
	table "github.com/calyptia/go-bubble-table"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"optom-terminal-ehr/db"
	"os"
)

type patients struct {
	FirstName string
	LastName  string
}

type model struct {
	viewtype int
	table    table.Model
}

type keyMap struct {
	Up key.Binding
}

var styleDoc = lipgloss.NewStyle().Padding(1)

func initialModel() model {
	pxs := []patients{
		{FirstName: "One", LastName: "PatientOne"},
		{FirstName: "Two", LastName: "PatientTwo"},
		{FirstName: "Three", LastName: "PatientThree"},
		{FirstName: "Four", LastName: "PatientFour"},
		{FirstName: "Five", LastName: "PatientFive"},
		{FirstName: "Six", LastName: "PatientSix"},
		{FirstName: "Seven", LastName: "PatientSeven"},
		{FirstName: "Eight", LastName: "PatientEight"},
	}
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
		h = 24
	}
	top, right, bottom, left := styleDoc.GetPadding()
	w = w - left - right
	h = h - top - bottom
	tbl := table.New([]string{"LASTNAME", "FIRSTNAME"}, w, h)
	rows := make([]table.Row, len(pxs))
	for i, px := range pxs {
		rows[i] = table.SimpleRow{
			px.FirstName,
			px.LastName,
		}
	}
	tbl.SetRows(rows)
	return model{
		table: tbl,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := styleDoc.GetPadding()
		m.table.SetSize(
			msg.Width-left-right,
			msg.Height-top-bottom,
		)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(message)
	return m, cmd
}

func (m model) View() string {
	return styleDoc.Render(
		m.table.View(),
	)
}

func main() {
	_, err := db.SetUpDb()

	if err != nil {
		panic(err)
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	err = p.Start()
	if err != nil {
		_ = fmt.Errorf("Alas, there's been an error: %v", err)
	}
}
