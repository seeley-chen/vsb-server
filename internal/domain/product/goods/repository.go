package goods

import (
	"context"
	"time"

	"github.com/Vanselyn/vsb-server/internal/models"
	"github.com/Vanselyn/vsb-server/pkg/idgen"
	"github.com/Vanselyn/vsb-server/pkg/mongoFilter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GoodsRepository struct {
	collection *mongo.Collection
}

func NewGoodsRepo(db *mongo.Database) *GoodsRepository {
	return &GoodsRepository{collection: db.Collection("goods")}
}

// 确保索引
func (r *GoodsRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "goods_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_goods_id"),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_name"),
		},
		{
			Keys:    bson.D{{Key: "sku", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_sku"),
		},
	})
	return err
}

// 创建商品
func (r *GoodsRepository) CreateGoods(ctx context.Context, data *GoodsRequest) (*GoodsResponse, error) {
	now := time.Now()

	good := &GoodsResponse{
		GoodsID:    idgen.GenerateUuid(),
		Name:       data.Name,
		Sku:        data.Sku,
		Remark:     data.Remark,
		Images:     data.Images,
		CategoryID: data.CategoryID,
		Specs:      r.CreateSpecsId(data.Specs),
		Status:     models.StatusUnlisted,
		Tags:       data.Tags,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if _, err := r.collection.InsertOne(ctx, good); err != nil {
		return nil, err
	}
	return good, nil
}

// 获取商品列表
func (r *GoodsRepository) GetGoodsList(ctx context.Context, pageIndex, pageSize int, query *GoodsQueryParams) ([]*GoodsResponse, int64, error) {
	init := []mongoFilter.FieldRule{
		{ParamValue: query.Name, DBKey: "name", Op: "$regex", RegexOptions: "i", SkipEmpty: true},
		{ParamValue: query.Sku, DBKey: "sku", Op: "$regex", RegexOptions: "i", SkipEmpty: true},
		{ParamValue: query.CategoryID, DBKey: "category_id", Op: "$eq", SkipEmpty: true},
		{ParamValue: query.Status, DBKey: "status", Op: "$eq", SkipEmpty: true},
	}

	filter := mongoFilter.BuildFilter(init)

	total, err := r.collection.CountDocuments(ctx, filter)

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

	goods := make([]*GoodsResponse, 0)
	if err := cursor.All(ctx, &goods); err != nil {
		return nil, 0, err
	}

	return goods, total, nil
}

// 根据ID获取商品
func (r *GoodsRepository) GetGoodsById(ctx context.Context, goodsId string) (*GoodsResponse, error) {
	var good GoodsResponse

	err := r.collection.FindOne(ctx, bson.M{"goods_id": goodsId}).Decode(&good)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &good, nil
}

// 根据SKU获取商品
func (r *GoodsRepository) GetGoodsBySku(ctx context.Context, sku string) (*GoodsResponse, error) {
	var good GoodsResponse

	err := r.collection.FindOne(ctx, bson.M{"sku": sku}).Decode(&good)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &good, nil
}

// 更新商品
func (r *GoodsRepository) UpdateGoods(ctx context.Context, data *GoodsUpdateRequest) (*GoodsResponse, error) {
	update := mongoFilter.BuildUpdate([]mongoFilter.UpdateRule{
		{Value: data.Name, DBKey: "name"},
		{Value: data.Sku, DBKey: "sku"},
		{Value: data.Remark, DBKey: "remark"},
		{Value: data.Images, DBKey: "images"},
		{Value: data.CategoryID, DBKey: "category_id"},
		{Value: r.CreateSpecsId(data.Specs), DBKey: "specs"},
		{Value: data.Status, DBKey: "status"},
		{Value: data.Tags, DBKey: "tags"},
	})

	result, err := r.collection.UpdateOne(ctx, bson.M{"goods_id": data.GoodsID}, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, nil
	}
	return r.GetGoodsById(ctx, data.GoodsID)
}

// 删除商品
func (r *GoodsRepository) DeleteGoods(ctx context.Context, goodsId string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"goods_id": goodsId})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// 创建子规格ID
// 循环specs，如果specsId为空,则specsId = idgen.GenerateUuid()
// 返回更新后的specs列表
func (r *GoodsRepository) CreateSpecsId(specs []GoodsSpec) []GoodsSpec {
	for _, spec := range specs {
		if spec.SpecID == "" {
			spec.SpecID = idgen.GenerateUuid()
		}
	}
	return specs
}
