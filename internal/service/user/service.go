package user

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/seeley-chen/vsb-server/internal/model"
	repo "github.com/seeley-chen/vsb-server/internal/repository/user"
	"github.com/seeley-chen/vsb-server/pkg/jwt"
)

// 自定义业务错误
var (
	ErrUsernameExists  = errors.New("username already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

// Service 用户业务逻辑层
type Service struct {
	repo      *repo.Repository
	jwtSecret string
	jwtExpire time.Duration
}

// NewService 创建 Service 实例
func NewService(r *repo.Repository, jwtSecret string, jwtExpire time.Duration) *Service {
	return &Service{
		repo:      r,
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
	}
}

// Register 用户注册
func (s *Service) Create(ctx context.Context, req model.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	existingUser, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err // 系统错误
	}
	if existingUser != nil {
		return nil, ErrUsernameExists
	}

	// 创建新用户
	user, err := model.NewUser(req)
	if err != nil {
		return nil, err
	}

	// 插入数据库
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录，返回 JWT token 和用户信息
func (s *Service) Login(ctx context.Context, username, password string) (string, *model.User, error) {
	// 查找用户
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", nil, ErrUserNotFound
		}
		return "", nil, err
	}

	// 验证密码
	if !user.CheckPassword(password) {
		return "", nil, ErrInvalidPassword
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(user.UserId, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// ListUsers 分页获取用户列表
func (s *Service) ListUsers(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // 默认每页20条
	}

	users, total, err := s.repo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
