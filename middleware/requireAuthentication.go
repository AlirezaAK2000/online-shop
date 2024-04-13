package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AlirezaAK2000/online-shop/initializers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func RequireAuth(c *gin.Context) {

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		log.Println(err)
		return
	}
	token, err1 := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err1 != nil {
		log.Println(err1)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Using MongoDB to validate token
		// user, err2 := repo.FindUserByID(claims["sub"].(string))

		// if err2 != nil {
		// 	log.Println(err2)
		// 	c.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }

		// c.Set("user", user)

		cachedToken, err2 := initializers.RedisDB.Get(initializers.RedisCtx, "token:"+claims["sub"].(string)).Result()
		if err2 != nil {
			log.Println(err2)
			if err2 == redis.Nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if cachedToken != tokenString {
			log.Println("does not match the cached token")
			log.Println("cached token: ", cachedToken)
			// log.Println("token : ", tokenString)

			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

}
