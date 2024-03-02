// main package
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Alpensin/learn-alphabet/alphabet"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Input            []rune
	StartTime        time.Time
	FinishTime       time.Time
	DisplayTimer     bool
	Mistakes         int
	LastInputMistake string
	Done             bool
	CurPosition      int
	alphabet         []rune
}

func initialModel(alphabet string) model {
	a := []rune(alphabet)
	return model{
		Input:        make([]rune, 0, len(a)),
		StartTime:    time.Time{},
		FinishTime:   time.Time{},
		DisplayTimer: false,
		alphabet:     a,
	}
}

func main() {
	p := tea.NewProgram(initialModel(alphabet.EN))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running the program: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func printFinalStatus(m model) string {
	m.FinishTime = time.Now()
	elapsed := m.FinishTime.Sub(m.StartTime)
	return fmt.Sprintf("Finished! Time taken: %v.\nMistakes made: %d\nFinal result: %s\n", elapsed, m.Mistakes, string(m.Input))
}

func prepareCurrentStatus(m model) string {
	mistakeInfo := ""
	if m.LastInputMistake != "" {
		mistakeInfo = fmt.Sprintf("\n%s not expected. Try again\n", m.LastInputMistake)
	}
	return fmt.Sprintf("type the alphabet letter by letter\n%s%s", string(m.Input), mistakeInfo)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.DisplayTimer {
		m.DisplayTimer = true
		m.StartTime = time.Now()
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "enter":
			return m, tea.Quit
		default:
			tmp := msg.Runes
			if len(tmp) != 1 {
				return m, tea.Printf("unexpected symbols: %s\nExpecting one letter", string(tmp))
			}
			if tmp[0] == m.alphabet[m.CurPosition] {
				m.LastInputMistake = ""
				m.Input = append(m.Input, tmp[0])
				m.CurPosition++
			} else {
				m.LastInputMistake = msg.String()
				m.Mistakes++
			}
			if len(m.Input) == len(m.alphabet) {
				m.Done = true
				return m, tea.Quit
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
