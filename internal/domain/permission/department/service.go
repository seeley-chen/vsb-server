package department

import (
	"context"
	"errors"

	"github.com/Vanselyn/vsb-server/internal/tools"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrDepartmentExists    = errors.New("department already exists")
	ErrDepartmentNotFound  = errors.New("department not found")
	ErrDepartmentNameEmpty = errors.New("department name is empty")
)

type DepartmentService struct {
	repo *DepartmentRepo
}

func NewDepartmentService(repo *DepartmentRepo) *DepartmentService {
	return &DepartmentService{repo: repo}
}

/** 创建部门 */
func (s *DepartmentService) CreateDepartment(ctx context.Context, data *DepartmentRequest) (*DepartmentResponse, error) {
	data.Name = tools.TrimSpace(data.Name)
	data.Description = tools.TrimSpace(data.Description)
	if data.Name == "" {
		return nil, ErrDepartmentNameEmpty
	}

	existing, err := s.repo.FindByName(ctx, data.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDepartmentExists
	}

	department, err := s.repo.CreateDepartment(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrDepartmentExists
		}
		return nil, err
	}
	return department, nil
}

/** 根据部门ID获取部门 */
func (s *DepartmentService) GetDepartmentById(ctx context.Context, departmentId string) (*DepartmentResponse, error) {
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

/** 获取部门列表 */
func (s *DepartmentService) GetDepartmentList(ctx context.Context, pageIndex, pageSize int, name string) ([]*DepartmentResponse, int64, error) {
	return s.repo.GetDepartmentList(ctx, pageIndex, pageSize, name)
}

/** 更新部门 */
func (s *DepartmentService) UpdateDepartment(ctx context.Context, departmentId string, data *DepartmentUpdateRequest) (*DepartmentResponse, error) {
	if departmentId == "" {
		return nil, ErrDepartmentNotFound
	}

	if data.Name != "" {
		data.Name = tools.TrimSpace(data.Name)
		if data.Name == "" {
			return nil, ErrDepartmentNameEmpty
		}
	}
	data.Description = tools.TrimSpace(data.Description)
	data.DepartmentId = departmentId

	if data.Name != "" {
		existing, err := s.repo.FindByName(ctx, data.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.DepartmentId != departmentId {
			return nil, ErrDepartmentExists
		}
	}

	department, err := s.repo.UpdateDepartment(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrDepartmentExists
		}
		return nil, err
	}
	if department == nil {
		return nil, ErrDepartmentNotFound
	}
	return department, nil
}

/** 删除部门 */
func (s *DepartmentService) DeleteDepartment(ctx context.Context, departmentId string) error {
	if departmentId == "" {
		return ErrDepartmentNotFound
	}

	err := s.repo.DeleteDepartment(ctx, departmentId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrDepartmentNotFound
		}
		return err
	}
	return nil
}
