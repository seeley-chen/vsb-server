package user

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	model "github.com/seeley-chen/vsb-server/internal/model/permission"
	repo "github.com/seeley-chen/vsb-server/internal/repository/permission/user"
	"github.com/seeley-chen/vsb-server/pkg/jwt"
)

// 自定义业务错误
var (
	/** 用户名已存在 */
	ErrAccountExists = errors.New("account already exists")
	/** 用户未找到 */
	ErrUserNotFound = errors.New("user not found")
	/** 密码不正确 */
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
	existingUser, err := s.repo.FindByAccount(ctx, req.Account)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err // 系统错误
	}
	if existingUser != nil {
		return nil, ErrAccountExists
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
func (s *Service) Login(ctx context.Context, account, password string) (string, *model.User, error) {
	// 查找用户
	user, err := s.repo.FindByAccount(ctx, account)
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

// DeleteUserById 根据用户ID删除用户
func (s *Service) DeleteUserById(ctx context.Context, userId string) error {
	if userId == "" {
		return ErrUserNotFound
	}
	err := s.repo.DeleteUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

// UpdateUserById 根据用户ID更新用户
func (s *Service) UpdateUserById(ctx context.Context, userId string, update model.User) error {
	if userId == "" {
		return ErrUserNotFound
	}
	err := s.repo.UpdateUserById(ctx, userId, &update)
	if err != nil {
		return err
	}
	return nil
}
