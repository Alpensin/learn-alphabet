// main package
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Alpensin/learn-alphabet/alphabet"
	tea "github.com/charmbracelet/bubbletea"
)

type timer struct {
	startTime    time.Time
	finishTime   time.Time
	displayTimer bool
}

type mistakes struct {
	count           int
	lastMistakeText string
}

type status struct {
	input           []rune
	done            bool
	currentPosition int
	mistakes        mistakes
}

type model struct {
	timer    timer
	alphabet []rune
	status   status
}

func initialModel(alphabet string) model {
	a := []rune(alphabet)
	return model{
		status: status{
			input: make([]rune, 0, len(a)),
		},
		alphabet: a,
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
	elapsed := m.timer.finishTime.Sub(m.timer.startTime)
	return fmt.Sprintf("Finished! Time taken: %v.\nMistakes made: %d\nFinal result: %s\n", elapsed, m.status.mistakes.count, string(m.status.input))
}

func prepareCurrentStatus(m model) string {
	mistakeInfo := ""
	if m.status.mistakes.lastMistakeText != "" {
		mistakeInfo = fmt.Sprintf("\n%s not expected. Try again\n", m.status.mistakes.lastMistakeText)
	}
	return fmt.Sprintf("type the alphabet letter by letter\n%s%s", string(m.status.input), mistakeInfo)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.timer.displayTimer {
		m.timer.displayTimer = true
		m.timer.startTime = time.Now()
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		default:
			if len(msg.Runes) != 1 {
				return m, tea.Printf("unexpected symbols: %s\nExpecting one letter", string(msg.String()))
			}
			if msg.Runes[0] == m.alphabet[m.status.currentPosition] {
				m.status.input = append(m.status.input, msg.Runes[0])
				m.status.currentPosition++
			} else {
				m.status.mistakes.lastMistakeText = msg.String()
				m.status.mistakes.count++
			}
			if len(m.status.input) == len(m.alphabet) {
				m.timer.finishTime = time.Now()
				m.status.done = true
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.status.done {
		return printFinalStatus(m)
	}
	return prepareCurrentStatus(m)
}
