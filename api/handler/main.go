package handler

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hochitai/jpl/internal/model"
	_ "github.com/lib/pq"
)

type Server struct {}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) DefaultConfiguation() *gin.Engine {
	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.LoadHTMLGlob("web/templates/*")

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"}
	r.Use(cors.Default())
	return r
}



type WordCmdModel struct {
	TextInput textinput.Model
	Words     []model.Word
	Word      model.Word
	Result    string
}

func InitialModel(words []model.Word, word model.Word) WordCmdModel {
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
			m.Word = model.GetRandomCharacter(m.Words)
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

