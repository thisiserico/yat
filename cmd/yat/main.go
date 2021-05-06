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

func main() {
	debug := flag.Bool("d", false, "use ./yat.log as log output")
	flag.Parse()

	defer prepareLooger(*debug)()

	files := flag.Args()
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
