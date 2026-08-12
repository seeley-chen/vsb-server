package permission

import "time"

// PermissionItem 权限树节点
type PermissionItem struct {
	Path     string           `bson:"path" json:"path"`
	Type     string           `bson:"type" json:"type"` // read | write
	Children []PermissionItem `bson:"children" json:"children"`
}

// RoleRequest 创建角色请求
type RoleRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Permissions []PermissionItem `json:"permissions"`
}

// RoleResponse 角色响应 / 存储模型
type RoleResponse struct {
	RoleId      string           `bson:"role_id" json:"roleId"`
	Name        string           `bson:"name" json:"name"`
	Description string           `bson:"description" json:"description"`
	Permissions []PermissionItem `bson:"permissions" json:"permissions"`
	CreatedAt   time.Time        `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time        `bson:"updated_at" json:"updatedAt"`
}

// RoleUpdateRequest 更新角色请求（空字段表示不修改；permissions 为 nil 表示不修改）
type RoleUpdateRequest struct {
	RoleId      string           `json:"-"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Permissions []PermissionItem `json:"permissions"`
}
