package main 

import (
"os"

"github.com/gin-gonic/gin"
"github.com/joho/godotenv"
"github.com/moos3/gin-rest-api/api"
"github.com/moos3/gin-rest-api/database"
"github.com/moos3/gin-rest-api/lib/middlewares"

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
	app.Use(database.Inject(db))
	app.Use(middlewares.JWTMiddleware())
	api.ApplyRoutes(app)
	app.Run(":" + port)
}