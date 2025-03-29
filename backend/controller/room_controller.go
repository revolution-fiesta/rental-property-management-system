package controller

import (
	"math"
	"net/http"
	"rental-property-management-system/backend/store"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 管理员更新房间的信息
type UpdateRoomInfoRequest struct {
	RoomID    uint     `json:"room_id" binding:"required"`
	Type      *string  `json:"type"`
	Price     *float64 `json:"price"`
	Available *bool    `json:"available"`
	Tags      *string  `json:"tags"`
	Area      *float64 `json:"area"`
}

func UpdateRoomInfo(c *gin.Context) {
	var request UpdateRoomInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 查找房间,  如果找不到房间，则创建新的房间
	var room store.Room
	result := store.GetDB().Find(&room, request.RoomID)
	if result.RowsAffected == 0 {
		room = store.Room{
			Type:      store.RoomType(*request.Type),
			Price:     *request.Price,
			Available: *request.Available,
			Tags:      *request.Tags,
			Area:      *request.Area,
		}
		if err := store.GetDB().Create(&room).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Room created successfully", "room": room})
		return
	}

	// 数据清洗
	if request.Type != nil {
		room.Type = store.RoomType(*request.Type)
	}
	if request.Price != nil {
		room.Price = *request.Price
	}
	if request.Available != nil {
		room.Available = *request.Available
	}
	if request.Tags != nil {
		room.Tags = *request.Tags
	}
	if request.Area != nil {
		room.Area = *request.Area
	}

	// 更新数据库
	if err := store.GetDB().Save(&room).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room updated successfully", "room": room})
}

// 查询所有房间接口
func ListRooms(c *gin.Context) {
	var rooms []store.Room
	// 查询所有房间数据
	if err := store.GetDB().Find(&rooms).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}
	// 返回房间信息
	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}

// 查询房间信息，支持多个过滤条件
type GetFilteredRoomsRequest struct {
	RoomType *store.RoomType `json:"room_type"`
	MinPrice *uint32         `json:"min_price"`
	MaxPrice *uint32         `json:"max_price"`
	MinArea  *uint32         `json:"min_area"`
	MaxArea  *uint32         `json:"max_area"`
	Keyword  *string         `json:"keyword"`
}

func ListFilteredRooms(c *gin.Context) {
	var request GetFilteredRoomsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to parse request"})
		return
	}

	// 数据清洗
	var rooms []store.Room
	db := store.GetDB()

	if request.MinPrice == nil {
		*request.MinPrice = 0
	}
	if request.MaxPrice == nil {
		*request.MaxPrice = math.MaxUint32
	}
	if request.MinArea == nil {
		*request.MinArea = 0
	}
	if request.MaxArea == nil {
		*request.MaxArea = math.MaxUint32
	}

	// 设置数据库查询条件
	db = db.Where("price >= ? AND price <= ?", request.MinPrice, request.MaxPrice)
	db = db.Where("area >= ? AND area <= ?", request.MinArea, request.MaxArea)
	if request.RoomType != nil {
		db = db.Where("type = ?", request.RoomType)
	}
	if request.Keyword != nil {
		db = db.Where("name LIKE ?", "%"+*request.Keyword+"%")
	}

	// 查询符合条件的房间数据
	if err := db.Find(&rooms).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}
