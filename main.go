package main

import (
	"fmt"
	"os"
	"strings"

	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	promptStyle         = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4"))
	// focusedButton = focusedStyle.Render("[ Submit ]")
	// blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type Item struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

type Todo struct {
	textInput textinput.Model
	items     []Item
	cursor    int
	changed   bool
}

func (t Todo) View() string {
	var s strings.Builder
	s.WriteString("\nTODO:\n")

	for i, item := range t.items {
		cursor := " "
		if t.cursor == i {
			cursor = ">"
		}

		checked := " "
		if t.items[i].Completed {
			checked = "󰄬"
		}

		s.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, item.Name))
	}

	s.WriteString(fmt.Sprintf("\n%s\n", t.textInput.View()))
	s.WriteString(helpStyle.Render("\n'q' - quit"))
	s.WriteString(helpStyle.Render("\n'TAB' - switch to prompt"))
	s.WriteString(helpStyle.Render("\n'u' - update"))
	s.WriteString(helpStyle.Render("\n'd' - delete"))
	return s.String()
}

func InitialModel() Todo {
	ti := textinput.New()
	ti.Cursor.Style = cursorStyle
	ti.Placeholder = "Enter new item: "

	items := GetItemsFromFile()

	return Todo{
		textInput: ti,
		items:     items,
	}
}

func (t Todo) Init() tea.Cmd {
	return nil
}

func (t Todo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", tea.KeyCtrlC.String():
			if !t.textInput.Focused() {
				WriteTodos(t.items)
				return t, tea.Quit
			}

		case "k", "up":
			if t.cursor > 0 && !t.textInput.Focused() {
				t.cursor--
			}

		case "j", "down":
			if t.cursor < len(t.items)-1 && !t.textInput.Focused() {
				t.cursor++
			}

		case "u":
			// FIXME: appends 'u' to the textinput
			(&t).changed = true
			if !t.textInput.Focused() {
				t.textInput.Focus()
				t.textInput.SetValue(t.items[t.cursor].Name)
			}

		case "d", tea.KeyDelete.String():
			if !t.textInput.Focused() {
				t.items = append(t.items[0:t.cursor], t.items[t.cursor+1:]...)
			}

		case tea.KeyEnter.String():
			if t.textInput.Focused() {
				item := t.textInput.Value()
				if item != "" {
					if t.changed {
						(&t.items[t.cursor]).Name = item
					} else {
						t.items = append(t.items, Item{Name: item})
						t.textInput.Reset()
					}
				}
			} else {
				isCompleted := t.items[t.cursor].Completed
				(&t.items[t.cursor]).Completed = !isCompleted
			}

		case tea.KeyTab.String():
			if t.textInput.Focused() {
				t.textInput.PromptStyle = noStyle
				t.textInput.TextStyle = noStyle
				t.textInput.Reset()
				t.textInput.Blur()
			} else {
				cmd = t.textInput.Focus()
				t.textInput.PromptStyle = focusedStyle
				t.textInput.TextStyle = focusedStyle
			}
		}

		t.textInput, cmd = t.textInput.Update(msg)
	}

	return t, cmd
}

func main() {
	// SeedFile()
	p := tea.NewProgram(InitialModel())
	_, err := p.Run()

	if err != nil {
		os.Exit(1)
	}
}
