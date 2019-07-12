package apiv1

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/moos3/gin-rest-api/api/v1.0/auth"
)

var region string

func ping(c *gin.Context) {
	region = os.Getenv("REGION")

	c.JSON(200, gin.H{
		"message": "pong",
		"region":  region,
	})
}

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1.0")
	{
		v1.GET("/ping", ping)
		auth.ApplyRoutes(v1)
	}
}
