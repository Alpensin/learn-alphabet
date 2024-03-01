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
	Input            string
	StartTime        time.Time
	FinishTime       time.Time
	DisplayTimer     bool
	Mistakes         int
	LastInputMistake string
	Done             bool
	CurPosition      int
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

func printFinalStatus(m model) string {
	m.FinishTime = time.Now()
	elapsed := m.FinishTime.Sub(m.StartTime)
	return fmt.Sprintf("Finished! Time taken: %v.\nMistakes made: %d\nFinal result: %q\n", elapsed, m.Mistakes, m.Input)
}

func prepareCurrentStatus(m model) string {
	mistakeInfo := ""
	if m.LastInputMistake != "" {
		mistakeInfo = fmt.Sprintf("\n%s not expected. Try again\n", m.LastInputMistake)
	}
	return fmt.Sprintf("Type the alphabet: a to z\n%s%s", m.Input, mistakeInfo)
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
			tmp := msg.String()
			if tmp == string(alphabet[m.CurPosition]) {
				m.LastInputMistake = ""
				m.Input = m.Input + tmp
				m.CurPosition++
			} else {
				m.LastInputMistake = msg.String()
				m.Mistakes++
			}
			if len(m.Input) == len(alphabet) {
				m.Done = true
				return m, btea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.Done {
		return printFinalStatus(m)
	}
	return prepareCurrentStatus(m)
}
