package api

import (
   
    "net/http"
    "rental-property-management-system/internal/database"
    "rental-property-management-system/internal/models"
    
    "github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
    var loginData struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := database.DB.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
        return
    }

    // 验证密码
    if !verifyPassword(user.Password, loginData.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
        return
    }

    // 生成JWT token
    token := "your-generated-jwt-token" 

    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "token":   token,
    })
}
