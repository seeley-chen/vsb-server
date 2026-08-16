package category

import (
	"time"

	"github.com/Vanselyn/vsb-server/internal/models"
	"github.com/Vanselyn/vsb-server/pkg/i18n"
)

type CategoryRequest struct {
	Name        i18n.Locale `json:"name" validate:"required"`
	Description i18n.Locale `json:"description"`
	ParentId    string      `json:"parentId"`
}

type CategoryResponse struct {
	CategoryId  string            `json:"categoryId" bson:"category_id"`
	Name        i18n.Locale       `json:"name" bson:"name"`
	Description i18n.Locale       `json:"description" bson:"description"`
	ParentId    string            `json:"parentId" bson:"parent_id"`
	Status      models.StatusEnum `json:"status"`
	CreatedAt   time.Time         `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updatedAt" bson:"updated_at"`
}

type CategoryUpdateRequest struct {
	CategoryId  string            `json:"categoryId"`
	Name        i18n.Locale       `json:"name"`
	Description i18n.Locale       `json:"description"`
	ParentId    string            `json:"parentId"`
	Status      models.StatusEnum `json:"status"`
}
