package goods

import (
	"context"
	"errors"

	"github.com/Vanselyn/vsb-server/internal/tools"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrGoodsSkuAlreadyExists = errors.New("sku already exists")
	ErrGoodsNotFound         = errors.New("goods not found")
)

type GoodsService struct {
	repo *GoodsRepository
}

func NewGoodsService(repo *GoodsRepository) *GoodsService {
	return &GoodsService{repo: repo}
}

// 创建商品
func (s *GoodsService) CreateGoods(ctx context.Context, req *GoodsRequest) (*GoodsResponse, error) {
	if err := tools.ValidateStruct(req); err != nil {
		return nil, err
	}

	if req.Sku != "" {
		existing, err := s.repo.GetGoodsBySku(ctx, req.Sku)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrGoodsSkuAlreadyExists
		}
	}

	good, err := s.repo.CreateGoods(ctx, req)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrGoodsSkuAlreadyExists
		}
		return nil, err
	}

	return good, nil
}

// 获取商品列表
func (s *GoodsService) GetGoodsList(ctx context.Context, pageIndex, pageSize int, params GoodsQueryParams) ([]*GoodsResponse, int64, error) {
	return s.repo.GetGoodsList(ctx, pageIndex, pageSize, &params)
}

// 获取商品详情
func (s *GoodsService) GetGoodsById(ctx context.Context, goodsId string) (*GoodsResponse, error) {
	return s.repo.GetGoodsById(ctx, goodsId)
}

// 更新商品
func (s *GoodsService) UpdateGoods(ctx context.Context, goodsId string, req *GoodsUpdateRequest) (*GoodsResponse, error) {
	if goodsId == "" {
		return nil, ErrGoodsNotFound
	}
	if err := tools.ValidateStruct(req); err != nil {
		return nil, err
	}

	req.GoodsID = goodsId

	if req.GoodsID != "" {
		existing, err := s.repo.GetGoodsById(ctx, req.GoodsID)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.GoodsID != goodsId {
			return nil, ErrGoodsNotFound
		}
	}

	good, err := s.repo.UpdateGoods(ctx, req)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrGoodsSkuAlreadyExists
		}
		return nil, err
	}
	if good == nil {
		return nil, ErrGoodsNotFound
	}
	return good, nil
}

// 删除商品
func (s *GoodsService) DeleteGoods(ctx context.Context, goodsId string) error {
	if goodsId == "" {
		return ErrGoodsNotFound
	}

	err := s.repo.DeleteGoods(ctx, goodsId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrGoodsNotFound
		}
		return err
	}
	return nil
}
