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

)

var region string

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	region = os.Getenv("REGION")

	db, _ := database.Initialize()

	port := os.Getenv("PORT")
	app := gin.Default()
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
	app.Use(middlewares.JWTMiddleware())
	api.ApplyRoutes(app)
	app.Run(":" + port)
}