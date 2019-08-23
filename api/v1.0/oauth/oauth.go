package oauth



import (
	"github.com/gin-gonic/gin"
)

// TODO: Add GET/POST routes for method via gin library
func ApplyRoutes(r *gin.RouterGroup) {
	oauth := r.Group("/oauth")
	{
		oauth.GET("/token",getToken)

		oauth.GET("/authorize",getAuthorize)
		oauth.POST("/authorize",postAuthorize)

		oauth.GET("/revoke",getRevoke)
		oauth.POST("/revoke",postRevoke)

	}
}

