package user

import (
	"github.com/gin-gonic/gin"
)

// ApplyRoutes - User functions
func ApplyRoutes(r *gin.RouterGroup) {
	user := r.Group("/users")
	{
		user.GET("/me", profile)
		user.POST("/profile/edit", edit)
		user.GET("/fetch/region", getMyRegion)
		user.GET("/profile/:username", fetchProfile)
	}
}
