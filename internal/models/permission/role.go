package permission

import "time"

type PermissionItem struct {
	Path    string   `json:"path"`
	Actions []string `json:"actions"`
}

type RoleRequest struct {
	Name        string           `json:"name" validate:"required"`
	Description string           `json:"description"`
	Permissions []PermissionItem `json:"permissions"`
}

type RoleResponse struct {
	RoleId      string           `bson:"role_id" json:"roleId"`
	Name        string           `bson:"name" json:"name"`
	Description string           `bson:"description" json:"description"`
	Permissions []PermissionItem `bson:"permissions" json:"permissions"`
	CreatedAt   time.Time        `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time        `bson:"updated_at" json:"updatedAt"`
}

type RoleUpdateRequest struct {
	RoleId      string           `bson:"role_id" json:"roleId" validate:"required"`
	Name        string           `json:"name" validate:"required"`
	Description string           `json:"description"`
	Permissions []PermissionItem `json:"permissions"`
}
