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

	tasks
	index int

	isEditing bool
	taskInput textinput.Model

	logs []string
}

func NewUI(store Store) *Model {
	return &Model{
		store:     store,
		tasks:     store.LoadTasks(),
		taskInput: textinput.NewModel(),
	}
}

func (m *Model) log(msg string) {
	m.logs = append(m.logs, msg)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.taskInput.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEsc:
				m.isEditing = false
				m.taskInput.Reset()
				m.taskInput.Blur()

			case tea.KeyEnter:
				if m.isEditing {
					m.isEditing = false
					m.tasks[m.index].replace(m.taskInput.Value())
				} else {
					m.tasks.append(m.taskInput.Value())
				}

				m.taskInput.Reset()
				m.taskInput.Blur()

			default:
				var cmd tea.Cmd
				m.taskInput, cmd = m.taskInput.Update(msg)

				return m, cmd
			}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.index = max(0, m.index-1)

		case "down", "j":
			m.index = min(len(m.tasks)-1, m.index+1)

		case add:
			m.taskInput.Placeholder = "describe the task..."
			m.taskInput.Focus()
			m.taskInput.Prompt = "> "
			m.index = max(m.index, 0)

		case toggle:
			m.tasks[m.index].toggle()

		case change:
			summary := m.tasks[m.index].summary
			m.isEditing = true
			m.taskInput.SetValue(summary)
			m.taskInput.Focus()
			m.taskInput.SetCursor(len(summary))
			m.taskInput.Prompt = ""

		case delete:
			m.tasks.delete(m.index)
			m.index = min(m.index, len(m.tasks)-1)

		case "ctrl+c", quit:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) Flush() {
	m.store.SaveTasks(m.tasks)
}

func (m *Model) View() string {
	lines := []string{"generic tasks"}
	lines = append(lines, m.logs...)

	for i, t := range m.tasks {
		lines = append(lines, m.renderTask(i, t))
	}

	lines = append(lines, m.renderInputField()...)
	lines = append(lines, m.renderCommands()...)

	return strings.Join(lines, "\n")
}

func (m *Model) renderTask(index int, t *task) string {
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
