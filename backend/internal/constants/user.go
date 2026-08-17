package constants

// UserRole 用户角色枚举（README 枚举出现位置清单必列）。
const (
	RoleAdmin   = "admin"   // 博主管理员
	RoleVisitor = "visitor" // 游客
)

// Roles 全部角色。
var Roles = []string{RoleAdmin, RoleVisitor}
