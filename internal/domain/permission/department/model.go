package department

import "time"

// 创建部门请求
type DepartmentRequest struct {
	Name        string `json:"name" validate:"required"` // 部门名称
	Description string `json:"description"`              // 部门描述
}

// 部门响应
type DepartmentResponse struct {
	Name         string    `json:"name" validate:"required"`                              // 部门名称
	Description  string    `bson:"description" json:"description"`                        // 部门描述
	DepartmentId string    `bson:"department_id" json:"departmentId" validate:"required"` // 部门ID
	CreatedAt    time.Time `bson:"created_at" json:"createdAt" validate:"required"`       // 创建时间
	UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt" validate:"required"`       // 更新时间
}

// 更新部门请求（空字段表示不修改）
type DepartmentUpdateRequest struct {
	DepartmentId string `json:"-"`           // 部门ID，以路径参数为准
	Name         string `json:"name"`        // 部门名称，空表示不修改
	Description  string `json:"description"` // 部门描述
}
