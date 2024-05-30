package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/hochitai/japanese_learning_tool/database"
	"github.com/spf13/cobra"
)

// practiceCmd represents the practice command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run web application server",
	Run: func(cmd *cobra.Command, args []string) {
		db := database.ConnectDB()
		defer db.Close()

		r := gin.Default()
		v1 := r.Group("/v1")
		{
			v1.GET("/words", database.GetWords(db))
			v1.POST("/word", database.AddWord(db))
			v1.PUT("/word", database.UpdateWord(db))
			v1.DELETE("/word", database.DeleteWord(db))
		}

		r.Run()

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}