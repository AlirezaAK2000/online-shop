package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/AlirezaAK2000/online-shop/initializers"
	"github.com/AlirezaAK2000/online-shop/repo"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignInController(c *gin.Context) {

	var body struct {
		Email    string
		Password string
		Address  string
		Phone    string
		Type     string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body.",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to Hash password.",
		})
		return
	}

	var user repo.User = repo.User{
		Email:    body.Email,
		Password: string(hash),
		Address:  body.Address,
		Phone:    body.Phone,
		Type:     body.Type,
	}

	if _, err := user.InsertUser(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.Status(200)

}

func LogInController(c *gin.Context) {

	var body struct {
		Email    string
		Password string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body.",
		})
		return
	}

	user, err := repo.FindUserByEmail(body.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user or password",
		})
		return
	}

	err1 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.Email,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
		"type": user.Type,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err2 := token.SignedString([]byte(os.Getenv("SECRET")))

	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	_, err3 := initializers.RedisDB.SetEx(initializers.RedisCtx, "token:"+user.Email, tokenString, time.Second*60*60*24*30).Result()
	if err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})

}
