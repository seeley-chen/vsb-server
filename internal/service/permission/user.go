package permission

import (
	"context"
	"errors"
	"time"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	repo "github.com/seeley-chen/vsb-server/internal/repository/permission"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrUserNameEmpty       = errors.New("user name is empty")
	ErrUserAccountEmpty    = errors.New("user account is empty")
	ErrUserRoleEmpty       = errors.New("user role is empty")
	ErrUserDepartmentEmpty = errors.New("user department is empty")
)

type UserService struct {
	userRepo  *repo.UserRepo
	jwtSecret string
	jwtExpire time.Duration
}

func NewUserService(userRepo *repo.UserRepo, jwtSecret string, jwtExpire time.Duration) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
	}
}

// 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *model.UserRequest) (*model.UserResponse, error) {
	req.Username = normalizeName(req.Username)

	if req.Username == "" {
		return nil, ErrUserNameEmpty
	}

	if req.Account == "" {
		return nil, ErrUserAccountEmpty
	}

	if req.Password == "" {
		return nil, ErrInvalidPassword
	}

	if req.RoleId == "" {
		return nil, ErrUserRoleEmpty
	}

	if req.DepartmentId == "" {
		return nil, ErrUserDepartmentEmpty
	}

	existing, err := s.userRepo.GetUserByAccount(ctx, req.Account)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	user, err := s.userRepo.CreateUser(ctx, req)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}
	return user, nil
}

// 根据用户ID获取用户
func (s *UserService) GetUserById(ctx context.Context, userId string) (*model.UserResponse, error) {
	if userId == "" {
		return nil, ErrUserNotFound
	}

	user, err := s.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// 获取用户列表（username 模糊搜索，account/email 精确匹配，均可选）
func (s *UserService) GetUserList(ctx context.Context, pageIndex, pageSize int, username, account, email string) ([]*model.UserResponse, int64, error) {
	return s.userRepo.GetUserList(ctx, pageIndex, pageSize, username, account, email)
}

// 更新用户
func (s *UserService) UpdateUser(ctx context.Context, userId string, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	if userId == "" {
		return nil, ErrUserNotFound
	}
	if req.Username != "" {
		req.Username = normalizeName(req.Username)
		if req.Username == "" {
			return nil, ErrUserNameEmpty
		}
	}
	if req.Account != "" {
		req.Account = normalizeName(req.Account)
		if req.Account == "" {
			return nil, ErrUserAccountEmpty
		}
	}

	req.UserId = userId

	if req.UserId != "" {
		existing, err := s.userRepo.GetUserById(ctx, req.UserId)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.UserId != userId {
			return nil, ErrUserAlreadyExists
		}
	}

	user, err := s.userRepo.UpdateUser(ctx, req)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}
	if user == nil {
		return nil, ErrDepartmentNotFound
	}
	return user, nil
}

// 删除用户
func (s *UserService) DeleteUser(ctx context.Context, userId string) error {
	if userId == "" {
		return ErrUserNotFound
	}

	err := s.userRepo.DeleteUser(ctx, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
