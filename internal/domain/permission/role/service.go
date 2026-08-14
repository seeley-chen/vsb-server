package role

import (
	"context"
	"errors"
	"strings"

	"github.com/Vanselyn/vsb-server/internal/tools"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleExists        = errors.New("role already exists")
	ErrRoleNameEmpty     = errors.New("role name is empty")
	ErrInvalidPermission = errors.New("invalid permission")
)

const maxPermissionDepth = 5

type RoleService struct {
	repo *RoleRepo
}

func NewRoleService(repo *RoleRepo) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) CreateRole(ctx context.Context, data *RoleRequest) (*RoleResponse, error) {
	data.Name = tools.TrimSpace(data.Name)
	data.Description = tools.TrimSpace(data.Description)
	if data.Name == "" {
		return nil, ErrRoleNameEmpty
	}
	if err := validatePermissions(data.Permissions); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByName(ctx, data.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrRoleExists
	}

	role, err := s.repo.CreateRole(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrRoleExists
		}
		return nil, err
	}
	return role, nil
}

func (s *RoleService) GetRoleById(ctx context.Context, roleId string) (*RoleResponse, error) {
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

// GetRolesByIds 根据角色ID列表批量获取角色
func (s *RoleService) GetRolesByIds(ctx context.Context, roleIds []string) ([]*RoleResponse, error) {
	return s.repo.GetRolesByIds(ctx, roleIds)
}

// MergePermissions 合并多个角色的权限（简单拼接，保留所有权限项）
func (s *RoleService) MergePermissions(roles []*RoleResponse) []PermissionItem {
	merged := make([]PermissionItem, 0)
	for _, role := range roles {
		if role != nil && role.Permissions != nil {
			merged = append(merged, role.Permissions...)
		}
	}
	return merged
}

func (s *RoleService) GetRoleList(ctx context.Context, pageIndex, pageSize int) ([]*RoleResponse, int64, error) {
	return s.repo.GetRoleList(ctx, pageIndex, pageSize)
}

func (s *RoleService) UpdateRole(ctx context.Context, roleId string, data *RoleUpdateRequest) (*RoleResponse, error) {
	if roleId == "" {
		return nil, ErrRoleNotFound
	}

	if data.Name != "" {
		data.Name = tools.TrimSpace(data.Name)
		if data.Name == "" {
			return nil, ErrRoleNameEmpty
		}
	}
	data.Description = tools.TrimSpace(data.Description)
	data.RoleId = roleId

	if data.Permissions != nil {
		if err := validatePermissions(data.Permissions); err != nil {
			return nil, err
		}
	}

	if data.Name != "" {
		existing, err := s.repo.FindByName(ctx, data.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.RoleId != roleId {
			return nil, ErrRoleExists
		}
	}

	role, err := s.repo.UpdateRole(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrRoleExists
		}
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, roleId string) error {
	if roleId == "" {
		return ErrRoleNotFound
	}

	err := s.repo.DeleteRole(ctx, roleId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrRoleNotFound
		}
		return err
	}
	return nil
}

// validatePermissions 校验权限树：path 非空、type 仅 read/write、限制深度。
func validatePermissions(items []PermissionItem) error {
	return validatePermissionsDepth(items, 0)
}

// 校验权限树深度
func validatePermissionsDepth(items []PermissionItem, depth int) error {
	if depth > maxPermissionDepth {
		return ErrInvalidPermission
	}
	for i := range items {
		item := &items[i]
		item.Path = strings.TrimSpace(item.Path)
		item.Type = strings.TrimSpace(item.Type)
		if item.Path == "" {
			return ErrInvalidPermission
		}
		if item.Type != "read" && item.Type != "write" {
			return ErrInvalidPermission
		}
		if item.Children == nil {
			item.Children = []PermissionItem{}
		}
		if err := validatePermissionsDepth(item.Children, depth+1); err != nil {
			return err
		}
	}
	return nil
}
