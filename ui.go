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

	add    = "a"
	toggle = "t"
	change = "c"
	delete = "d"
	quit   = "q"
)

type color string

type Model struct {
	store Store

	collection taskCollection
	index      int

	isEditing bool
	taskInput textinput.Model
}

func NewUI(store Store) *Model {
	return &Model{
		store:      store,
		collection: store.LoadTasks(),
		taskInput:  textinput.NewModel(),
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
	m.store.SaveTasks(m.collection)
}

func (m *Model) updateTaskInputField(msg tea.Msg) (tea.Model, tea.Cmd, bool) {
	if !m.taskInput.Focused() {
		return nil, nil, false
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.isEditing = false
			m.taskInput.Reset()
			m.taskInput.Blur()

		case tea.KeyEnter:
			value := m.taskInput.Value()
			if editingExistingTask := m.isEditing; editingExistingTask {
				m.collection.change(m.index, value)
			} else {
				m.collection.append(value)
			}

			m.isEditing = false
			m.taskInput.Reset()
			m.taskInput.Blur()
		}
	}

	var cmd tea.Cmd
	m.taskInput, cmd = m.taskInput.Update(msg)
	return m, cmd, true
}

func (m *Model) updateTaskNavigator(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.index = max(0, m.index-1)

		case "down", "j":
			m.index = min(m.collection.len()-1, m.index+1)

		case add:
			m.index = max(m.index, 0)
			m.taskInput.Prompt = "> "
			m.taskInput.Placeholder = "describe the task..."
			m.taskInput.Focus()

		case toggle:
			m.collection.toggle(m.index)

		case change:
			summary := m.collection.summary(m.index)
			m.isEditing = true
			m.taskInput.Prompt = ""
			m.taskInput.SetValue(summary)
			m.taskInput.SetCursor(len(summary))
			m.taskInput.Focus()

		case delete:
			m.collection.delete(m.index)
			m.index = min(m.index, m.collection.len()-1)

		case "ctrl+c", quit:
			cmd = tea.Quit
		}
	}

	return m, cmd
}

func (m *Model) View() string {
	lines := []string{m.collection.name}

	for i, t := range m.collection.tasks {
		lines = append(lines, m.renderTask(i, t))
	}

	lines = append(lines, m.renderInputField()...)
	lines = append(lines, m.renderCommands()...)

	return strings.Join(lines, "\n")
}

func (m *Model) renderTask(index int, t task) string {
	color := yellow
	checked := " "
	if t.isCompleted {
		color = green
		checked = "x"
	}

	cursor := " "
	if index == m.index {
		cursor = ">"
	}

	summary := t.summary
	if m.isEditing && m.index == index {
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
		add, toggle, change, delete, quit,
	)}
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
