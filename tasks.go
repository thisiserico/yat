package yat

import "time"

type tasks []*task

type task struct {
	summary     string
	isCompleted bool
	addedAt     time.Time
}

func fakeTaskList() tasks {
	return tasks{
		{
			summary: "a brand new trask",
			addedAt: time.Now(),
		},
		{
			summary: "this one is not so cool tho...",
			addedAt: time.Now(),
		},
		{
			summary:     "but this one is completed!",
			isCompleted: true,
			addedAt:     time.Now(),
		},
	}
}

func (t *tasks) append(summary string) {
	*t = append(*t, &task{
		summary: summary,
		addedAt: time.Now(),
	})
}

func (t *task) toggle() {
	t.isCompleted = !t.isCompleted
}
