package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hochitai/japanese_learning_tool/database"
	"github.com/spf13/cobra"
)

var hiraganaCmd = &cobra.Command{
	Use:   "hiragana",
	Short: "Practice hiragana character",
	Run: func(cmd *cobra.Command, args []string) {
		words := database.GetHiraganaCharacters()
		word := database.GetRandomCharacter(words)
		p := tea.NewProgram(database.InitialModel(words, word))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(hiraganaCmd)
}