package main

import (
	"fmt"
	"os"

	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	Name string `json:"name"`
}

type Todo struct {
	textInput textinput.Model
	items     []Item
	cursor    int
	selected  map[int]struct{}
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

	s += fmt.Sprintf("\n%s", t.textInput.View())
	s += "\n('q' - quit)"
	return s
}

func InitialModel() Todo {
	ti := textinput.New()
	ti.Placeholder = "Enter new item: "

	items := GetItemsFromFile()
	return Todo{
		textInput: ti,
		items:     items,
		selected:  make(map[int]struct{}),
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
			return t, tea.Quit

		case "k", "up":
			cmd = tea.ClearScreen
			if t.cursor > 0 && !t.textInput.Focused() {
				t.cursor--
			}

		case "j", "down":
			cmd = tea.ClearScreen
			if t.cursor < len(t.items)-1 && !t.textInput.Focused() {
				t.cursor++
			}

		case " ", tea.KeyEnter.String():
			if t.textInput.Focused() {
				item := t.textInput.Value()
				if item != "" {
					t.items = append(t.items, Item{item})
					t.textInput.Reset()
				}
			} else {
				_, ok := t.selected[t.cursor]
				if ok {
					delete(t.selected, t.cursor)
				} else {
					t.selected[t.cursor] = struct{}{}
				}
			}

		case tea.KeyTab.String():
			if t.textInput.Focused() {
				t.textInput.Blur()
			} else {
				cmd = t.textInput.Focus()
			}
		}

		t.textInput, cmd = t.textInput.Update(msg)
	}

	return t, cmd
}

func main() {
	SeedFile()
	p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	_, err := p.Run()

	if err != nil {
		os.Exit(1)
	}
}
