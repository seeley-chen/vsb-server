package permission

import (
	"context"
	"errors"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	repo "github.com/seeley-chen/vsb-server/internal/repository/permission"
)

var (
	/** 未找到角色 */
	ErrRoleNotFound = errors.New("role not found")
	/** 角色已存在 */
	ErrRoleExist = errors.New("role already exist")
	/** 角色名不能为空 */
	ErrRoleNameEmpty = errors.New("role name is empty")
)

type RoleService struct {
	repo *repo.RoleRepo
}

func NewRoleService(repo *repo.RoleRepo) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

// 创建角色
func (s *RoleService) CreateRole(ctx context.Context, data *model.RoleRequest) (*model.RoleResponse, error) {
	// 检查角色名是否已存在
	existingRole, err := s.repo.FindByRoleName(ctx, data.Name)
	if err != nil {
		return nil, err
	}

	if existingRole != nil {
		return nil, ErrRoleExist
	}

	// 创建新用户
	if data.Name == "" {
		return nil, ErrRoleNameEmpty
	}

	return s.repo.CreateRole(ctx, data)
}

// 通过角色ID获取角色
func (s *RoleService) GetRoleById(ctx context.Context, roleId string) (*model.RoleResponse, error) {
	if roleId == "" {
		return nil, ErrRoleNotFound
	}

	role, err := s.repo.GetRoleById(ctx, roleId)

	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, ErrRoleNotFound
	}

	return role, nil
}

// 分页获取角色列表
func (s *RoleService) GetRoleList(ctx context.Context, pageIndex, pageSize int) ([]*model.RoleResponse, int64, error) {
	return s.repo.GetRoleList(ctx, pageIndex, pageSize)
}

// 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, roleId string, data *model.RoleUpdateRequest) (*model.RoleResponse, error) {
	if _, err := s.GetRoleById(ctx, roleId); err != nil {
		return nil, err
	}

	data.RoleId = roleId
	return s.repo.UpdateRole(ctx, data)
}

// 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, roleId string) error {
	if _, err := s.GetRoleById(ctx, roleId); err != nil {
		return err
	}

	return s.repo.DeleteRole(ctx, roleId)
}
