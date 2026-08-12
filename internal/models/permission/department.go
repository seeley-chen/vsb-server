package permission

import "time"

// DepartmentRequest 创建部门请求
type DepartmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DepartmentResponse 部门响应 / 存储模型
type DepartmentResponse struct {
	Name         string    `bson:"name" json:"name"`                  // 部门名称
	Description  string    `bson:"description" json:"description"`    // 部门描述
	DepartmentId string    `bson:"department_id" json:"departmentId"` // 部门ID
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`       // 创建时间
	UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt"`       // 更新时间
}

// DepartmentUpdateRequest 更新部门请求（空字段表示不修改）
type DepartmentUpdateRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DepartmentId string `json:"-"`
}
