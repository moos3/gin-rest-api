package api

import (
	"github.com/gin-gonic/gin"
	apiv1 "github.com/moos3/gin-rest-api/api/v1.0"
)

func ApplyRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		apiv1.ApplyRoutes(api)
	}
}
