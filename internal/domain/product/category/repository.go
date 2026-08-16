package category

import (
	"context"
	"time"

	"github.com/Vanselyn/vsb-server/internal/models"
	"github.com/Vanselyn/vsb-server/pkg/i18n"
	"github.com/Vanselyn/vsb-server/pkg/idgen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryRepo struct {
	collection *mongo.Collection
}

func NewCategoryRepo(db *mongo.Database) *CategoryRepo {
	return &CategoryRepo{
		collection: db.Collection("product_categories"),
	}
}

// 确保索引
func (r *CategoryRepo) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "category_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_category_id"),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_category_name"),
		},
	})
	return err
}

// 创建
func (r *CategoryRepo) CreateCategory(ctx context.Context, data *CategoryRequest) (*CategoryResponse, error) {
	now := time.Now()
	category := &CategoryResponse{
		Name:        data.Name,
		Description: data.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
		CategoryId:  idgen.GenerateUuid(),
		Status:      models.StatusActive,
	}

	if data.ParentId != "" {
		category.ParentId = data.ParentId
	}

	if _, err := r.collection.InsertOne(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

// 获取列表
func (r *CategoryRepo) GetCategoryList(ctx context.Context) ([]*CategoryResponse, int64, error) {
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	categories := make([]*CategoryResponse, 0)
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

// 根据ID查询
func (r *CategoryRepo) GetCategoryById(ctx context.Context, categoryId string) (*CategoryResponse, error) {
	var category CategoryResponse
	err := r.collection.FindOne(ctx, bson.M{"category_id": categoryId}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// 根据名称查询
func (r *CategoryRepo) FindByName(ctx context.Context, name i18n.Locale) (*CategoryResponse, error) {
	var category CategoryResponse
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// 更新
func (r *CategoryRepo) UpdateCategory(ctx context.Context, data *CategoryUpdateRequest) (*CategoryResponse, error) {
	update := bson.M{"updated_at": time.Now()}
	if !data.Name.IsEmpty() {
		update["name"] = data.Name
	}
	if !data.Description.IsEmpty() {
		update["description"] = data.Description
	}
	if data.ParentId != "" {
		update["parent_id"] = data.ParentId
	}
	if data.Status != "" {
		update["status"] = data.Status
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"category_id": data.CategoryId},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, nil
	}
	return r.GetCategoryById(ctx, data.CategoryId)
}

// 删除
func (r *CategoryRepo) DeleteCategory(ctx context.Context, categoryId string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"category_id": categoryId})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
