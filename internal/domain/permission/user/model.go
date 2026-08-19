package user

import (
	"time"

	"github.com/Vanselyn/vsb-server/internal/models"
)

// 创建用户请求
type UserRequest struct {
	Username     string            `json:"username" validate:"required" `                         // 用户名
	Account      string            `json:"account" validate:"required"`                           // 账号
	Password     string            `json:"password" validate:"required"`                          // 密码
	Identity     string            `json:"identity" validate:"required" enums:"admin,user"`       // 身份标识
	RoleId       string            `bson:"role_id" json:"roleId" validate:"required"`             // 角色ID
	DepartmentId string            `bson:"department_id" json:"departmentId" validate:"required"` // 部门ID
	Email        string            `json:"email"`                                                 // 邮箱
	Phone        string            `json:"phone"`
	Status       models.StatusEnum `json:"status" enums:"active,inactive"` // 用户状态
	CreatedAt    time.Time         `bson:"created_at" json:"createdAt"`    // 创建时间，由服务端生成
	UpdatedAt    time.Time         `bson:"updated_at" json:"updatedAt"`    // 更新时间，由服务端生成
}

// 用户响应 / 存储模型
type UserResponse struct {
	UserId       string            `bson:"user_id" json:"userId" validate:"required"`                        // 用户ID
	Username     string            `json:"username" validate:"required"`                                     // 用户名
	Account      string            `json:"account" validate:"required"`                                      // 账号
	PasswordHash string            `bson:"passwordHash" json:"-"`                                            // 密码哈希，不对外返回
	Identity     string            `json:"identity" validate:"required"`                                     // 身份标识
	RoleId       string            `bson:"role_id" json:"roleId" validate:"required"`                        // 角色ID
	DepartmentId string            `bson:"department_id" json:"departmentId" validate:"required"`            // 部门ID
	Status       models.StatusEnum `bson:"status" json:"status" enums:"active,inactive" validate:"required"` // 用户状态
	Email        string            `json:"email"`                                                            // 邮箱
	Phone        string            `json:"phone"`                                                            // 手机号
	CreatedAt    time.Time         `bson:"created_at" json:"createdAt" validate:"required"`                  // 创建时间
	UpdatedAt    time.Time         `bson:"updated_at" json:"updatedAt" validate:"required"`                  // 更新时间
}

// 更新用户请求（空字段表示不修改）
type UserUpdateRequest struct {
	UserId       string            `json:"-"`                              // 用户ID，以路径参数为准
	Username     string            `json:"username"`                       // 用户名，空表示不修改
	Account      string            `json:"-"`                              // 登录账号
	Password     string            `json:"password"`                       // 登录密码
	Identity     string            `json:"identity" enums:"admin,user"`    // 身份标识
	Email        string            `json:"email"`                          // 邮箱
	Phone        string            `json:"phone"`                          // 手机号
	RoleId       string            `json:"roleId"`                         // 角色ID
	DepartmentId string            `json:"departmentId"`                   // 部门ID
	Status       models.StatusEnum `json:"status" enums:"active,inactive"` // 用户状态
}

// 用户查询参数
type UserQueryParams struct {
	Username     string            // 用户名
	Account      string            // 账号
	Email        string            // 邮箱
	RoleID       string            // 角色ID
	DepartmentID string            // 部门ID
	Status       models.StatusEnum // 用户状态
}
