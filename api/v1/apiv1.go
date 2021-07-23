package apiv1

import (
	"github.com/gin-gonic/gin"
	"github.com/moos3/gin-rest-api/api/v1/auth"
	"github.com/moos3/gin-rest-api/api/v1/user"
)

// ApplyRoutes attaches /v1 routs
func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	user.ApplyRoutes(v1)
	auth.ApplyRoutes(v1)
}
