package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thisiserico/yat"
)

func main() {
	store := yat.NewFileStore(expandPath("~/.yat"))
	model := yat.NewUI(store.LoadTasks())

	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func expandPath(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}

	return path
}
