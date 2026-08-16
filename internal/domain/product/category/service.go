package category

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category exists")
)

type CategoryService struct {
	repo *CategoryRepo
}

func NewCategoryService(repo *CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

// 创建
// 直接靠唯一索引判重（InsertOne 的 duplicate key 错误），不使用前置 FindByName 查重：
// 避免查重与插入之间的竞态（并发请求可能同时通过查重），且唯一索引本身已保证唯一性。
func (s *CategoryService) CreateCategory(ctx context.Context, data *CategoryRequest) (*CategoryResponse, error) {
	data.Name = data.Name.TrimSpace()
	data.Description = data.Description.TrimSpace()

	category, err := s.repo.CreateCategory(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrCategoryExists
		}
		return nil, err
	}

	return s.repo.GetCategoryById(ctx, category.CategoryId)
}

// 查询列表
func (s *CategoryService) CategoryList(ctx context.Context) ([]*CategoryResponse, int64, error) {
	return s.repo.GetCategoryList(ctx)
}

// 更新
func (s *CategoryService) UpdateCategory(ctx context.Context, categoryId string, data *CategoryUpdateRequest) (*CategoryResponse, error) {
	if categoryId == "" {
		return nil, ErrCategoryNotFound
	}

	data.Description = data.Description.TrimSpace()
	data.CategoryId = categoryId

	if !data.Name.IsEmpty() {
		existing, err := s.repo.FindByName(ctx, data.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.CategoryId != categoryId {
			return nil, ErrCategoryExists
		}
	}

	category, err := s.repo.UpdateCategory(ctx, data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrCategoryExists
		}
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, categoryId string) error {
	if categoryId == "" {
		return ErrCategoryNotFound
	}

	err := s.repo.DeleteCategory(ctx, categoryId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrCategoryNotFound
		}
		return err
	}
	return nil
}
