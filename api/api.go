package api

import (
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		apiv1.ApplyRoutes(api)
	}
}
