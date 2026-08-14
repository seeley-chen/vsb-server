package user

import (
	"context"
	"errors"

	"github.com/Vanselyn/vsb-server/internal/tools"
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
	repo *UserRepo
}

func NewUserService(repo *UserRepo) *UserService {
	return &UserService{repo: repo}
}

// 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	req.Username = tools.TrimSpace(req.Username)

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

	existing, err := s.repo.GetUserByAccount(ctx, req.Account)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	user, err := s.repo.CreateUser(ctx, req)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}
	return user, nil
}

// 根据用户ID获取用户
func (s *UserService) GetUserById(ctx context.Context, userId string) (*UserResponse, error) {
	if userId == "" {
		return nil, ErrUserNotFound
	}

	user, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// 获取用户列表（username/account 模糊搜索，email 精确匹配，均可选）
func (s *UserService) GetUserList(ctx context.Context, pageIndex, pageSize int, username, account, email string) ([]*UserResponse, int64, error) {
	return s.repo.GetUserList(ctx, pageIndex, pageSize, username, account, email)
}

// 更新用户
func (s *UserService) UpdateUser(ctx context.Context, userId string, req *UserUpdateRequest) (*UserResponse, error) {
	if userId == "" {
		return nil, ErrUserNotFound
	}
	if req.Username != "" {
		req.Username = tools.TrimSpace(req.Username)
		if req.Username == "" {
			return nil, ErrUserNameEmpty
		}
	}
	if req.Account != "" {
		req.Account = tools.TrimSpace(req.Account)
		if req.Account == "" {
			return nil, ErrUserAccountEmpty
		}
	}

	req.UserId = userId

	if req.Account != "" {
		existing, err := s.repo.GetUserByAccount(ctx, req.Account)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.UserId != userId {
			return nil, ErrUserAlreadyExists
		}
	}

	user, err := s.repo.UpdateUser(ctx, req)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// 删除用户
func (s *UserService) DeleteUser(ctx context.Context, userId string) error {
	if userId == "" {
		return ErrUserNotFound
	}

	err := s.repo.DeleteUser(ctx, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
