package handler

import (
	"GoProject/common"
	"GoProject/model"
	"GoProject/module/jwt"
	"GoProject/module/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// binging validation: https://blog.csdn.net/weixin_40022980/article/details/122796567
// https://medium.com/@maktoobgar/how-to-validate-api-inputs-in-gin-f2af4a3ce43e

type UserHttpHandler struct {
	svc user.Service
}

func NewUserHttpHandler(svc user.Service, server *gin.Engine) UserHttpHandler {
	var handler = UserHttpHandler{svc: svc}

	// find one user
	server.GET("/user/:username", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.findUser(c)
	})

	// find users
	server.GET("/users", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.findUsers(c)
	})

	//create user
	server.POST("/user", func(c *gin.Context) {
		log.Printf("[%s]%s with request body %s", c.Request.Method, c.Request.URL, c.Request.Body)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.createUsers(c)
	})

	//find user roles
	server.GET("/user_roles", func(c *gin.Context) {
		log.Printf("[%s]%s", c.Request.Method, c.Request.URL)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)

	}, func(c *gin.Context) {
		handler.findUserRoles(c)
	})

	//delete user
	server.DELETE("/user/:username", func(c *gin.Context) {
		log.Printf("[%s]%s with request body %s", c.Request.Method, c.Request.URL, c.Request.Body)
		//_ = jwt.Authenticate(c) && jwt.IsAdmin(c)
	}, func(c *gin.Context) {
		handler.deleteUser(c)
	})

	//update user email
	server.PUT("/user/:username", func(c *gin.Context) {
		log.Printf("[%s]%s with request body %s", c.Request.Method, c.Request.URL, c.Request.Body)
		//jwt.Authenticate(c)
	}, func(c *gin.Context) {
		handler.updateUserEmail(c)
	})

	return handler
}

func (g *UserHttpHandler) findUser(c *gin.Context) {
	var username = c.Param("username")
	users, err := g.svc.FindUser(username)
	if err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
	} else {
		c.JSON(http.StatusOK, common.Success(users))
	}
}

func (g *UserHttpHandler) findUsers(c *gin.Context) {
	var userType = c.DefaultQuery("role_type", "USER")
	users, err := g.svc.FindUserByUserType(userType)
	if err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
	} else {
		c.JSON(http.StatusOK, common.Success(users))
	}
}

type CreateUserRequestBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (g *UserHttpHandler) createUsers(c *gin.Context) {
	var requestBody CreateUserRequestBody

	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
		return
	}

	userRoleType, userRoleErr := g.svc.FindUserRoleTypeByCode(common.USER_ROLE_USER)
	if userRoleErr != nil {
		c.JSON(http.StatusOK, common.Fail(userRoleErr.Error()))
		return
	}

	encryptedPassword, _ := getEncryptedPwd(requestBody.Password)
	users := model.Users{
		Username: requestBody.Username,
		Password: encryptedPassword,
		Email:    requestBody.Email,
		UserType: common.USER_ROLE_USER,
	}

	userErr := g.svc.CreateUserData(&users, userRoleType)
	if userErr != nil {
		c.JSON(http.StatusOK, common.Fail(userErr.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(jwt.GenerateToken(users)))
}

func (g *UserHttpHandler) findUserRoles(c *gin.Context) {
	var username = c.DefaultQuery("username", "")
	var userRoles, err = g.svc.FindUserRoleTypeByUsername(username)
	if err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Success(userRoles))
}

func (g *UserHttpHandler) deleteUser(c *gin.Context) {
	var username = c.Param("username")
	users, err := g.svc.FindUser(username)
	if err != nil {
		c.JSON(http.StatusOK, common.Fail(err.Error()))
		return
	}

	deleteErr := g.svc.DeleteUser(users)
	if deleteErr != nil {
		c.JSON(http.StatusOK, common.Fail(deleteErr.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(""))
}

type UpdateUserRequestBody struct {
	Email string
}

func (g *UserHttpHandler) updateUserEmail(c *gin.Context) {
	var username = c.Param("username")
	var requestBody UpdateUserRequestBody

	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		data := common.Fail(err.Error())
		c.JSON(http.StatusOK, data)
		return
	}
	if requestBody.Email == "" {
		c.JSON(http.StatusOK, common.Fail("No Email Input"))
		return
	}

	users, userErr := g.svc.FindUser(username)
	if userErr != nil {
		c.JSON(http.StatusOK, common.Fail(userErr.Error()))
		return
	}

	datas := map[string]interface{}{
		"Email": requestBody.Email,
	}
	users, updateErr := g.svc.UpdateUser(users, datas)
	if updateErr != nil {
		c.JSON(http.StatusOK, common.Fail(updateErr.Error()))
		return
	}

	c.JSON(http.StatusOK, common.Success(users))
}

func getEncryptedPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}
