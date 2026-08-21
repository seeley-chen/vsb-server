package goods

import (
	"time"

	"github.com/Vanselyn/vsb-server/internal/models"
	"github.com/Vanselyn/vsb-server/pkg/i18n"
)

// GoodsRequest 商品请求参数
type GoodsRequest struct {
	Name       i18n.Locale       `json:"name" validate:"required"`                           // 商品名称（多语言）
	Sku        string            `json:"sku"`                                                // 商品 SKU
	Remark     i18n.Locale       `json:"remark"`                                             // 备注（多语言）
	Images     []string          `json:"images"`                                             // 商品图片 URL 列表
	CategoryID string            `bson:"category_id" json:"categoryId"`                      // 所属分类 ID
	Specs      []GoodsSpec       `json:"specs"`                                              // 商品规格列表
	Status     models.StatusEnum `json:"status" enums:"active,inactive" validate:"required"` // 商品状态：active 或 inactive
	Tags       []string          `json:"tags"`                                               // 商品标签
}

// GoodsResponse 商品响应
type GoodsResponse struct {
	GoodsID    string            `bson:"goods_id" json:"goodsId"`         // 商品 ID
	Name       i18n.Locale       `json:"name" `                           // 商品名称（多语言）
	Sku        string            `json:"sku"`                             // 商品 SKU
	Remark     i18n.Locale       `json:"remark"`                          // 备注（多语言）
	Images     []string          `json:"images"`                          // 商品图片 URL 列表
	CategoryID string            `bson:"category_id" json:"categoryId"`   // 所属分类 ID
	Specs      []GoodsSpec       `json:"specs" `                          // 商品规格列表
	Status     models.StatusEnum `json:"status" enums:"active,inactive" ` // 商品状态：active 或 inactive
	Tags       []string          `json:"tags"`                            // 商品标签
	CreatedAt  time.Time         `bson:"created_at" json:"createdAt"`     // 创建时间
	UpdatedAt  time.Time         `bson:"updated_at" json:"updatedAt"`     // 更新时间
}

type GoodsUpdateRequest struct {
	GoodsID    string            `bson:"goods_id" json:"-"`              // 商品 ID
	Name       i18n.Locale       `json:"name"`                           // 商品名称（多语言）
	Sku        string            `json:"sku"`                            // 商品 SKU
	Remark     i18n.Locale       `json:"remark"`                         // 备注（多语言）
	Images     []string          `json:"images"`                         // 商品图片 URL 列表
	CategoryID string            `bson:"category_id" json:"categoryId"`  // 所属分类 ID
	Specs      []GoodsSpec       `json:"specs"`                          // 商品规格列表
	Status     models.StatusEnum `json:"status" enums:"active,inactive"` // 商品状态：active 或 inactive
	Tags       []string          `json:"tags"`                           // 商品标签
}

// GoodsParams 商品列表查询参数
type GoodsQueryParams struct {
	Sku        string            `bson:"sku" json:"sku"`                 // 商品 SKU
	CategoryID string            `bson:"category_id" json:"categoryId"`  // 分类 ID
	Name       i18n.Locale       `json:"name"`                           // 商品名称（多语言）
	Status     models.StatusEnum `json:"status" enums:"active,inactive"` // 商品状态：active 或 inactive
}

// GoodsSpec 商品规格
type GoodsSpec struct {
	Price    float64 `json:"price" validate:"required"`  // 价格
	Size     string  `json:"size" validate:"required"`   // 尺寸
	Status   bool    `json:"status" validate:"required"` // 启用状态
	Weight   string  `json:"weight" `                    // 重量
	Variety  string  `json:"variety"`                    // 品种，可选
	Color    string  `json:"color"`                      // 颜色，可选
	Material string  `json:"material"`                   // 材质，可选
}
