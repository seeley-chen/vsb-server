package department

import (
	"context"
	"regexp"
	"time"

	"github.com/Vanselyn/vsb-server/pkg/idgen"
	"github.com/Vanselyn/vsb-server/pkg/mongoFilter"
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
func (r *DepartmentRepo) CreateDepartment(ctx context.Context, data *DepartmentRequest) (*DepartmentResponse, error) {
	now := time.Now()
	department := &DepartmentResponse{
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
func (r *DepartmentRepo) GetDepartmentList(ctx context.Context, pageIndex, pageSize int, name string) ([]*DepartmentResponse, int64, error) {
	filter := bson.M{}
	if name != "" {
		filter["name"] = bson.M{
			"$regex":   regexp.QuoteMeta(name),
			"$options": "i",
		}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetSkip(int64((pageIndex - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	departments := make([]*DepartmentResponse, 0)
	if err := cursor.All(ctx, &departments); err != nil {
		return nil, 0, err
	}
	return departments, total, nil
}

// 根据部门ID查询
func (r *DepartmentRepo) GetDepartmentById(ctx context.Context, departmentId string) (*DepartmentResponse, error) {
	var department DepartmentResponse
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
func (r *DepartmentRepo) FindByName(ctx context.Context, name string) (*DepartmentResponse, error) {
	var department DepartmentResponse
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
func (r *DepartmentRepo) UpdateDepartment(ctx context.Context, data *DepartmentUpdateRequest) (*DepartmentResponse, error) {
	update := mongoFilter.BuildUpdate([]mongoFilter.UpdateRule{
		{Value: data.Name, DBKey: "name"},
		{Value: data.Description, DBKey: "description"},
	})

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
