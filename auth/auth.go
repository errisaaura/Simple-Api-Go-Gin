package auth

import (
	"learn-gome/models"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
)

const (
	USER     = "admin"
	PASSWORD = "Password123!"
	SECRET   = "secret"
)

//fitur login
func LoginHandler(c *gin.Context) {
	var user models.Credential

	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}

	if user.Username != USER {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user invalid"})

	} else if user.Password != PASSWORD {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Password Invalid"})
	} else {
		//token
		claim := jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(), //expirednya
			Issuer:    "test",
			IssuedAt:  time.Now().Unix(), //kapan token dibuat
		}

		sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		token, err := sign.SignedString([]byte(SECRET))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"token":   token,
		})

	}
}
