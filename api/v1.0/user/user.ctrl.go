package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/moos3/gin-rest-api/database/models"
	"github.com/moos3/gin-rest-api/lib/common"
)

// User is alias for models.User
type User = models.User

// UserProfile is alias for models.UserProfile
type UserProfile = models.UserProfile

//
// @Summary Fetch Profile for a given user
// @Description get profile by username or user id
// @Accept json
// @Produce json
// @Param some_id path string true "Some ID"
// @Success 200 {object} common.JSON
// @Failure 404 {object} err
// @Failure 400 {object} err

func fetchProfile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	username := c.Param("username")
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusBadRequest) // user not found
		return
	}

	var userProfile UserProfile
	if err := db.Where("user_id = ?", user.ID).First(&userProfile).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound) // user profile not found
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, common.JSON{
		"username":     user.Username,
		"uid":          user.ID,
		"display_name": user.DisplayName,
		"bio":          userProfile.Bio,
		"twitter":      userProfile.TwitterHandle,
	})

}

//
// @Summary Fetch Profile for a given user
// @Description get profile by username or user id
// @Accept json
// @Produce json
// @Param some_id path string true "Some ID"
// @Success 200 {object} common.JSON
// @Failure 404 {object} err
// @Failure 400 {object} err

func profile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	userRaw, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user := userRaw.(User)
	if err := db.Where("username = ?", user.Username).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound) // user not found
		return
	}

	var userProfile UserProfile
	if err := db.Where("user_id = ?", user.ID).First(&userProfile).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound) // user profile not found
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, common.JSON{
		"username":     user.Username,
		"uid":          user.ID,
		"display_name": user.DisplayName,
		"bio":          userProfile.Bio,
		"twitter":      userProfile.TwitterHandle,
	})

}

//
// @Summary Edit Profile for a given user
// @Description edit profile by username or user id
// @Accept json
// @Produce json
// @Param some_id path string true "Some ID"
// @Success 200 {object} common.JSON
// @Failure 404 {object} err
// @Failure 400 {object} err
// @Failure 401 {object} err

func edit(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userRaw, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user := userRaw.(User)
	if err := db.Where("username = ?", user.Username).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound) // user not found
		return
	}

	type RequestBody struct {
		Bio           string `json:"bio"`
		TwitterHandle string `json:"twitter"`
		Email         string `json:"email"`
		Region        string `json:"region"`
		DisplayName   string `json:"display_name"`
	}

	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db.First(&user)
	if body.Email != "" {
		user.Email = body.Email
	}
	if body.DisplayName != "" {
		user.DisplayName = body.DisplayName
	}
	if body.Region != "" {
		user.Region = body.Region
	}
	err := db.Save(&user).Error
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, common.JSON{
			"message": "Profile failed to be updated",
			"action":  true,
		})
		return
	}

	var userProfile UserProfile
	db.First(&userProfile)
	userProfile.UserID = user.ID
	if body.Bio != "" {
		userProfile.Bio = body.Bio
	}
	if body.TwitterHandle != "" {
		userProfile.TwitterHandle = body.TwitterHandle
	}
	err = db.Save(&userProfile).Error
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, common.JSON{
			"message": "profile failed to be updated",
			"action":  true,
		})
		return
	}

	c.JSON(http.StatusOK, common.JSON{
		"message": "profile updated",
		"action":  true,
	})
}

//
// @Summary Fetch Region for a given user
// @Description get the region of a user
// @Accept json
// @Produce json
// @Param some_id path string true "Some ID"
// @Success 200 {object} common.JSON
// @Failure 404 {object} err
// @Failure 400 {object} err
// @Failure 401 {object} err

func getMyRegion(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userRaw, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user := userRaw.(User)
	if err := db.Where("username = ?", user.Username).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound) // user not found
		return
	}

	c.JSON(http.StatusOK, common.JSON{
		"username": user.Username,
		"region":   user.Region,
		"uuid":     user.ID,
	})
}
