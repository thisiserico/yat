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

	taskAbove       = "k"
	collectionAbove = "K"
	taskBelow       = "j"
	collectionBelow = "J"
	addTask         = "a"
	toggleTask      = "t"
	changeTask      = "c"
	deleteTask      = "d"
	quit            = "q"
)

type color string

type collection struct {
	store Store
	model taskCollection
	index int
}

type Model struct {
	collections []collection
	index       int

	isEditing bool
	taskInput textinput.Model
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
		collections: collections,
		taskInput:   textinput.NewModel(),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m, cmd, handled := m.updateTaskInputField(msg); handled {
		return m, cmd
	}

	return m.updateTaskNavigator(msg)
}

func (m *Model) Flush() {
	current := m.currentCollection()
	current.store.SaveTasks(current.model)
}

func (m *Model) currentCollection() *collection {
	return &m.collections[m.index]
}

func (m *Model) updateTaskInputField(msg tea.Msg) (tea.Model, tea.Cmd, bool) {
	if !m.taskInput.Focused() {
		return nil, nil, false
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.resetInputField()

		case tea.KeyEnter:
			current := m.currentCollection()
			value := m.taskInput.Value()
			if editingExistingTask := m.isEditing; editingExistingTask {
				current.model.change(current.index, value)
			} else {
				current.model.append(value)
			}

			m.resetInputField()
		}
	}

	var cmd tea.Cmd
	m.taskInput, cmd = m.taskInput.Update(msg)
	return m, cmd, true
}

func (m *Model) resetInputField() {
	m.isEditing = false
	m.taskInput.Reset()
	m.taskInput.Blur()
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
			m.taskInput.Prompt = "> "
			m.taskInput.Placeholder = "describe the task..."
			m.taskInput.Focus()

		case toggleTask:
			current.model.toggle(current.index)

		case changeTask:
			summary := current.model.summary(current.index)
			m.isEditing = true
			m.taskInput.Prompt = ""
			m.taskInput.SetValue(summary)
			m.taskInput.SetCursor(len(summary))
			m.taskInput.Focus()

		case deleteTask:
			current.model.delete(current.index)
			current.index = min(current.index, current.model.len()-1)

		case collectionAbove:
			m.index = max(0, m.index-1)

		case collectionBelow:
			m.index = min(len(m.collections)-1, m.index+1)

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

	lines := []string{prepended + current.model.name}
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
	if m.isEditing && focusedOnIt {
		summary = m.taskInput.View()
	}

	return fmt.Sprintf(
		"%s %s[%s] %s%s",
		cursor, color, checked, summary, nocolor,
	)
}

func (m *Model) renderInputField() []string {
	if m.isEditing || !m.taskInput.Focused() {
		return nil
	}

	return []string{fmt.Sprintf("\n%s", m.taskInput.View())}
}

func (m *Model) renderCommands() []string {
	if m.taskInput.Focused() {
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
