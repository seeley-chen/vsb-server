package permission

import "time"

type UserRequest struct {
	Username     string    `json:"username" validate:"required"`
	Account      string    `json:"account" validate:"required"`
	Password     string    `json:"password" validate:"required"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	RoleId       string    `bson:"role_id" json:"roleId" validate:"required"`
	DepartmentId string    `bson:"department_id" json:"departmentId" validate:"required"`
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt"`
}

type UserResponse struct {
	UserId       string    `bson:"user_id" json:"userId"`
	Username     string    `bson:"username" json:"username"`
	Account      string    `bson:"account" json:"account"`
	PasswordHash string    `bson:"passwordHash" json:"-"`
	Email        string    `bson:"email" json:"email"`
	Phone        string    `bson:"phone" json:"phone"`
	RoleId       string    `bson:"role_id" json:"roleId"`
	DepartmentId string    `bson:"department_id" json:"departmentId"`
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updatedAt"`
}

type UserUpdateRequest struct {
	UserId       string `bson:"user_id" json:"userId" validate:"required"`
	Username     string `json:"username" validate:"required"`
	Account      string `json:"account" validate:"required"`
	Password     string `bson:"Password" json:"-"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	RoleId       string `bson:"role_id" json:"roleId" validate:"required"`
	DepartmentId string `bson:"department_id" json:"departmentId" validate:"required"`
}

type UserLoginRequest struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}
