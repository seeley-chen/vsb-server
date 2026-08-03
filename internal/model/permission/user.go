package model

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string
	Password string
	Email    string
	Account  string
	Phone    string
}

// User 用户模型，对应 MongoDB users 集合
type User struct {
	UserId       string    `bson:"userId" json:"userId"`
	Username     string    `bson:"username" json:"username"`
	PasswordHash string    `bson:"passwordHash" json:"-"`
	Email        string    `bson:"email" json:"email"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt" json:"updatedAt"`
	Account      string    `bson:"account" json:"account"`
	Phone        string    `bson:"phone" json:"phone"`
}

// generateUserID 生成纯数字字符串用户 ID（毫秒时间戳 + 4 位随机数）
func generateUserID() string {
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixMilli())
	}
	return fmt.Sprintf("%d%04d", time.Now().UnixMilli(), n.Int64())
}

// NewUser 创建新用户，自动生成数字 ID、加密密码、设置时间戳
func NewUser(req RegisterRequest) (*User, error) {
	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err // 如果生成密码哈希失败，则返回错误
	}

	// 设置时间戳
	now := time.Now()

	// 返回用户模型
	return &User{
		UserId:       generateUserID(),
		Username:     req.Username,
		Account:      req.Account,
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),
		Email:        req.Email,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// CheckPassword 验证密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
