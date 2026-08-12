package permission

import (
	"context"
	"errors"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	repo "github.com/seeley-chen/vsb-server/internal/repository/permission"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleExists        = errors.New("role already exists")
	ErrRoleNameEmpty     = errors.New("role name is empty")
	ErrInvalidPermission = errors.New("invalid permission")
)

type RoleService struct {
	repo *repo.RoleRepo
}

func NewRoleService(repo *repo.RoleRepo) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) CreateRole(ctx context.Context, data *model.RoleRequest) (*model.RoleResponse, error) {
	data.Name = normalizeName(data.Name)
	data.Description = normalizeDescription(data.Description)
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

func (s *RoleService) GetRoleList(ctx context.Context, pageIndex, pageSize int) ([]*model.RoleResponse, int64, error) {
	return s.repo.GetRoleList(ctx, pageIndex, pageSize)
}

func (s *RoleService) UpdateRole(ctx context.Context, roleId string, data *model.RoleUpdateRequest) (*model.RoleResponse, error) {
	if roleId == "" {
		return nil, ErrRoleNotFound
	}

	if data.Name != "" {
		data.Name = normalizeName(data.Name)
		if data.Name == "" {
			return nil, ErrRoleNameEmpty
		}
	}
	data.Description = normalizeDescription(data.Description)
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
