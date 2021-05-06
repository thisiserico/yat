package yat

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	nocolor color = "\033[0m"
	green   color = "\033[0;32m"
	yellow  color = "\033[0;33m"

	taskAbove        = "k"
	taskBelow        = "j"
	addTask          = "a"
	toggleTask       = "t"
	changeTask       = "c"
	deleteTask       = "d"
	collectionAbove  = "K"
	collectionBelow  = "J"
	changeCollection = "C"
	quit             = "q"

	none ongoingModification = iota
	editingTask
	addingTask
	renamingCollection
)

type color string

type ongoingModification int

type collection struct {
	store Store
	model taskCollection
	index int
}

type Model struct {
	collections []collection
	index       int

	modification ongoingModification
	inputField   textinput.Model
}

func NewUI(stores ...Store) *Model {
	collections := make([]collection, 0, len(stores))
	for _, store := range stores {
		collections = append(collections, collection{
			store: store,
			model: store.LoadTasks(),
		})
	}

	return &Model{
		collections:  collections,
		inputField:   textinput.NewModel(),
		modification: none,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.modification != none {
		return m.updateTaskInputField(msg)
	}

	return m.updateTaskNavigator(msg)
}

func (m *Model) Flush() {
	for _, collection := range m.collections {
		collection.store.SaveTasks(collection.model)
	}
}

func (m *Model) currentCollection() *collection {
	return &m.collections[m.index]
}

func (m *Model) updateTaskInputField(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.resetInputField()

		case tea.KeyEnter:
			value := m.inputField.Value()
			current := m.currentCollection()

			switch m.modification {
			case editingTask:
				current.model.change(current.index, value)

			case addingTask:
				current.model.append(value)

			case renamingCollection:
				current.model.rename(value)
			}

			m.resetInputField()
		}
	}

	var cmd tea.Cmd
	m.inputField, cmd = m.inputField.Update(msg)
	return m, cmd
}

func (m *Model) resetInputField() {
	m.modification = none
	m.inputField.Reset()
	m.inputField.Blur()
}

func (m *Model) updateTaskNavigator(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := m.currentCollection()

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case taskAbove:
			current.index = max(0, current.index-1)

		case taskBelow:
			current.index = min(current.model.len()-1, current.index+1)

		case addTask:
			current.index = max(current.index, 0)
			m.inputField.Prompt = "> "
			m.inputField.Placeholder = "describe the task..."
			m.inputField.Focus()
			m.modification = addingTask

		case toggleTask:
			current.model.toggle(current.index)

		case changeTask:
			summary := current.model.summary(current.index)
			m.inputField.Prompt = ""
			m.inputField.SetValue(summary)
			m.inputField.SetCursor(len(summary))
			m.inputField.Focus()
			m.modification = editingTask

		case deleteTask:
			current.model.delete(current.index)
			current.index = min(current.index, current.model.len()-1)

		case collectionAbove:
			m.index = max(0, m.index-1)

		case collectionBelow:
			m.index = min(len(m.collections)-1, m.index+1)

		case changeCollection:
			name := current.model.name
			m.inputField.Prompt = ""
			if name == "" {
				m.inputField.Placeholder = "name the collection..."
			}
			m.inputField.SetValue(name)
			m.inputField.SetCursor(len(name))
			m.inputField.Focus()
			m.modification = renamingCollection

		case "ctrl+c", quit:
			cmd = tea.Quit
		}
	}

	return m, cmd
}

func (m *Model) View() string {
	var lines []string

	for i, collection := range m.collections {
		lines = append(lines, m.renderCollection(collection, i == m.index, i > 0)...)
	}

	lines = append(lines, m.renderInputField()...)
	lines = append(lines, m.renderCommands()...)

	return strings.Join(lines, "\n")
}

func (m *Model) renderCollection(current collection, focusedOnIt, prependEmptyLine bool) []string {
	prepended := ""
	if prependEmptyLine {
		prepended = "\n"
	}

	name := current.model.name
	if m.modification == renamingCollection && focusedOnIt {
		name = m.inputField.View()
	}

	lines := []string{prepended + name}
	for i, t := range current.model.tasks {
		lines = append(lines, m.renderTask(t, focusedOnIt && i == current.index))
	}

	return lines
}

func (m *Model) renderTask(t task, focusedOnIt bool) string {
	color := yellow
	checked := " "
	if t.isCompleted {
		color = green
		checked = "x"
	}

	cursor := " "
	if focusedOnIt {
		cursor = ">"
	}

	summary := t.summary
	if m.modification == editingTask && focusedOnIt {
		summary = m.inputField.View()
	}

	return fmt.Sprintf(
		"%s %s[%s] %s%s",
		cursor, color, checked, summary, nocolor,
	)
}

func (m *Model) renderInputField() []string {
	if m.modification != addingTask {
		return nil
	}

	return []string{fmt.Sprintf("\n%s", m.inputField.View())}
}

func (m *Model) renderCommands() []string {
	if m.modification != none {
		return nil
	}

	return []string{fmt.Sprintf(
		"\n%s appends, %s toggles, %s changes, %s deletes and %s quits",
		addTask, toggleTask, changeTask, deleteTask, quit,
	)}
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
