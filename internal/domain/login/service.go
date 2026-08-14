package login

import (
	"context"
	"errors"
	"time"

	"github.com/Vanselyn/vsb-server/internal/domain/permission/user"
	"github.com/Vanselyn/vsb-server/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidAccount  = errors.New("invalid account")
	ErrInvalidPassword = errors.New("invalid password")
)

type LoginService struct {
	userRepo  *user.UserRepo
	jwtSecret string
	jwtExpire time.Duration
}

func NewLoginService(userRepo *user.UserRepo, jwtSecret string, jwtExpire time.Duration) *LoginService {
	return &LoginService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
	}
}

// Login 校验账号密码并签发 JWT token
func (s *LoginService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	if req.Account == "" {
		return nil, ErrInvalidAccount
	}

	if req.Password == "" {
		return nil, ErrInvalidPassword
	}

	u, err := s.userRepo.GetUserByAccount(ctx, req.Account)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrInvalidAccount
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(u.UserId, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  u,
	}, nil
}
