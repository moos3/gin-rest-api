package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/moos3/gin-rest-api/database/models"
	"github.com/moos3/gin-rest-api/lib/common"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"

)

var gdprCountries = []string{"AT", "BE", "HR", "BG", "CY", "CZ", "DK", "EE", "FI", "FR", "DE", "GR", "HU", "IE",
	"IT", "LV", "LT", "LU", "MT", "NL", "PL", "PT", "RO", "SK", "SI", "ES", "SE", "GB"}

// User is alias for models.User
type User = models.User

func hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJwtToken(data common.JSON) (string, error) {
	region := os.Getenv("REGION")
	date := time.Now().Add(time.Hour * 24 * 7)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":        data,
		"exp":         date.Unix(),
		"host_region": region,
	})

	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"
	key, readErr := ioutil.ReadFile(keyPath)
	if readErr != nil {
		return "", readErr
	}
	tokenString, err := token.SignedString(key)
	return tokenString, err
}

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

	hash, hashErr := hash(body.Password)
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
	token, _ := generateJwtToken(serialized)
	c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)

	c.JSON(200, common.JSON{
		"user":   user.Serialize(),
		"token":  token,
		"region": user.Region,
	})
}

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
	if !checkHash(body.OldPasword, user.PasswordHash) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	hash, hashErr := hash(body.NewPassword)
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

	if !checkHash(body.Password, user.PasswordHash) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	serialized := user.Serialize()
	token, _ := generateJwtToken(serialized)

	c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, common.JSON{
		"user":  user.Serialize(),
		"token": token,
	})
}

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
		token, _ := generateJwtToken(user.Serialize())
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

// This is Forgot password functionality 
func forgotPassword(c *gin.Context){
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
			"message": "username not found"
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
		Token: id.String(),
		UserID: user.ID,
		RequestedByIP: c.ClientIP(),
		Claimed: false,
		Expiration: timesUp.Unix(),
	}

	db.NewRecord(&r)
	if err := db.Create(&r).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.JSON{
			"success": false,
			"message": err,
		})
	}

	c.JSON(http.StatusOK, common.JSON{
		"reset_token": r.Token,
		"claimed": r.Claimed,
		"expiration": r.Expiration,
		"requested_by": r.RequestedByIP,
	})
}

func resetPassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Username string `json:"username" binding:"required"`
		Token string `json:"reset_token" binding:"required"`
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
			"message": "username not found"
		})
		return
	}

	var resetToken ResetPasswordToken
	if err = db.Where("token = ?", body.Token).First(&resetToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, common.JSON{
			"success": false,
			"message": "token not found",
		})
	}

	// See if token has been already used
	if resetToken.Claimed {
		c.JSON(http.StatusBadRequest, common.JSON{
			"success": false,
			"message": "Token has already been claimed"
		})
	}

	// check if token as expired
	now := time.Now()


}

// TODO: locateMe a function to see if the user is here, if not locate the user in the other
// regions

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
			"message": "username not found"
		})
		return
	}

	
}