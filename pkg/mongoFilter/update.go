package mongoFilter

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateRule 定义单个字段在 $set 更新中的规则
type UpdateRule struct {
	// Value 字段值；零值（空字符串、0、nil、空切片等）将被跳过，符合「空字段表示不修改」语义
	Value interface{}
	// DBKey MongoDB 文档字段名
	DBKey string
}

// BuildUpdate 根据一组 UpdateRule 构建可直接用于 bson.M{"$set": ...} 的更新文档
// 零值字段会被自动跳过，调用方无需逐个判空
func BuildUpdate(rules []UpdateRule) bson.M {
	update := bson.M{"updated_at": time.Now()}
	for _, rule := range rules {
		if isEmpty(rule.Value) {
			continue
		}
		update[rule.DBKey] = rule.Value
	}
	return update
}
