package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	model "github.com/seeley-chen/vsb-server/internal/model/permission"
)

// Repository 用户数据访问层
type Repository struct {
	collection *mongo.Collection
}

// NewRepository 创建 Repository 实例
func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection("users"),
	}
}

// Create 插入新用户
func (r *Repository) Create(ctx context.Context, user *model.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err // 如果插入失败，则返回错误
}

// FindByAccount 根据用户名查找用户
func (r *Repository) FindByAccount(ctx context.Context, account string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"account": account}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据用户 ID 查找用户
func (r *Repository) FindByID(ctx context.Context, userId string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"userId": userId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll 分页查询用户列表（不返回密码哈希）
func (r *Repository) FindAll(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	// 计算总数
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err // 如果计算总数失败，则返回错误
	}

	// 分页查询，排除 password_hash 字段
	findOptions := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetProjection(bson.M{"passwordHash": 0})

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err // 如果分页查询失败，则返回错误
	}
	defer cursor.Close(ctx)

	var users []model.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err // 如果获取用户列表失败，则返回错误
	}

	return users, total, nil // 返回用户列表和总数
}

// EnsureIndexes 创建必要索引（account 唯一索引）
func (r *Repository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"account": 1},
		Options: options.Index().SetUnique(true),
	})
	return err
}

// DeleteById 根据用户ID删除用户
func (r *Repository) DeleteUserById(ctx context.Context, userId string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"userId": userId})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments // 没找到要删除的用户
	}

	return nil
}

// UpdateUserById 根据用户ID更新用户
func (r *Repository) UpdateUserById(ctx context.Context, userId string, update *model.User) error {
	setFields := bson.M{
		"username":  update.Username,
		"email":     update.Email,
		"account":   update.Account,
		"phone":     update.Phone,
		"updatedAt": time.Now(),
	}

	if update.PasswordHash != "" {
		setFields["passwordHash"] = update.PasswordHash
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"userId": userId}, bson.M{"$set": setFields})

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments // 没找到要更新的用户
	}

	return nil
}
