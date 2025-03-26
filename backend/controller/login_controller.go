package controller

import (
	"fmt"
	"math/rand"
	"net/http"

	"rental-property-management-system/backend/config"
	"rental-property-management-system/backend/store"
	"rental-property-management-system/backend/utils"

	"github.com/gin-gonic/gin"
)

// 登录接口
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var reqeust LoginRequest
	if err := c.ShouldBindJSON(&reqeust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	// 检查用户名是否存在
	var user store.User
	if err := store.GetDB().Where("username = ?", reqeust.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("User %q does not exist", reqeust.Username)})
		c.Abort()
		return
	}

	// 验证密码是否正确并生成
	if user.PasswordHash != utils.Sha256(reqeust.Password, user.Salt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong username or password"})
		c.Abort()
		return
	}

	// 生成 access token
	token, err := utils.GenerateAccessToken(int(user.ID), config.PrivateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		c.Abort()
		return
	}

	// TODO: 如果需要登陆状态的话
	// sessionId := uuid.NewString()
	// if err := store.SetSession(ctx, strconv.Itoa(user.Id), []byte(sessionId)); err != nil {
	// 	return nil, err
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successfully",
		"token":   token,
		// 返回用户角色，后续可以根据角色做权限验证
		"role": user.Role,
	})
}

// 注册用户接口
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Salt     string `json:"salt"`
}

// 普通用户注册接口
func Register(c *gin.Context) {
	// NOTES: 判断是否已经有管理员, 如果没有管理员，第一位注册的用户自动成为管理员!
	role := string(store.UserRoleMember)
	var firstAdmin store.User
	if err := store.GetDB().Where("role = ?", store.UserRoleAdmin).First(&firstAdmin).Error; err != nil && err.Error() == "record not found" {
		role = string(store.UserRoleAdmin)
	}
	register(c, role)
}

// 管理员注册接口
func RegisterAdmin(c *gin.Context) {
	register(c, string(store.UserRoleAdmin))
}

func register(c *gin.Context, role string) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否不合法或已存在
	if err := utils.CheckUsername(request.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
		c.Abort()
		return
	}
	var existingUser store.User
	if err := store.GetDB().Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		c.Abort()
		return
	}

	// 生成盐值与密码
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)
	saltString := fmt.Sprintf("%x", salt)
	hashedPasswd := utils.Sha256(request.Password, saltString)

	// 插入用户
	newUser := store.User{
		Username:     request.Username,
		PasswordHash: hashedPasswd,
		Email:        request.Email,
		Role:         role,
		Salt:         saltString,
	}
	if err := store.GetDB().Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
