package models

type StatusEnum string

const (
	// 上架
	StatusActive StatusEnum = "active"
	// 未激活
	StatusInactive StatusEnum = "inactive"
	// 已删除
	StatusDeleted StatusEnum = "deleted"
	// 已过期
	StatusExpired StatusEnum = "expired"
	// 停用
	StatusDormant StatusEnum = "dormant"
	// 下架
	StatusDown StatusEnum = "down"
	// 禁用
	StatusDisabled StatusEnum = "disabled"
	// 启用
	StatusEnabled StatusEnum = "enabled"
)
