/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hochitai/japanese_learning_tool/db"
	"github.com/spf13/cobra"
)

// hiraganaCmd represents the hiragana command
var hiraganaCmd = &cobra.Command{
	Use:   "hiragana",
	Short: "Learn japanses character in hiragana",
	Run: func(cmd *cobra.Command, args []string) {
		// db.HiraganaExam()
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(hiraganaCmd)
}

type model struct {
	textInput 	textinput.Model
	words		[]db.Word
	word 		db.Word
	result		string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type a Romaji representation and press enter..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	words := db.GetCharactersInHiragana()

	return model{
		textInput: ti,
		words: db.GetCharactersInHiragana(),
		word: db.GetWordInHiragana(words),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	var styleCorrect = lipgloss.NewStyle().
						Bold(true).
						Foreground(lipgloss.Color("#1aff66"))

	var styleFail = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#ff4d4d"))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if strings.Compare(m.word.Sound, m.textInput.Value()) == 0 {
				m.result = styleCorrect.Render("Correct!")
			} else {
				m.result = styleFail.Render("Fail! Correct is: " + m.word.Sound)
			}
			m.textInput.Reset()
			m.word = db.GetWordInHiragana(m.words)
		default:
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}
	return m, cmd
}

func (m model) View() string {

	var style = lipgloss.NewStyle().
						Bold(true).
						Foreground(lipgloss.Color("#66b3ff"))

	return fmt.Sprintf(
		"Test your memory, Hiragana Exam!\n\n %s %s\n\n%s\n\n%s\n\n%s",
		style.Render("Word:"),
		m.word.Characters,
		m.textInput.View(),
		m.result,
		"(esc to quit)",
	) + "\n"
}