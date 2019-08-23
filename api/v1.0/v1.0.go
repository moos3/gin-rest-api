package apiv1

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/moos3/gin-rest-api/api/v1.0/auth"
	user "github.com/moos3/gin-rest-api/api/v1.0/user"
)

var region string

//
// @Summary Ping
// @Description ping healthcheck

// @Success 200 {object}
// @Failure 401 {object} err
// @Failure 500 {object} err

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
		user.ApplyRoutes(v1)
	}
}
