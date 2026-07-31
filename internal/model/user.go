package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User 用户模型，对应 MongoDB users 集合
type User struct {
	UserId       string    `bson:"user_id" json:"user_id"`
	Username     string    `bson:"username" json:"username"`
	PasswordHash string    `bson:"password_hash" json:"-"`
	Email        string    `bson:"email" json:"email"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}

// NewUser 创建新用户，自动生成 UUID、加密密码、设置时间戳
func NewUser(username, password, email string) (*User, error) {
	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err // 如果生成密码哈希失败，则返回错误
	}

	// 设置时间戳
	now := time.Now()

	// 返回用户模型
	return &User{
		UserId:       uuid.New().String(),
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        email,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// CheckPassword 验证密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
