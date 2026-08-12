package permission

import "time"

// 创建-部门请求模型
type DepartmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// 部门模型
type DepartmentResponse struct {
	Name         string    `bson:"name" json:"name"`                  // 部门名称
	Description  string    `bson:"description" json:"description"`    // 部门描述
	DepartmentId string    `bson:"department_id" json:"departmentId"` // 部门ID
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`       // 创建时间
	UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt"`       // 更新时间
}

// 更新-部门请求模型
type DepartmentUpdateRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DepartmentId string `json:"-"`
}
