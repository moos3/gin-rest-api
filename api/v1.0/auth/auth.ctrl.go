package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/moos3/gin-rest-api/database/models"
	"github.com/moos3/gin-rest-api/lib/common"
)

// User is alias for models.User
type User = models.User

// ResetPasswordToken is alias for models.ResetPasswordToken
type ResetPasswordToken = models.ResetPasswordToken

//
// @Summary Register User
// @Description user sign up
// @Accept json
// @Produce json
// @Param
// @Success 200 {object} common.JSON
// @Failure 404 {object} err
// @Failure 400 {object} err

func register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Username    string `json:"username" binding:"required"`
		DisplayName string `json:"display_name" binding:"required"`
		Password    string `json:"password" binding:"required"`
		Region      string `json:"region" binding:"required"`
	}

	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var exists User
	if err := db.Where("username = ?", body.Username).First(&exists).Error; err == nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	hash, hashErr := common.Hash(body.Password)
	if hashErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// validate country
	// query := gountries.New()
	// data, _ := query.FindCountryByName(body.Region)
	// check the country code here

	user := User{
		Username:     body.Username,
		DisplayName:  body.DisplayName,
		PasswordHash: hash,
		Region:       body.Region,
	}

	db.NewRecord(user)
	db.Create(&user)

	serialized := user.Serialize()
	token, _ := common.GenerateJwtToken(serialized)
	c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)

	c.JSON(200, common.JSON{
		"user":   user.Serialize(),
		"token":  token,
		"region": user.Region,
	})
}

//
// @Summary Change User password
// @Description update the users password
// @Accept json
// @Produce json
// @Param
// @Success 200 {object} common.JSON
// @Failure 404 {object} err
// @Failure 400 {object} err
// @Failure 500 {object} err

func changePassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	type RequestBody struct {
		Username    string `json:"username" binding:"required"`
		OldPasword  string `json:"old_pasword" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check existancy
	var user User
	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound) // user not found
		return
	}

	// check old password
	if !common.CheckHash(body.OldPasword, user.PasswordHash) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	hash, hashErr := common.Hash(body.NewPassword)
	if hashErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db.First(&user)
	user.PasswordHash = hash
	err := db.Save(&user).Error
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, common.JSON{
		"message": "password changed",
		"action":  true,
	})

}


//
// @Summary Login User password
// @Description login user and get JWT token
// @Accept json
// @Produce json
// @Param
// @Success 200 {object} common.JSON
// @Failure 404 {object} err
// @Failure 400 {object} err

func login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	type RequestBody struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check existancy
	var user User
	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound) // user not found
		return
	}

	if !common.CheckHash(body.Password, user.PasswordHash) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	serialized := user.Serialize()
	token, _ := common.GenerateJwtToken(serialized)

	c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, common.JSON{
		"user":  user.Serialize(),
		"token": token,
	})
}

//
// @Summary Check username name availablity
// @Description see if username has already been taken
// @Accept json
// @Produce json
// @Param
// @Success 200 {object} common.JSON
// @Failure 404 {object} err
// @Failure 400 {object} err

func usernameAvailability(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Username string `json:"username" binding:"required"`
	}
	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check existancy
	var user User
	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		// user not found
		c.JSON(http.StatusOK, common.JSON{
			"username":    "available",
			"host_region": user.Region,
		})
		return
	}

	c.JSON(http.StatusOK, common.JSON{
		"username":    "not available",
		"host_region": user.Region,
	})
}

//
// @Summary Check token
// @Description renew token when token life is less than 3 days, otherwise, return null for token
// @Accept json
// @Produce json
// @Param
// @Success 200 {object} common.JSON
// @Failure 401 {object} err

// check API will renew token when token life is less than 3 days, otherwise, return null for token
func check(c *gin.Context) {
	userRaw, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user := userRaw.(User)

	tokenExpire := int64(c.MustGet("token_expire").(float64))
	now := time.Now().Unix()
	diff := tokenExpire - now

	fmt.Println(diff)
	if diff < 60*60*24*3 {
		// renew token
		token, _ := common.GenerateJwtToken(user.Serialize())
		c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)
		c.JSON(http.StatusOK, common.JSON{
			"token":  token,
			"user":   user.Serialize(),
			"region": user.Region,
		})
		return
	}

	c.JSON(http.StatusOK, common.JSON{
		"token":  nil,
		"user":   user.Serialize(),
		"regoin": user.Region,
	})
}
//
// @Summary Forgot Password
// @Description generates a reset password token
// @Accept json
// @Produce json
// @Param
// @Success 200 {object} common.JSON
// @Failure 401 {object} err
// @Failure 500 {object} err

// This is Forgot password functionality
func forgotPassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Username string `json:"username" binding:"required"`
	}

	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// check existancy
	var user User
	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		// user not found
		c.JSON(http.StatusBadRequest, common.JSON{
			"success": false,
			"message": "username not found",
		})
		return
	}
	now := time.Now()
	expireTime := time.Minute * 30
	timesUp := now.Add(expireTime)

	// generate one time token
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("Failed to make token")
	}

	r := ResetPasswordToken{
		Token:         id.String(),
		UserID:        user.ID,
		RequestedByIP: c.ClientIP(),
		Claimed:       false,
		Expiration:    timesUp.Unix(),
	}

	db.NewRecord(&r)
	if err := db.Create(&r).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.JSON{
			"success": false,
			"message": err,
		})
	}

	c.JSON(http.StatusOK, common.JSON{
		"reset_token":  r.Token,
		"claimed":      r.Claimed,
		"expiration":   r.Expiration,
		"requested_by": r.RequestedByIP,
	})
}

//
// @Summary Reset Password
// @Description Takes Reset Password Token and updates given password
// @Accept json
// @Produce json
// @Param
// @Success 200 {object} common.JSON
// @Failure 401 {object} err

func resetPassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Username string `json:"username" binding:"required"`
		Token    string `json:"reset_token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// check existancy
	var user User
	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		// user not found
		c.JSON(http.StatusBadRequest, common.JSON{
			"success": false,
			"message": "username not found",
		})
		return
	}

	var resetToken ResetPasswordToken
	if err := db.Where("token = ?", body.Token).First(&resetToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, common.JSON{
			"success": false,
			"message": "token not found",
		})
	}

	// See if token has been already used
	if resetToken.Claimed {
		c.JSON(http.StatusBadRequest, common.JSON{
			"success": false,
			"message": "Token has already been claimed",
		})
	}

	// check if token as expired
	//now := time.Now()

}

// TODO: locateMe a function to see if the user is here, if not locate the user in the other
// regions
//
// @Summary Locates User
// @Description checks server to see if user is here
// @Accept json
// @Produce json
// @Param
// @Success 200 {object} common.JSON
// @Failure 401 {object} err
// @Failure 500 {object} err

func locateMe(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Username string `json:"username" binding:"required"`
	}

	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// check existancy
	var user User
	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		// user not found
		c.JSON(http.StatusBadRequest, common.JSON{
			"success": false,
			"message": "username not found",
		})
		return
	}

}
