package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hochitai/jpl/api/handler"
	"github.com/hochitai/jpl/internal/model"
	"github.com/spf13/cobra"
)

var hiraganaCmd = &cobra.Command{
	Use:   "hiragana",
	Short: "Practice hiragana character",
	Run: func(cmd *cobra.Command, args []string) {
		words := model.GetHiraganaCharacters()
		word := model.GetRandomCharacter(words)
		p := tea.NewProgram(handler.InitialModel(words, word))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(hiraganaCmd)
}