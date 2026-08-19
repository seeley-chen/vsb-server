package user

import (
	"context"
	"time"

	"github.com/Vanselyn/vsb-server/internal/models"
	"github.com/Vanselyn/vsb-server/pkg/idgen"
	"github.com/Vanselyn/vsb-server/pkg/mongoFilter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		collection: db.Collection("users"),
	}
}

func (r *UserRepo) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_user_id"),
		},
		{
			Keys:    bson.D{{Key: "account", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_account"),
		},
	})
	return err
}

// 创建用户
func (r *UserRepo) CreateUser(ctx context.Context, data *UserRequest) (*UserResponse, error) {
	now := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &UserResponse{
		UserId:       idgen.GenerateUniqueID(),
		Username:     data.Username,
		Account:      data.Account,
		Identity:     data.Identity,
		PasswordHash: string(hashedPassword),
		Email:        data.Email,
		Phone:        data.Phone,
		RoleId:       data.RoleId,
		DepartmentId: data.DepartmentId,
		Status:       models.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if _, err := r.collection.InsertOne(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// 获取用户列表（username 模糊搜索，account/email 精确匹配，均可选）
func (r *UserRepo) GetUserList(ctx context.Context, pageIndex, pageSize int, params UserQueryParams) ([]*UserResponse, int64, error) {
	init := []mongoFilter.FieldRule{
		{ParamValue: params.Username, DBKey: "username", Op: "$regex", RegexOptions: "i", SkipEmpty: true},
		{ParamValue: params.Account, DBKey: "account", Op: "$regex", RegexOptions: "i", SkipEmpty: true},
		{ParamValue: params.Email, DBKey: "email", Op: "$regex", RegexOptions: "i", SkipEmpty: true},
		{ParamValue: params.RoleID, DBKey: "role_id", Op: "$eq", SkipEmpty: true},
		{ParamValue: params.DepartmentID, DBKey: "department_id", Op: "$eq", SkipEmpty: true},
		{ParamValue: params.Status, DBKey: "status", Op: "$eq", SkipEmpty: true},
	}

	filter := mongoFilter.BuildFilter(init)

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetSkip(int64((pageIndex - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSort(bson.D{{Key: "name", Value: 1}}).
		SetProjection(bson.M{"passwordHash": 0}).
		SetCollation(&options.Collation{
			Locale:   "zh", // 中文 locale
			Strength: 1,    // 1 表示忽略大小写和变音符号，常用
		})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close(ctx)

	users := make([]*UserResponse, 0)
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// 根据用户ID获取用户
func (r *UserRepo) GetUserById(ctx context.Context, userId string) (*UserResponse, error) {
	var user UserResponse
	err := r.collection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 根据账号获取用户
func (r *UserRepo) GetUserByAccount(ctx context.Context, account string) (*UserResponse, error) {
	var user UserResponse
	err := r.collection.FindOne(ctx, bson.M{"account": account}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 根据邮箱获取用户
func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*UserResponse, error) {
	var user UserResponse
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 根据手机号获取用户
func (r *UserRepo) GetUserByPhone(ctx context.Context, phone string) (*UserResponse, error) {
	var user UserResponse
	err := r.collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// 更新用户（空字段表示不修改）
func (r *UserRepo) UpdateUser(ctx context.Context, data *UserUpdateRequest) (*UserResponse, error) {
	// 密码需先哈希；空则保持空字符串，BuildUpdate 会自动跳过
	var hashedPassword string
	if data.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashedPassword = string(hashed)
	}

	update := mongoFilter.BuildUpdate([]mongoFilter.UpdateRule{
		{Value: data.Username, DBKey: "username"},
		{Value: data.Status, DBKey: "status"},
		{Value: hashedPassword, DBKey: "passwordHash"},
		{Value: data.Email, DBKey: "email"},
		{Value: data.Phone, DBKey: "phone"},
		{Value: data.RoleId, DBKey: "role_id"},
		{Value: data.DepartmentId, DBKey: "department_id"},
		{Value: data.Identity, DBKey: "identity"},
	})

	if _, err := r.collection.UpdateOne(ctx, bson.M{"user_id": data.UserId}, bson.M{"$set": update}); err != nil {
		return nil, err
	}
	return r.GetUserById(ctx, data.UserId)
}

// 删除用户
func (r *UserRepo) DeleteUser(ctx context.Context, userId string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userId})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
