package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"rental-property-management-system/backend/config"
	"rental-property-management-system/backend/store"
	"rental-property-management-system/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 登录接口
type AuthMethod string

const (
	AuthMethodPlain  AuthMethod = "plain"
	AuthMethodWechat AuthMethod = "wechat"
)

type LoginRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	AuthMethod string `json:"auth_method"`
	OAuthCode  string `json:"code"`
}

func Login(c *gin.Context) {
	var reqeust LoginRequest
	if err := c.ShouldBindJSON(&reqeust); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 检查用户名是否存在
	var user store.User
	// 微信登录
	if reqeust.AuthMethod == string(AuthMethodWechat) {
		openID, err := getWechatOpenID(reqeust.OAuthCode)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 如果该 OpenID 没有关联的账号，则创建新账号
		if tx := store.GetDB().Where("open_id = ?", openID).Find(&user); tx.RowsAffected == 0 {
			if err := store.CreateUser("", "", "", string(store.UserRoleMember), openID); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				return
			}
		}
		// 重新获取用户
		if tx := store.GetDB().Where("open_id = ?", openID).Find(&user); tx.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
			return
		}
		// 用户名密码登录
	} else {
		if reqeust.Username == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
			return
		}
		if err := store.GetDB().Where("username = ?", reqeust.Username).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("User %q does not exist", reqeust.Username)})
			return
		}
		// 验证密码是否正确并生成
		if user.PasswordHash != utils.Sha256(reqeust.Password, user.Salt) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Wrong username or password"})
			return
		}
	}

	// 生成 access token
	token, err := utils.GenerateAccessToken(int(user.ID), config.PrivateKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
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
		"role":    user.Role,
		"open_id": user.OpenID,
	})
}

type WechatOAuthResponse struct {
	SessionKey string `json:"session_key"`
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid,omitempty"`
}

func getWechatOpenID(code string) (string, error) {
	params := url.Values{
		"appid":      {config.AppConfig.Wechat.AppID},
		"secret":     {config.AppConfig.Wechat.AppSecret},
		"js_code":    {code},
		"grant_type": {"authorization_code"},
	}
	url := "https://api.weixin.qq.com/sns/jscode2session"
	resp, err := http.Get(fmt.Sprintf("%s?%s", url, params.Encode()))
	if err != nil {
		return "", errors.Wrapf(err, "failed to send oauth request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var oauthResp WechatOAuthResponse
	if err := json.Unmarshal(body, &oauthResp); err != nil {
		return "", err
	}

	if oauthResp.OpenID == "" {
		return "", errors.New("failed to get OpenID")
	}

	return oauthResp.OpenID, nil
}

// 注册用户接口
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// 普通用户注册接口
func Register(c *gin.Context) {
	// NOTES: 判断是否已经有管理员, 如果没有管理员，第一位注册的用户自动成为管理员!
	role := string(store.UserRoleMember)
	var firstAdmin store.User

	result := store.GetDB().Where("role = ?", store.UserRoleAdmin).Find(&firstAdmin)
	if result.RowsAffected == 0 {
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否不合法或已存在
	if err := utils.CheckUsername(request.Username); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
		return
	}
	var existingUser store.User
	result := store.GetDB().Where("username = ?", request.Username).Find(&existingUser)
	if result.RowsAffected > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	if err := store.CreateUser(request.Username, request.Password, request.Email, role, ""); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
