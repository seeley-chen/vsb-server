package goods

type GoodsService struct {
	repo *GoodsRepository
}

func NewGoodsService(repo *GoodsRepository) *GoodsService {
	return &GoodsService{repo: repo}
}
