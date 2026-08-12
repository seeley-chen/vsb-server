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

type DepartmentRepo struct {
	collection *mongo.Collection
}

func NewDepartmentRepo(db *mongo.Database) *DepartmentRepo {
	return &DepartmentRepo{
		collection: db.Collection("departments"),
	}
}

// EnsureIndexes 确保部门集合索引
func (r *DepartmentRepo) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "department_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_department_id"),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_department_name"),
		},
	})
	return err
}

// 创建部门
func (r *DepartmentRepo) CreateDepartment(ctx context.Context, data *model.DepartmentRequest) (*model.DepartmentResponse, error) {
	now := time.Now()
	department := &model.DepartmentResponse{
		Name:         data.Name,
		Description:  data.Description,
		CreatedAt:    now,
		UpdatedAt:    now,
		DepartmentId: idgen.GenerateUuid(),
	}

	if _, err := r.collection.InsertOne(ctx, department); err != nil {
		return nil, err
	}
	return department, nil
}

// 获取部门列表
func (r *DepartmentRepo) GetDepartmentList(ctx context.Context, pageIndex, pageSize int) ([]*model.DepartmentResponse, int64, error) {
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetSkip(int64((pageIndex - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	departments := make([]*model.DepartmentResponse, 0)
	if err := cursor.All(ctx, &departments); err != nil {
		return nil, 0, err
	}
	return departments, total, nil
}

// 根据部门ID查询
func (r *DepartmentRepo) GetDepartmentById(ctx context.Context, departmentId string) (*model.DepartmentResponse, error) {
	var department model.DepartmentResponse
	err := r.collection.FindOne(ctx, bson.M{"department_id": departmentId}).Decode(&department)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &department, nil
}

// 根据部门名称查询
func (r *DepartmentRepo) FindByName(ctx context.Context, name string) (*model.DepartmentResponse, error) {
	var department model.DepartmentResponse
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&department)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &department, nil
}

// 更新部门
func (r *DepartmentRepo) UpdateDepartment(ctx context.Context, data *model.DepartmentUpdateRequest) (*model.DepartmentResponse, error) {
	update := bson.M{"updated_at": time.Now()}
	if data.Name != "" {
		update["name"] = data.Name
	}
	if data.Description != "" {
		update["description"] = data.Description
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"department_id": data.DepartmentId},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, nil
	}
	return r.GetDepartmentById(ctx, data.DepartmentId)
}

// 删除部门
func (r *DepartmentRepo) DeleteDepartment(ctx context.Context, departmentId string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"department_id": departmentId})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
