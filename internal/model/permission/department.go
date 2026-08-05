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
	DepartmentId string    `bson:"departmentId" json:"departmentId"`
	Name         string    `bson:"name" json:"name"`
	Description  string    `bson:"description" json:"description"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt" json:"updatedAt"`
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
