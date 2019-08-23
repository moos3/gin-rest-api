package main 

import (
	"os"
	//"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/moos3/gin-rest-api/api"
	"github.com/moos3/gin-rest-api/database"
	"github.com/moos3/gin-rest-api/lib/middlewares"
	"github.com/gin-contrib/cors"

	"github.com/gin-contrib/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"github.com/moos3/gin-rest-api/docs"

)

var region string

// @title Simple Gin Rest API with JWT
// @version 1.0
// @description this is a simple web application that supports api versioning, jwt and multi-regional deployments

// @contact.name Moos3
// @contact.url https://github.com/moos3/gin-rest-api/issues
// @contact.email github@guthnur.net

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host myapp.guthnur.net
// @BasePath /api/v1/

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	region = os.Getenv("REGION")

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Gin Rest API"
	docs.SwaggerInfo.Description = "This is a simple gin rest api with authentication."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "api.guthnur.net"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	db, _ := database.Initialize()

	port := os.Getenv("PORT")
	app := gin.Default()

	// Add a logger middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.	
	app.Use(logger.SetLogger())

	// Custom logger
	subLog := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Logger()

	app.Use(logger.SetLogger(logger.Config{
		Logger:         &subLog,
		UTC:            true,
		SkipPath:       []string{"/skip"},
	}))

	/*
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://localhost"},
		AllowMethods: []string{"PUT", "PATCH", "GET"},
		AllowHeaders: []string{"Origin"},
		ExposeHeaders: []string{"Content-Lenght"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	*/
	app.Use(cors.Default())
	app.Use(database.Inject(db))
	app.Use(middlewares.RequestIdMiddleware())
	app.Use(middlewares.RevisionMiddleware())
	app.Use(middlewares.JWTMiddleware())
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api.ApplyRoutes(app)

	app.Run(":" + port)
}