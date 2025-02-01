package main

import (
	"fmt"
	"log"
	"os"
	"github.com/Software78/encryption-test/docs"
	handler "github.com/Software78/encryption-test/src/controllers"
	db "github.com/Software78/encryption-test/src/db"
	middleware "github.com/Software78/encryption-test/src/middleware"
	"github.com/Software78/encryption-test/src/models"
	repository "github.com/Software78/encryption-test/src/repository"
	service "github.com/Software78/encryption-test/src/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	pagination "github.com/webstradev/gin-pagination"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//	@title						encryption-test API
//	@version					1.0
//	@description				encryption-test API Documentation
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//	@host						localhost:8080
//	@BasePath					/api/v1
//	@securityDefinitions.apiKey	BearerAuth
//	@in							header
//	@name						Authorization

func main() {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = os.Getenv("HOST")
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{os.Getenv("SCHEMES")}
	var err error
	postgresDb := os.Getenv("POSTGRES_URL")
	gormDB, err := gorm.Open(postgres.Open(postgresDb), &gorm.Config{})
	if err != nil {
		log.Fatal("🚨🚨🚨---failed to connect to database---🚨🚨🚨")
		log.Panic(err)
	} else {
		fmt.Println("🚀🚀🚀---ASCENDE SUPERIUS---🚀🚀🚀")
	}
	database := db.NewGormDB(gormDB)
	database.AutoMigrate(&models.User{})
	userRepository := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepository)
	userController := handler.NewUserController(*userService)
	r := gin.Default()
	r.Use(pagination.Default())
	r.Use(middleware.ErrorHandler())

	

	crypto, err := middleware.NewCryptoMiddlewareFromEnv( `/docs/`)
	if err != nil {
		log.Fatal("🚨🚨🚨---failed to create crypto middleware---🚨🚨🚨")
		fmt.Println(err)
		log.Panic(err)
	}
    r.Use(crypto.DecryptRequestMiddleware())
	// r.Use(crypto.EncryptResponseMiddleware())


	//v1 group
	v1 := r.Group("/api/v1")
	v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
 // Exclude Swaggo documentation paths (docs/*any)

	//auth group
	auth := v1.Group("/auth")
	auth.POST("/login", userController.Login)
	auth.POST("/register", userController.Register)

	r.Run(":8080")
}
