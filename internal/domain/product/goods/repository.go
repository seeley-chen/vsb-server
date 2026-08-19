package goods

import "go.mongodb.org/mongo-driver/mongo"

type GoodsRepository struct {
	collection *mongo.Collection
}

func NewGoodsRepository(db *mongo.Database) *GoodsRepository {
	return &GoodsRepository{collection: db.Collection("goods")}
}
