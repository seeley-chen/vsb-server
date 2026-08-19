package user

import (
	"context"
	"errors"

	"github.com/Vanselyn/vsb-server/internal/tools"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserService struct {
	repo *UserRepo
}

func NewUserService(repo *UserRepo) *UserService {
	return &UserService{repo: repo}
}

// 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	if err := tools.ValidateStruct(req); err != nil {
		return nil, err
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
func (s *UserService) GetUserList(ctx context.Context, pageIndex, pageSize int, params UserQueryParams) ([]*UserResponse, int64, error) {
	return s.repo.GetUserList(ctx, pageIndex, pageSize, params)
}

// 更新用户
func (s *UserService) UpdateUser(ctx context.Context, userId string, req *UserUpdateRequest) (*UserResponse, error) {
	if userId == "" {
		return nil, ErrUserNotFound
	}
	if err := tools.ValidateStruct(req); err != nil {
		return nil, err
	}

	req.UserId = userId

	if req.UserId != "" {
		existing, err := s.repo.GetUserById(ctx, req.UserId)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.UserId != userId {
			return nil, ErrUserNotFound
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
