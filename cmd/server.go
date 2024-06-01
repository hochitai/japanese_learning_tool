package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/hochitai/jpl/database"
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
		r.ForwardedByClientIP = true
		r.SetTrustedProxies([]string{"127.0.0.1"})
		v1 := r.Group("/v1")
		{
			// Word
			v1.GET("/words", database.GetWords(db))
			v1.POST("/words", database.AddWord(db))
			v1.PUT("/words/:id", database.UpdateWord(db))
			v1.DELETE("/words/:id", database.DeleteWord(db))

			// User
			v1.POST("/users/register", database.AddUser(db))
			v1.POST("/users/login", database.Login(db))
			v1.PUT("/users/:id", database.UpdateWord(db))
			v1.DELETE("/users/:id", database.DeleteWord(db))

		}

		r.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}