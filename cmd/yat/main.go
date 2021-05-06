package main

import (
	"flag"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thisiserico/yat"
)

func main() {
	debug := flag.Bool("d", false, "use ./yat.log as log output")
	flag.Parse()

	defer prepareLooger(*debug)()

	files := flag.Args()
	stores := make([]yat.Store, 0, len(files))
	for _, file := range files {
		stores = append(stores, yat.NewTomlStore(file))
	}

	model := yat.NewUI(stores...)
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
