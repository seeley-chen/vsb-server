package permission

import (
	"context"
	"errors"

	models "github.com/seeley-chen/vsb-server/internal/models/permission"
	repo "github.com/seeley-chen/vsb-server/internal/repository/permission"
)

var (
	ErrDepartmentExists     = errors.New("department already exists")
	ErrDepartmentNotFound   = errors.New("department not found")
	ErrDepartmentNameFailed = errors.New("department name is empty")
)

type DepartmentService struct {
	repo *repo.DepartmentRepo
}

// NewDepartmentService 创建部门服务
func NewDepartmentService(repo *repo.DepartmentRepo) *DepartmentService {
	return &DepartmentService{
		repo: repo,
	}
}

// CreateDepartment 创建部门
func (s *DepartmentService) CreateDepartment(ctx context.Context, data *models.DepartmentRequest) (*models.Department, error) {
	if data.Name == "" {
		return nil, ErrDepartmentNameFailed
	}
	return s.repo.CreateDepartment(ctx, data)
}

// GetDepartmentById 通过ID获取部门
func (s *DepartmentService) GetDepartmentById(ctx context.Context, departmentId string) (*models.Department, error) {
	if departmentId == "" {
		return nil, ErrDepartmentNotFound
	}

	department, err := s.repo.GetDepartmentById(ctx, departmentId)
	if err != nil {
		return nil, err
	}
	if department == nil {
		return nil, ErrDepartmentNotFound
	}

	return department, nil
}

// GetDepartmentList 分页获取部门列表
func (s *DepartmentService) GetDepartmentList(ctx context.Context, pageIndex, pageSize int) ([]*models.Department, int64, error) {
	return s.repo.GetDepartmentList(ctx, pageIndex, pageSize)
}

// UpdateDepartment 更新部门
func (s *DepartmentService) UpdateDepartment(ctx context.Context, departmentId string, data *models.DepartmentUpdateRequest) (*models.Department, error) {
	if _, err := s.GetDepartmentById(ctx, departmentId); err != nil {
		return nil, err
	}

	data.DepartmentId = departmentId
	return s.repo.UpdateDepartment(ctx, data)
}

// DeleteDepartment 删除部门
func (s *DepartmentService) DeleteDepartment(ctx context.Context, departmentId string) error {
	if _, err := s.GetDepartmentById(ctx, departmentId); err != nil {
		return err
	}
	return s.repo.DeleteDepartment(ctx, departmentId)
}
