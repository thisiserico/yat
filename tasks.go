package yat

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type taskCollection struct {
	name  string
	tasks []task
}

type task struct {
	sortableID  string
	summary     string
	isCompleted bool
	addedAt     time.Time
}

func newTaskID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), nil).String()
}

func (col *taskCollection) summary(index int) string {
	return col.tasks[index].summary
}

func (col *taskCollection) len() int {
	return len(col.tasks)
}

func (col *taskCollection) append(summary string) {
	col.tasks = append(col.tasks, task{
		sortableID: newTaskID(),
		summary:    summary,
		addedAt:    time.Now(),
	})
}

func (col *taskCollection) toggle(index int) {
	col.tasks[index].isCompleted = !col.tasks[index].isCompleted
}

func (col *taskCollection) rename(name string) {
	col.name = name
}

func (col *taskCollection) change(index int, summary string) {
	col.tasks[index].summary = summary
}

func (col *taskCollection) delete(index int) {
	col.tasks = append(col.tasks[:index], col.tasks[index+1:]...)
}
