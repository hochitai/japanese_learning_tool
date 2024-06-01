package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hochitai/japanese_learning_tool/database"
	"github.com/spf13/cobra"
)

var practiceCmd = &cobra.Command{
	Use:   "practice",
	Short: "Practice japanese vocabulary",
	Run: func(cmd *cobra.Command, args []string) {
		db := database.ConnectDB()
		defer db.Close()

		words, err := database.GetVocabularies(db)
		if err != nil {
			panic(err)
		}
		word := database.GetRandomCharacter(words)

		p := tea.NewProgram(database.InitialModel(words, word))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(practiceCmd)
}