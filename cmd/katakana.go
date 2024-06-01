package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hochitai/jpl/database"
	"github.com/spf13/cobra"
)

var katakanaCmd = &cobra.Command{
	Use:   "katakana",
	Short: "Practice katakana character",
	Run: func(cmd *cobra.Command, args []string) {
		words := database.GetKatakanaCharacters()
		word := database.GetRandomCharacter(words)
		p := tea.NewProgram(database.InitialModel(words, word))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(katakanaCmd)
}



