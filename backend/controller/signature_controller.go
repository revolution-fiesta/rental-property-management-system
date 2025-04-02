package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"rental-property-management-system/backend/store"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 上传手写签名
type UploadSignatureRequest struct {
	ImageData string `json:"image_data"`
}

func UploadSignature(c *gin.Context) {
	var request UploadSignatureRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 解码 Base64 图片数据
	imageData := request.ImageData
	// 去除可能的前缀部分，比如 data:image/png;base64, 或 data:image/jpeg;base64,
	if strings.HasPrefix(imageData, "data:image/") {
		// 找到 Base64 字符串的起始部分并去除
		imageData = strings.Split(imageData, ",")[1]
	}
	decodedImage, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode Base64 image"})
		return
	}

	// 确保文件夹存在,
	// 并保存文件到本地 statics 文件夹
	fileName := fmt.Sprintf("signature_%d.png", time.Now().Unix())
	if err := os.MkdirAll("statics", os.ModePerm); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}
	filePath := fmt.Sprintf("statics/%s", fileName)
	if err := os.WriteFile(filePath, decodedImage, 0644); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// 保存文件路径到数据库
	signature := store.Signature{
		ImageData: filePath, // 保存图片的文件路径
	}
	if err := store.GetDB().Create(&signature).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save signature to database"})
		return
	}

	// 返回签名 ID 和文件路径
	c.JSON(http.StatusOK, gin.H{
		"signature_id": signature.ID,
	})
}
