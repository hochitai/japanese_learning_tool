package cmd

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hochitai/jpl/api/handler"
	"github.com/hochitai/jpl/api/middleware"
	"github.com/hochitai/jpl/internal/database"
	"github.com/spf13/cobra"
)

// practiceCmd represents the practice command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run web application server",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.ConnectDB()
		if err != nil {
			fmt.Printf("alas, there's been an error: %v", err)
			os.Exit(1)
		}

		r := gin.Default()
		r.ForwardedByClientIP = true
		r.SetTrustedProxies([]string{"127.0.0.1"})
		r.LoadHTMLGlob("web/templates/*")

		config := cors.DefaultConfig()
  		config.AllowOrigins = []string{"http://localhost:8080"}
		r.Use(cors.Default())

		// pages
		r.GET("/", handler.HomePage(db))

		// v1 rest api
		v1 := r.Group("/v1")
		{
			// Word
			v1.GET("/words", handler.GetWords(db))
			v1.POST("/words", handler.AddWord(db))
			v1.PUT("/words/:id", handler.UpdateWord(db))
			v1.DELETE("/words/:id", handler.DeleteWord(db))

			// User
			v1.POST("/users/register", handler.AddUser(db))
			v1.POST("/users/admin", middleware.CheckTokenAndPermission(), handler.AddAdmin(db))
			v1.POST("/users/login", handler.Login(db))
			v1.POST("/users/token", middleware.RefreshToken())
			v1.PUT("/users/:id", middleware.CheckToken(), handler.UpdateUser(db))
			v1.DELETE("/users/:id", middleware.CheckTokenAndPermission(), handler.DeleteUser(db))
			v1.POST("/users/verify", handler.Verify(db))

			//Admin
			v1.GET("/admin/users", middleware.CheckTokenAndPermission(), handler.GetUsers(db))

		}

		r.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}