package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// App state
type Todo struct {
	items    []string
	cursor   int
	selected map[int]struct{}
}

// View implements tea.Model.
func (t Todo) View() string {
	s := "TODO:\n"

	for i, item := range t.items {
		cursor := " "
		if t.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := t.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, item)
	}

	s += "\nPress 'Q' to quit!"
	return s
}

func initialModel() Todo {
	return Todo{
		items:    []string{"Learn Go", "Use Github APIs"},
		selected: make(map[int]struct{}),
	}
}

// Kick-off the event loop
func (t Todo) Init() tea.Cmd {
	return nil
}

func (t Todo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return t, tea.Quit

		case "k", "up":
			if t.cursor > 0 {
				t.cursor--
			}
		case "j", "down":
			if t.cursor < len(t.items)-1 {
				t.cursor++
			}
		case "enter", " ":
			_, ok := t.selected[t.cursor]
			if ok {
				delete(t.selected, t.cursor)
			} else {
				t.selected[t.cursor] = struct{}{}
			}
		}
	}

	return t, nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
