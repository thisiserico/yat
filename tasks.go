package yat

import "time"

type tasks []*task

type task struct {
	summary     string
	isCompleted bool
	addedAt     time.Time
}

func (t *tasks) append(summary string) {
	*t = append(*t, &task{
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
