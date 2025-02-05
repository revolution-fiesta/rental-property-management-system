package api

import (
    "fmt"
    "net/http"
    "rental-property-management-system/internal/database"
    "rental-property-management-system/internal/models"
    "golang.org/x/crypto/argon2"
    "github.com/gin-gonic/gin"
)

func hashPassword(password string) (string, error) {
    salt := []byte("some_random_salt") // 你可以使用更复杂的盐值或从用户获取
    hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32) // 生成哈希值
    return fmt.Sprintf("%x", hash), nil
}

func verifyPassword(storedHash, password string) bool {
    salt := []byte("some_random_salt") // 必须使用相同的盐值
    hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
    return fmt.Sprintf("%x", hash) == storedHash
}

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
    hashedPassword, err := hashPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user.Password = hashedPassword

    // 插入用户
    if err := database.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

