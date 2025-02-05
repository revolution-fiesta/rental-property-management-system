package api
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	_"github.com/lib/pq"
	"rental-property-management-system/internal/models"
	"rental-property-management-system/internal/database"
)
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.PasswordHash = string(hashedPassword)

	// 插入用户
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
