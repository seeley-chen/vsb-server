package models

// StatusEnum 资源状态
// active=上架，inactive=未激活，deleted=已删除，expired=已过期，dormant=停用，down=下架，disabled=禁用，enabled=启用
type StatusEnum string

const (
	StatusActive   StatusEnum = "active"   // 活跃中
	StatusUnlisted StatusEnum = "unlisted" // 未上架
	StatusListed   StatusEnum = "listed"   // 已上架
	StatusInactive StatusEnum = "inactive" // 已下架
	StatusDeleted  StatusEnum = "deleted"  // 已删除
	StatusExpired  StatusEnum = "expired"  // 已过期
	StatusPending  StatusEnum = "pending"  // 待审核
	StatusApproved StatusEnum = "approved" // 审核通过
	StatusDisabled StatusEnum = "disabled" // 禁用
	StatusEnabled  StatusEnum = "enabled"  // 启用
)
