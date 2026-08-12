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

type RoleRepo struct {
	collection *mongo.Collection
}

func NewRoleRepo(db *mongo.Database) *RoleRepo {
	return &RoleRepo{
		collection: db.Collection("roles"),
	}
}

// EnsureIndexes 确保角色集合索引
func (r *RoleRepo) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "role_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_role_id"),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_role_name"),
		},
	})
	return err
}

func (r *RoleRepo) CreateRole(ctx context.Context, data *model.RoleRequest) (*model.RoleResponse, error) {
	now := time.Now()
	permissions := data.Permissions
	if permissions == nil {
		permissions = []model.PermissionItem{}
	}

	role := &model.RoleResponse{
		RoleId:      idgen.GenerateUuid(),
		Name:        data.Name,
		Description: data.Description,
		Permissions: permissions,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if _, err := r.collection.InsertOne(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

// 获取角色列表
func (r *RoleRepo) GetRoleList(ctx context.Context, pageIndex, pageSize int) ([]*model.RoleResponse, int64, error) {
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

	roles := make([]*model.RoleResponse, 0)
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

// 根据角色ID获取角色
func (r *RoleRepo) GetRoleById(ctx context.Context, roleId string) (*model.RoleResponse, error) {
	var role model.RoleResponse
	err := r.collection.FindOne(ctx, bson.M{"role_id": roleId}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) FindByName(ctx context.Context, name string) (*model.RoleResponse, error) {
	var role model.RoleResponse
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) UpdateRole(ctx context.Context, data *model.RoleUpdateRequest) (*model.RoleResponse, error) {
	update := bson.M{"updated_at": time.Now()}
	if data.Name != "" {
		update["name"] = data.Name
	}
	if data.Description != "" {
		update["description"] = data.Description
	}
	if data.Permissions != nil {
		update["permissions"] = data.Permissions
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"role_id": data.RoleId},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, nil
	}
	return r.GetRoleById(ctx, data.RoleId)
}

func (r *RoleRepo) DeleteRole(ctx context.Context, roleId string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"role_id": roleId})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
