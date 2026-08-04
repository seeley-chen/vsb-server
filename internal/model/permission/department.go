package model

import (
	"time"

	"github.com/seeley-chen/vsb-server/pkg/idgen"
)

type DepartmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Department struct {
	DepartmentId string    `bson:"departmentId"`
	Name         string    `bson:"name"`
	Description  string    `bson:"description"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// 创建新部门
func NewDepartment(req *DepartmentRequest) *Department {
	now := time.Now()

	return &Department{
		DepartmentId: idgen.GenerateUuid(),
		Name:         req.Name,
		Description:  req.Description,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
