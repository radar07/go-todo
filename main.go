package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	Name string `json:"name"`
}

type Todo struct {
	items    []Item
	cursor   int
	selected map[int]struct{}
}

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

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, item.Name)
	}

	s += "\n('q' - quit)"
	return s
}

func InitialModel() Todo {
	items := GetItemsFromFile()
	return Todo{
		items:    items,
		selected: make(map[int]struct{}),
	}
}

func (t Todo) Init() tea.Cmd {
	return nil
}

func (t Todo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := tea.EnterAltScreen
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return t, tea.Quit

		case "k", "up":
			cmd = tea.ClearScreen
			if t.cursor > 0 {
				t.cursor--
			}
		case "j", "down":
			cmd = tea.ClearScreen
			if t.cursor < len(t.items)-1 {
				t.cursor++
			}
		case "n":
			panic("unimplemented!")
		case "u":
			fmt.Println(t.cursor)
		case "d":
			panic("unimplemented!")
		case "enter", " ":
			_, ok := t.selected[t.cursor]
			if ok {
				delete(t.selected, t.cursor)
			} else {
				t.selected[t.cursor] = struct{}{}
			}
		}
	}

	return t, cmd
}

func main() {
	SeedFile()
	p := tea.NewProgram(InitialModel())
	_, err := p.Run()

	if err != nil {
		os.Exit(1)
	}
}
