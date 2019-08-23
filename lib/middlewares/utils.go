package middlewares

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// RevisionMiddleware - Set Revision version Header aka GIT SHA
func RevisionMiddleware() gin.HandlerFunc {
	// Revision file contents will be only loaded once per process
	data, err := ioutil.ReadFile("REVISION")

	// if we can't read file, just skip to the next request handler
	// this is pretty much a NOOP middleware :)
	if err != nil {
		log.Println("revision middleware error: ", err)

		return func(c *gin.Context) {
			c.Next()
		}
	}
	// Clean up the value since it should contain line breaks
	revision := strings.TrimSpace(string(data))
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Revision", revision)
		c.Next()
	}
}

// RequestIdMiddleware - for searching logs easier
func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
	}
}
