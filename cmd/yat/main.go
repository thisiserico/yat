package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thisiserico/yat"
)

func main() {
	model := yat.NewUI()

	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
