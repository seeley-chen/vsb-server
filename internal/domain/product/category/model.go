package category

import (
	"time"

	"github.com/Vanselyn/vsb-server/internal/models"
	"github.com/Vanselyn/vsb-server/pkg/i18n"
)

// CategoryRequest 创建商品分类请求
type CategoryRequest struct {
	Name        i18n.Locale `json:"name" validate:"required"`  // 分类名称（多语言）
	Description i18n.Locale `json:"description"`               // 分类描述（多语言）
	ParentId    string      `bson:"parent_id" json:"parentId"` // 父分类 ID，为空表示顶级分类
}

// CategoryResponse 商品分类响应 / 存储模型
type CategoryResponse struct {
	CategoryId  string            `bson:"category_id" json:"categoryId"` // 分类 ID
	Name        i18n.Locale       `json:"name"`                          // 分类名称（多语言）
	Description i18n.Locale       `json:"description"`                   // 分类描述（多语言）
	ParentId    string            `bson:"parent_id" json:"parentId"`     // 父分类 ID，为空表示顶级分类
	Status      models.StatusEnum `json:"status"`                        // 分类状态
	CreatedAt   time.Time         `json:"createdAt"`                     // 创建时间
	UpdatedAt   time.Time         `json:"updatedAt"`                     // 更新时间
}

// CategoryUpdateRequest 更新商品分类请求（空字段表示不修改）
type CategoryUpdateRequest struct {
	CategoryId  string            `bson:"category_id" json:"categoryId"` // 分类 ID，以路径参数为准
	Name        i18n.Locale       `json:"name"`                          // 分类名称（多语言）
	Description i18n.Locale       `json:"description"`                   // 分类描述（多语言）
	ParentId    string            `json:"parentId"`                      // 父分类 ID
	Status      models.StatusEnum `json:"status"`                        // 分类状态
}
