package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hochitai/jpl/api/handler"
	"github.com/hochitai/jpl/internal/database"
	"github.com/hochitai/jpl/internal/model"
	"github.com/spf13/cobra"
)

var practiceCmd = &cobra.Command{
	Use:   "practice",
	Short: "Practice japanese vocabulary",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.ConnectDB()
		if err != nil {
			fmt.Println("alas, there's been an error: %v", err)
			os.Exit(1)
		}

		var wordModel model.Word
		words, err := wordModel.GetVocabularies(db)
		if err != nil {
			panic(err)
		}
		word := model.GetRandomCharacter(words)

		p := tea.NewProgram(handler.InitialModel(words, word))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(practiceCmd)
}
