package yat

import "time"

type Store interface {
	LoadTasks() tasks
	SaveTasks(tasks)
}

type dummyStore struct{}

func NewStore() Store {
	return &dummyStore{}
}

func (d *dummyStore) LoadTasks() tasks {
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

func (d *dummyStore) SaveTasks(_ tasks) {}
