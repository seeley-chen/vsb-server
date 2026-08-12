package permission

import (
	"strings"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
)

const maxPermissionDepth = 5

func normalizeName(name string) string {
	return strings.TrimSpace(name)
}

func normalizeDescription(desc string) string {
	return strings.TrimSpace(desc)
}

// validatePermissions 校验权限树：path 非空、type 仅 read/write、限制深度。
func validatePermissions(items []model.PermissionItem) error {
	return validatePermissionsDepth(items, 0)
}

func validatePermissionsDepth(items []model.PermissionItem, depth int) error {
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
			item.Children = []model.PermissionItem{}
		}
		if err := validatePermissionsDepth(item.Children, depth+1); err != nil {
			return err
		}
	}
	return nil
}
