package service

import (
	"errors"
	"os"
	"time"

	"fukubukuro/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "fukubukuro-secret-key"
}

// Register 注册新用户
func Register(db *gorm.DB, username, password string) (*model.User, error) {
	// 检查用户名是否已存在
	var exist model.User
	if err := db.Where("username = ?", username).First(&exist).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// Bcrypt 加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := model.User{
		Username:     username,
		PasswordHash: string(hashed),
		TokenVersion: 1,
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, errors.New("创建用户失败")
	}
	return &user, nil
}

// Login 登录，返回 JWT Token
func Login(db *gorm.DB, username, password string) (string, error) {
	var user model.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 签发 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":       user.ID,
		"username":      user.Username,
		"token_version": user.TokenVersion,
		"exp":           time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.New("Token生成失败")
	}
	return tokenString, nil
}
