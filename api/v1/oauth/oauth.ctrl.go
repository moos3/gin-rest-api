package oauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/moos3/gin-rest-api/database/models"
	"github.com/moos3/gin-rest-api/lib/common"
)

// User is alias for models.User
type User = models.User

type OauthToken = models.OauthToken

// UserProfile is alias for models.UserProfile
type UserProfile = models.UserProfile


// TODO: https://github.com/RichardKnop/go-oauth2-server/blob/master/models/oauth.go


func getToken(c *gin.Context) {
	userRaw, ok := c.Get("user")
	db := c.MustGet("db").(*gorm.DB)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}


	var oauthTokens OauthToken
	var user User
	if err := db.Where("username = ?", userRaw.(User).Username).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound) // user not found
		return
	}
	if err := db.Where("user_id = ?",user.ID).First(&oauthTokens).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound) // token not found
		return
	}

	c.JSON(http.StatusOK, common.JSON{
		"oauthtokens":oauthTokens,
	})
}


func postToken(c *gin.Context) {

}

func getAuthorize(c *gin.Context) {

}


func postAuthorize(c *gin.Context) {

}


func getRevoke(c *gin.Context) {

}

func postRevoke(c *gin.Context) {

}