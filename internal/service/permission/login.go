package permission

import (
	"context"
	"errors"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	"github.com/seeley-chen/vsb-server/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	/** 用户未找到 */
	ErrInvalidAccount = errors.New("invalid account")
)

func (s *UserService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.UserLoginResponse, error) {
	if req.Account == "" {
		return nil, ErrInvalidAccount
	}

	if req.Password == "" {
		return nil, ErrInvalidPassword
	}

	user, err := s.userRepo.GetUserByAccount(ctx, req.Account)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidAccount
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(user.UserId, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, err
	}

	return &model.UserLoginResponse{
		Token: token,
		User:  user,
	}, nil
}
