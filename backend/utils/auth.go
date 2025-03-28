package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	mathRand "math/rand"
	"rental-property-management-system/backend/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	HttpHeaderAuthorization string = "authorization"
	HttpHeaderBearer        string = "bearer"
)

type AudienceType string

const (
	audienceTypeAccessTokenUser AudienceType = "access_token_user"
)

func GenerateAccessToken(userId int, privateKey *rsa.PrivateKey) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userId),
		Issuer:    config.ProjectName,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(config.AccessTokenExpiration)),
		Audience:  jwt.ClaimStrings{string(audienceTypeAccessTokenUser)},
		// actually this is only used for key rotation.
		ID: uuid.NewString(),
	})
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.Wrapf(err, "failed to sign token")
	}
	return signedToken, nil
}

// validate access token and return user id.
func ValidateAccessToken(token string, publicKey *rsa.PublicKey) (int, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if alg := t.Method.Alg(); alg != jwt.SigningMethodRS256.Name {
			return nil, errors.Errorf("sighing method not supported: %q", alg)
		}
		return publicKey, nil
	})
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse claims")
	}

	if len(claims.Audience) != 1 || claims.Audience[0] != string(audienceTypeAccessTokenUser) {
		return 0, errors.New(fmt.Sprintf("audience type not supported: %q", claims.Audience[0]))
	}

	if claims.Issuer != config.ProjectName {
		return 0, errors.New(fmt.Sprintf("wrong issuer: %q", claims.Issuer))
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

// hash salt and password with SHA-256.
func Sha256(text string, salt string) string {
	hasher := sha256.New()
	_, _ = hasher.Write([]byte(salt))
	_, _ = hasher.Write([]byte(text))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func GeneratePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func CheckUsername(username string) error {
	//用户名长度限制
	if len(username) == 0 {
		return errors.New("The username is empty")
	} else if len(username) < 6 {
		return errors.New("The username length is too short")
	} else if len(username) > 15 {
		return errors.New("The username length is too long")
	}
	//用户名首个字符限制
	if username[0] >= '0' && username[0] <= '9' {
		return errors.New("Usernames start with numbers")
	}
	//用户名不允许含有特殊字符
	for _, val := range username {
		if !((val >= 'a' && val <= 'z') || (val >= 'A' && val <= 'Z') || (val >= '0' && val <= '9') || val == '_' || val == '@') {
			return errors.New("Usernames contain special characters")
		}
	}
	return nil
}

// 生成指定位数的密码
func GenerateRandomPassword(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
