package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

type model struct {
	cities  table.Model
	pokemon table.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func Unselect() table.Styles {
	unselected := table.DefaultStyles()

	unselected.Header = unselected.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	unselected.Selected = unselected.Selected.
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("0")).
		Bold(false)
	return unselected
}

func Select() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd, cmd2 tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.cities.Focused() {
				m.cities.Blur()
			} else {
				m.cities.Focus()
			}
		case "tab":
			if m.cities.Focused() {
				m.pokemon.Focus()
				m.pokemon.SetStyles(Select())
				m.cities.Blur()
				m.cities.SetStyles(Unselect())
			} else {
				m.cities.Focus()
				m.cities.SetStyles(Select())
				m.pokemon.Blur()
				m.pokemon.SetStyles(Unselect())
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.cities.Focused() {
				return m, tea.Printf("Let's go to %s!", m.cities.SelectedRow()[1])
			}
			return m, tea.Printf("I choose you %s!", m.pokemon.SelectedRow()[1])
		}
	}
	m.cities, cmd = m.cities.Update(msg)
	m.pokemon, cmd2 = m.pokemon.Update(msg)
	return m, tea.Batch(cmd, cmd2)
}

func (m model) View() string {
	return fmt.Sprintf("\tCities\n%s\n\n\tPokemon\n%s\n", baseStyle.Render(m.cities.View()), baseStyle.Render(m.pokemon.View()))
}

func main() {
	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "City", Width: 10},
		{Title: "Country", Width: 10},
		{Title: "Population", Width: 10},
	}
	cities, err := ParseCities("./data/cities.csv")
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	rows := make([]table.Row, len(cities))
	for i, city := range cities {
		rows[i] = []string{fmt.Sprint(city.Id), city.City, city.Country, fmt.Sprint(city.Population)}
	}

	t, t2 := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	), table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	t2.SetStyles(s)

	m := model{t, t2}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
