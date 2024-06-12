package cmd

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hochitai/jpl/api/handler"
	"github.com/hochitai/jpl/api/middleware"
	"github.com/hochitai/jpl/docs"
	"github.com/hochitai/jpl/internal/database"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

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

		// programmatically set swagger info
		docs.SwaggerInfo.Title = "Swagger Example API"
		docs.SwaggerInfo.Description = "This is a sample server Petstore server."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = "petstore.swagger.io"
		docs.SwaggerInfo.BasePath = "/v2"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}

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
			v1.POST("/words", middleware.VerifyToken(), handler.AddWord(db))
			v1.PUT("/words/:id", middleware.VerifyToken(), handler.UpdateWord(db))
			v1.DELETE("/words/:id", middleware.VerifyToken(), handler.DeleteWord(db))

			// User
			v1.POST("/users/register", handler.AddUser(db))
			v1.POST("/users/admin", middleware.VerifyTokenAndPermission(), handler.AddAdmin(db))
			v1.POST("/users/login", handler.Login(db))
			v1.POST("/users/token", middleware.RefreshToken())
			v1.PUT("/users/:id", middleware.VerifyToken(), handler.UpdateUser(db))
			v1.DELETE("/users/:id", middleware.VerifyTokenAndPermission(), handler.DeleteUser(db))

			//Admin
			v1.GET("/admin/users", middleware.VerifyTokenAndPermission(), handler.GetUsers(db))

		}
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		r.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}