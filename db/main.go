package db

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Word struct {
	Characters string
	Sound      string
}

type WordModel struct {
	TextInput textinput.Model
	Words     []Word
	Word      Word
	Result    string
}

func InitialModel(words []Word, word Word) WordModel {
	ti := textinput.New()
	ti.Placeholder = "Type a Romaji representation and press enter..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return WordModel{
		TextInput: ti,
		Words: words,
		Word: word,
	}
}

func (m WordModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m WordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if strings.Compare(m.Word.Sound, m.TextInput.Value()) == 0 {
				m.Result = styleCorrect.Render("Correct!")
			} else {
				m.Result = styleFail.Render("Fail! Correct is: " + m.Word.Sound)
			}
			m.TextInput.Reset()
			m.Word = GetRandomCharacter(m.Words)
		default:
			m.TextInput, cmd = m.TextInput.Update(msg)
			return m, cmd
		}
	}
	return m, cmd
}

func (m WordModel) View() string {

	var style = lipgloss.NewStyle().
						Bold(true).
						Foreground(lipgloss.Color("#66b3ff"))

	return fmt.Sprintf(
		"Test your memory, Alphabet Exam!\n\n %s %s\n\n%s\n\n%s\n\n%s",
		style.Render("Word:"),
		m.Word.Characters,
		m.TextInput.View(),
		m.Result,
		"(esc to quit)",
	) + "\n"
}

