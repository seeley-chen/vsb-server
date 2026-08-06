package permission

import (
	"context"
	"time"

	model "github.com/seeley-chen/vsb-server/internal/models/permission"
	"github.com/seeley-chen/vsb-server/pkg/idgen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 初始化结构体
type DepartmentRepo struct {
	collection *mongo.Collection
}

// 初始化部门仓库
func NewDepartmentRepo(db *mongo.Database) *DepartmentRepo {
	return &DepartmentRepo{
		collection: db.Collection("departments"),
	}
}

// 创建部门
func (r *DepartmentRepo) CreateDepartment(ctx context.Context, data *model.DepartmentRequest) (*model.Department, error) {
	now := time.Now()

	department := &model.Department{
		Name:         data.Name,
		Description:  data.Description,
		CreatedAt:    now,
		UpdatedAt:    now,
		DepartmentId: idgen.GenerateUuid(),
	}

	_, err := r.collection.InsertOne(ctx, department)

	if err != nil {
		return nil, err
	}

	return department, nil
}

// GetDepartmentList 分页获取部门列表
func (r *DepartmentRepo) GetDepartmentList(ctx context.Context, pageIndex, pageSize int) ([]*model.Department, int64, error) {
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetSkip(int64((pageIndex - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	departments := make([]*model.Department, 0)
	if err := cursor.All(ctx, &departments); err != nil {
		return nil, 0, err
	}

	return departments, total, nil
}

// 根据部门ID查询
func (r *DepartmentRepo) GetDepartmentById(ctx context.Context, departmentId string) (*model.Department, error) {
	var department model.Department

	err := r.collection.FindOne(ctx, bson.M{"department_id": departmentId}).Decode(&department)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &department, nil
}

// 更新部门
func (r *DepartmentRepo) UpdateDepartment(ctx context.Context, data *model.DepartmentUpdateRequest) (*model.Department, error) {
	update := bson.M{}

	if data.Name != "" {
		update["name"] = data.Name
	}

	if data.Description != "" {
		update["description"] = data.Description
	}

	update["updated_at"] = time.Now()

	_, err := r.collection.UpdateOne(ctx, bson.M{"department_id": data.DepartmentId}, bson.M{"$set": update})

	if err != nil {
		return nil, err
	}

	return r.GetDepartmentById(ctx, data.DepartmentId)
}

// 删除部门
func (r *DepartmentRepo) DeleteDepartment(ctx context.Context, departmentId string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"department_id": departmentId})

	if err != nil {
		return err
	}

	return nil
}
