package role

import "time"

// 创建角色请求
type RoleRequest struct {
	Name        string     `json:"name" validate:"required"` // 角色名称
	Description string     `json:"description"`              // 角色描述
	Permissions []RoleTree `json:"permissions"`              // 权限树
}

// 角色响应 / 存储模型
type RoleResponse struct {
	RoleId      string     `bson:"role_id" json:"roleId" validate:"required"` // 角色ID
	Name        string     `json:"name" validate:"required"`                  // 角色名称
	Description string     `json:"description"`                               // 角色描述
	Permissions []RoleTree `json:"permissions"`                               // 权限树
	CreatedAt   time.Time  `json:"createdAt" validate:"required"`             // 创建时间
	UpdatedAt   time.Time  `json:"updatedAt" validate:"required"`             // 更新时间
}

// 更新角色请求（空字段表示不修改；permissions 为 nil 表示不修改）
type RoleUpdateRequest struct {
	RoleId      string     `json:"-"`           // 角色ID，以路径参数为准
	Name        string     `json:"name"`        // 角色名称，空表示不修改
	Description string     `json:"description"` // 角色描述
	Permissions []RoleTree `json:"permissions"` // 权限树，不传表示不修改
}

// 权限树节点
type RoleTree struct {
	Path     string         `json:"path" validate:"required"`                                // 权限路径
	Type     string         `bson:"type" json:"type" validate:"required" enums:"read,write"` // 权限类型：read 或 write
	Children []RoleTreeItem `bson:"children" json:"children"`                                // 子权限节点
}

// 权限树节点
type RoleTreeItem struct {
	Path string `json:"path" validate:"required" `                               // 权限路径
	Type string `bson:"type" json:"type" validate:"required" enums:"read,write"` // 权限类型：read 或 write
}
