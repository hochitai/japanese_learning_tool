package database

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  password = "123456789"
  dbname   = "learning"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
						host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

type WordCmdModel struct {
	TextInput textinput.Model
	Words     []Word
	Word      Word
	Result    string
}

func InitialModel(words []Word, word Word) WordCmdModel {
	ti := textinput.New()
	ti.Placeholder = "Type a Romaji representation and press enter..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return WordCmdModel{
		TextInput: ti,
		Words: words,
		Word: word,
	}
}

func (m WordCmdModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m WordCmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	styleCorrect := lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#1aff66"))

	styleFail := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#ff4d4d"))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if strings.Compare(m.Word.Pronunciation, m.TextInput.Value()) == 0 {
				m.Result = styleCorrect.Render("Correct! " + m.Word.Meaning)
			} else {
				m.Result = styleFail.Render("Fail! Correct is: " + m.Word.Pronunciation + " " + m.Word.Meaning)
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

func (m WordCmdModel) View() string {

	style := lipgloss.NewStyle().
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

