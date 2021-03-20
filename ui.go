package yat

import (
	"fmt"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	nocolor color = "\033[0m"
	green   color = "\033[0;32m"
	yellow  color = "\033[0;33m"

	add    command = "a"
	toggle command = "t"
	change command = "c"
	quit   command = "q"
)

type color string

type command string

type Model struct {
	tasks
	index int
}

func NewUI() *Model {
	return &Model{
		tasks: fakeTaskList(),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.index = max(0, m.index-1)

		case "down", "j":
			m.index = min(len(m.tasks)-1, m.index+1)

		case "t":
			m.tasks[m.index].isCompleted = !m.tasks[m.index].isCompleted

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) View() string {
	lines := []string{"generic tasks\n"}

	for i, t := range m.tasks {
		lines = append(lines, m.renderTask(i, t))
	}

	lines = append(lines, renderCommands())
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

	return fmt.Sprintf("%s %s[%s] %s%s", cursor, color, checked, t.summary, nocolor)
}

func renderCommands() string {
	return fmt.Sprintf(
		"\n%s appends, %s toggles, %s changes and %s quits",
		add,
		toggle,
		change,
		quit,
	)
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
