package jwt

import (
	"GoProject/common"
	. "GoProject/model"
	"GoProject/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var jwtSecret = []byte("secret")

func GenerateToken(user Users) string {

	var userRoles = util.ConvertToStrings(user.UserRoles, func(role UserRole) string {
		return role.UserRoleType.Code
	})
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  user.Username,
		"userRoles": userRoles,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

func Authenticate(c *gin.Context) bool {
	tokenString := c.GetHeader("Authorization")
	idToken := strings.TrimSpace(strings.Replace(tokenString, "Bearer", "", 1))

	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return false
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["userID"].(string)
	//userRoles := claims["userRoles"].([]string)
	userRolesInterface := claims["userRoles"].([]interface{})
	var userRoles []string
	for _, roleInterface := range userRolesInterface {
		userRoles = append(userRoles, roleInterface.(string))
	}
	c.Set("userID", userID)
	c.Set("userRoles", userRoles)
	return true
}

func IsAdmin(c *gin.Context) bool {
	value, exists := c.Get("userRoles")
	if exists {
		userRoles, _ := value.([]string) // Type assertion here
		isAdmin := util.ContainsString(userRoles, common.USER_ROLE_ADMIN)
		if !isAdmin {
			c.JSON(http.StatusOK, common.Fail("user has no authorization"))
			c.Abort()
			return false
		}
	} else {
		c.JSON(http.StatusOK, common.Fail("user has no authorization"))
		c.Abort()
		return false
	}
	return true
}
