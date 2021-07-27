package exchange

import (
"github.com/gin-gonic/gin"
)

// ApplyRoutes - User functions
func ApplyRoutes(r *gin.RouterGroup) {
	exchange := r.Group("/exchanges")
	{
		exchange.GET("/ws", wsStream)
	}
}
