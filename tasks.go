package yat

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type tasks []*task

type task struct {
	id          string
	summary     string
	isCompleted bool
	addedAt     time.Time
}

func newTaskID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), nil).String()
}

func (t *tasks) append(summary string) {
	*t = append(*t, &task{
		id:      newTaskID(),
		summary: summary,
		addedAt: time.Now(),
	})
}

func (t *tasks) delete(index int) {
	*t = append((*t)[:index], (*t)[index+1:]...)
}

func (t *task) toggle() {
	t.isCompleted = !t.isCompleted
}

func (t *task) replace(summary string) {
	t.summary = summary
}
