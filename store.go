package yat

import (
	"os"
	"time"

	"github.com/pelletier/go-toml"
)

type Store interface {
	LoadTasks() tasks
	SaveTasks(tasks)
}

type tomlStore struct {
	path string
}

type tomlConfig struct {
	Tasks []tomlTask
}

type tomlTask struct {
	Summary     string
	IsCompleted bool
	AddedAt     time.Time
}

func NewTomlStore(path string) Store {
	return &tomlStore{
		path: path,
	}
}

func (t *tomlStore) LoadTasks() tasks {
	file, err := os.OpenFile(t.path, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config tomlConfig
	decoder := toml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		panic(err)
	}

	var tasks tasks
	for _, t := range config.Tasks {
		tasks = append(tasks, &task{
			summary:     t.Summary,
			isCompleted: t.IsCompleted,
			addedAt:     t.AddedAt,
		})
	}

	return tasks
}

func (t *tomlStore) SaveTasks(tasks tasks) {
	file, err := os.OpenFile(t.path, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	config := tomlConfig{}
	for _, t := range tasks {
		config.Tasks = append(config.Tasks, tomlTask{
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
