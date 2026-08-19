package mongoFilter

// 提供 MongoDB 查询过滤器的通用构建工具

import (
	"reflect"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

// FieldRule 定义单个字段的过滤规则
type FieldRule struct {
	// ParamValue 为从请求参数中获取的值，可为字符串、数字、切片等
	ParamValue interface{}
	// DBKey 为 MongoDB 文档中的字段名
	DBKey string
	// Op 为操作符，如 "$eq", "$ne", "$gt", "$gte", "$lt", "$lte", "$in", "$nin", "$regex", "$exists"
	Op string
	// RegexOptions 当 Op 为 "$regex" 时可指定正则选项，如 "i" 表示忽略大小写
	RegexOptions string
	// SkipEmpty 为 true 时，若 ParamValue 为零值（空字符串、0、nil、空切片等）则跳过该规则
	SkipEmpty bool
}

// BuildFilter 根据一组 FieldRule 构建 bson.M 过滤器
// 参数 ctx 预留，可用于日志追踪（当前未使用）
// 返回构建好的过滤器，可直接用于 Find, CountDocuments 等方法
func BuildFilter(rules []FieldRule) bson.M {
	filter := bson.M{}
	for _, rule := range rules {
		// 如果启用了跳过空值，且 ParamValue 为零值，则跳过
		if rule.SkipEmpty && isEmpty(rule.ParamValue) {
			continue
		}

		switch rule.Op {
		case "$eq":
			filter[rule.DBKey] = rule.ParamValue
		case "$ne":
			filter[rule.DBKey] = bson.M{"$ne": rule.ParamValue}
		case "$gt":
			filter[rule.DBKey] = bson.M{"$gt": rule.ParamValue}
		case "$gte":
			filter[rule.DBKey] = bson.M{"$gte": rule.ParamValue}
		case "$lt":
			filter[rule.DBKey] = bson.M{"$lt": rule.ParamValue}
		case "$lte":
			filter[rule.DBKey] = bson.M{"$lte": rule.ParamValue}
		case "$in":
			filter[rule.DBKey] = bson.M{"$in": toSlice(rule.ParamValue)}
		case "$nin":
			filter[rule.DBKey] = bson.M{"$nin": toSlice(rule.ParamValue)}
		case "$regex":
			// 使用 regexp.QuoteMeta 防止正则注入
			pattern := regexp.QuoteMeta(toString(rule.ParamValue))
			regexOpt := rule.RegexOptions
			if regexOpt == "" {
				regexOpt = "" // 默认无选项
			}
			filter[rule.DBKey] = bson.M{
				"$regex":   pattern,
				"$options": regexOpt,
			}
		case "$exists":
			// ParamValue 应为 bool 类型，true 表示存在，false 表示不存在
			if v, ok := rule.ParamValue.(bool); ok {
				filter[rule.DBKey] = bson.M{"$exists": v}
			}
		default:
			// 未知操作符，默认当作精确匹配
			filter[rule.DBKey] = rule.ParamValue
		}
	}
	return filter
}

// isEmpty 判断值是否为「空」（零值）。
// 支持：字符串（""）、数值（含 uint 家族及具名数值类型 0）、布尔（false）、
// nil 指针/接口/chan/func、空切片/数组/映射、零值结构体等。
// 切片/数组/映射/字符串以长度为 0 判定，保留对非 nil 空集合的跳过语义。
func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		return rv.Len() == 0
	default:
		return rv.IsZero()
	}
}

// toString 将 interface{} 转为字符串
func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// toSlice 将 interface{} 转为 []interface{}，用于 $in/$nin
func toSlice(v interface{}) []interface{} {
	if slice, ok := v.([]interface{}); ok {
		return slice
	}
	// 如果不是切片，包装成单元素切片
	return []interface{}{v}
}
