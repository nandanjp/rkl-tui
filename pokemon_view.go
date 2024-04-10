package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	blurredStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("236"))
	cursorStyle             = focusedStyle.Copy()
	placeholderStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	focusedPlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	noStyle                 = lipgloss.NewStyle()
	helpStyle               = blurredStyle.Copy()
	cursorModeHelpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	focusedButton = focusedStyle.Copy().Render("[Submit]")
	blurredButton = blurredStyle.Copy().Render("Submit")

	modelStyle = lipgloss.NewStyle().Width(30).Height(30).Align(lipgloss.Center, lipgloss.Center).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("216"))
)

type keymap = struct {
	next, prev, mode, search, scrollUp, scrollDown, change, quit key.Binding
}

type PokemonViewModel struct {
	width      int
	height     int
	keymap     keymap
	help       help.Model
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	output     viewport.Model
	onViewport bool
}

func newView(content string) (*viewport.Model, error) {
	const width = 100
	vp := viewport.New(width, 30)
	vp.Style = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("79")).Padding(2)

	str, err := glamour.Render(content, "dark")
	if err != nil {
		fmt.Printf("There was an error: %v", err)
		return nil, err
	}
	vp.SetContent(str)
	return &vp, nil
}

func newInput(placeHolder string) textinput.Model {
	in := textinput.New()
	in.Placeholder = placeHolder
	in.PromptStyle = focusedStyle
	in.TextStyle = focusedStyle
	in.Cursor.Style = cursorStyle
	in.CharLimit = 15
	return in
}

func NewPokemonViewModel(content string) (PokemonViewModel, error) {
	vp, err := newView(content)
	if err != nil {
		return PokemonViewModel{}, err
	}
	m := PokemonViewModel{
		inputs: []textinput.Model{newInput("Arceus")},
		keymap: keymap{
			next:       key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
			prev:       key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev")),
			mode:       key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "mode")),
			search:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "search")),
			change:     key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "change to viewport")),
			scrollUp:   key.NewBinding(key.WithKeys("k"), key.WithHelp("up", "scroll up")),
			scrollDown: key.NewBinding(key.WithKeys("j"), key.WithHelp("down", "scroll down")),
			quit:       key.NewBinding(key.WithKeys("esc", "ctrl+c", "q"), key.WithHelp("esc", "quit")),
		},
		output:     *vp,
		onViewport: false,
	}
	m.inputs[m.focusIndex].Focus()
	return m, nil
}

func (m PokemonViewModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PokemonViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.mode):
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds = make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keymap.search):
			if m.focusIndex == len(m.inputs)-1 {
				pokemon, err := Pokemon(m.inputs[m.focusIndex].Value())
				var vp *viewport.Model
				if err != nil {
					vp, err = newView(err.Error())
					if err != nil {
						return m, tea.Batch(cmds...)
					}
				} else {
					vp, err = newView(pokemon.ToMarkdown().String())
					if err != nil {
						return m, tea.Batch(cmds...)
					}
				}
				m.output = *vp
			}
		case key.Matches(msg, m.keymap.next):
			m.inputs[m.focusIndex].Blur()
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = 0
			}
			cmds = append(cmds, m.inputs[m.focusIndex].Focus())
		case key.Matches(msg, m.keymap.prev):
			m.inputs[m.focusIndex].Blur()
			m.focusIndex++
			if m.focusIndex > len(m.inputs)-1 {
				m.focusIndex = 0
			}
			cmds = append(cmds, m.inputs[m.focusIndex].Focus())
		case key.Matches(msg, m.keymap.scrollUp), key.Matches(msg, m.keymap.scrollDown):
			var cmd tea.Cmd
			m.output, cmd = m.output.Update(msg)
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.change):
			if m.onViewport {
				m.onViewport = false
				m.inputs[m.focusIndex].Focus()
			} else {
				m.onViewport = true
				m.inputs[m.focusIndex].Blur()
			}
			m.output.MouseWheelEnabled = m.onViewport
		case key.Matches(msg, m.keymap.quit):
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}
	inputCmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], inputCmds[i] = m.inputs[i].Update(msg)
	}
	cmds = append(cmds, inputCmds...)
	return m, tea.Batch(cmds...)
}

func (m PokemonViewModel) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.prev,
		m.keymap.mode,
		m.keymap.search,
		m.keymap.quit,
	})

	var views []string
	for i := range m.inputs {
		views = append(views, m.inputs[i].View())
	}

	view := lipgloss.JoinVertical(lipgloss.Center, modelStyle.Render(fmt.Sprintf("%4s", views)))

	var b strings.Builder
	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "%s\n\n", *button)
	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return lipgloss.JoinHorizontal(lipgloss.Top, view, m.output.View()) + "\n\n" + lipgloss.JoinVertical(lipgloss.Left, b.String(), helpStyle.Render(help)+"\n\n")
}
