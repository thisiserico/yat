package yat

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"time"
)

type Store interface {
	LoadTasks() tasks
	SaveTasks(tasks)
}

type fileStore struct {
	file *os.File
}

func NewFileStore(file string) Store {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}

	return &fileStore{
		file: f,
	}
}

func (f *fileStore) LoadTasks() tasks {
	offset, err := f.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil
	}

	scanner := bufio.NewScanner(f.file)
	scanner.Split(func(data []byte, atEof bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEof)
		if err == nil && token != nil {
			offset += int64(advance)
		}

		return advance, token, err
	})

	var tasks tasks
	regex, _ := regexp.Compile("^\\[([x ])\\] (.+) \\| (.+)$")
	for scanner.Scan() {
		line := scanner.Bytes()
		if !regex.Match(line) {
			continue
		}

		submatches := regex.FindSubmatch(line)

		addedAt := time.Now().UTC()
		if existingAddedAt, err := time.Parse(time.RFC3339, string(submatches[3])); err == nil {
			addedAt = existingAddedAt
		}
		tasks = append(tasks, &task{
			summary:     string(submatches[2]),
			isCompleted: string(submatches[1]) == "x",
			addedAt:     addedAt,
		})
	}

	return tasks
}

func (f *fileStore) SaveTasks(_ tasks) {
}
