package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtkey = []byte(get("JWT_SECRET", "my_secret"))

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

type Claims struct {
	UserID uint `json:"uid"`
	jwt.RegisteredClaims
}

func GenerateToken(uid uint) (string, error) { // fn
	claims := &Claims{ // claims reference
		UserID: uid, // assign uid
		RegisteredClaims: jwt.RegisteredClaims{ // time of expire and issue
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // add 24 hour in now time
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // now time
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // encode jwt
	return token.SignedString(jwtkey)                          // sign using key
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization") // geting authorization header
		parts := strings.SplitN(authz, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		tkn, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(t *jwt.Token) (any, error) {
			return jwtkey, nil
		})
		if err != nil || !tkn.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims := tkn.Claims.(*Claims)
		c.Set("uid", claims.UserID)
		c.Next()

	}
}

func UID(c *gin.Context) uint {
	v, ok := c.Get("uid")
	if !ok {
		return 0
	}
	return v.(uint)
}
