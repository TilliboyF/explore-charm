package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
