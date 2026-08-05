package department

import (
	"context"
	"errors"
	"time"

	model "github.com/seeley-chen/vsb-server/internal/model/permission"
	repo "github.com/seeley-chen/vsb-server/internal/repository/permission/department"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrDepartmentExists   = errors.New("department already exists")
	ErrDepartmentNotFound = errors.New("department not found")
)

// Service 用户业务逻辑层
type Service struct {
	repo      *repo.Repository
	jwtSecret string
	jwtExpire time.Duration
}

func NewService(r *repo.Repository, jwtSecret string, jwtExpire time.Duration) *Service {
	return &Service{
		repo:      r,
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
	}
}

func (s *Service) Create(ctx context.Context, req model.DepartmentRequest) (*model.Department, error) {
	// 检查部门名称是否存在
	existingDepartment, err := s.repo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err // 系统错误
	}
	if existingDepartment != nil {
		return nil, ErrDepartmentExists
	}

	// 创建新部门
	department := model.NewDepartment(&req)

	// 保存部门到数据库
	err = s.repo.Create(ctx, department)
	if err != nil {
		return nil, err // 系统错误
	}

	return department, nil
}

func (s *Service) ListDepartments(ctx context.Context, page, pageSize int) ([]*model.Department, int64, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // 默认每页20条
	}

	departments, total, err := s.repo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return departments, total, nil
}
