package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thisiserico/yat"
)

type filesFlag []string

func main() {
	var files filesFlag
	debug := flag.Bool("debug", false, "use ./yat.log as log output")
	flag.Var(&files, "file", "file to store tasks at")
	flag.Parse()

	defer prepareLooger(*debug)

	collection := expandPath("~/.yat")
	if len(files) > 0 {
		collection = expandPath(files[0])
	}

	store := yat.NewTomlStore(collection)
	model := yat.NewUI(store)
	defer model.Flush()

	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (i *filesFlag) String() string {
	return "[]"
}

func (i *filesFlag) Set(collection string) error {
	*i = append(*i, collection)
	return nil
}

type closer func() error

func prepareLooger(debugEnabled bool) closer {
	if !debugEnabled {
		return func() error {
			return nil
		}
	}

	f, err := tea.LogToFile("yat.log", "")
	if err != nil {
		panic(err)
	}

	return f.Close
}

func expandPath(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}

	return path
}
