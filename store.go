package yat

import (
	"os"
	"time"

	"github.com/pelletier/go-toml"
)

type Store interface {
	LoadTasks() taskCollection
	SaveTasks(taskCollection)
}

type tomlStore struct {
	path string
}

type tomlConfig struct {
	Collection string
	Tasks      []tomlTask
}

type tomlTask struct {
	ID          string
	Summary     string
	IsCompleted bool
	AddedAt     time.Time
}

func NewTomlStore(path string) Store {
	return &tomlStore{
		path: path,
	}
}

func (s *tomlStore) LoadTasks() taskCollection {
	file, err := os.OpenFile(s.path, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config tomlConfig
	decoder := toml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		panic(err)
	}

	collection := taskCollection{name: config.Collection}
	for _, t := range config.Tasks {
		collection.tasks = append(collection.tasks, task{
			sortableID:  t.ID,
			summary:     t.Summary,
			isCompleted: t.IsCompleted,
			addedAt:     t.AddedAt,
		})
	}

	return collection
}

func (s *tomlStore) SaveTasks(collection taskCollection) {
	file, err := os.OpenFile(s.path, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	config := tomlConfig{Collection: collection.name}
	for _, t := range collection.tasks {
		config.Tasks = append(config.Tasks, tomlTask{
			ID:          t.sortableID,
			Summary:     t.summary,
			IsCompleted: t.isCompleted,
			AddedAt:     t.addedAt,
		})
	}

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		panic(err)
	}
}
