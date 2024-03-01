// main package
package main

import (
	"fmt"
	"os"
	"time"

	btea "github.com/charmbracelet/bubbletea"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

type model struct {
	Input        string
	StartTime    time.Time
	FinishTime   time.Time
	DisplayTimer bool
	Done         bool
}

func initialModel() model {
	return model{
		Input:        "",
		StartTime:    time.Time{},
		FinishTime:   time.Time{},
		DisplayTimer: false,
	}
}

func main() {
	p := btea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running the program: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() btea.Cmd {
	return nil
}

func prepareFinalStatus(m model) string {
	m.FinishTime = time.Now()
	elapsed := m.FinishTime.Sub(m.StartTime)
	return fmt.Sprintf("\nFinished! Time taken: %v\nResult: %s", elapsed, m.Input)
}

func prepareCurrentStatus(m model) string {
	return fmt.Sprintf("Type the alphabet: a to z\n\n%s", m.Input)
}

func (m model) Update(msg btea.Msg) (btea.Model, btea.Cmd) {
	if !m.DisplayTimer {
		m.DisplayTimer = true
		m.StartTime = time.Now()
	}
	switch msg := msg.(type) {
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "enter":
			return m, btea.Quit
		default:
			tmp := m.Input + msg.String()
			if tmp == alphabet[0:len(tmp)] {
				m.Input = tmp
			}
			if len(tmp) == len(alphabet) {
				m.Done = true
				return m, btea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.Done {
		return prepareFinalStatus(m)
	}
	return prepareCurrentStatus(m)
}