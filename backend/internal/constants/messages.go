package constants

// messages.go 集中管理前端提示文案、后端返回文案与日志文案（屎山耦合点）。

const (
	MsgOK                = "ok"
	MsgLoginSuccess      = "登录成功"
	MsgRegisterSuccess   = "注册成功"
	MsgCreated           = "创建成功"
	MsgUpdated           = "更新成功"
	MsgDeleted           = "删除成功"
	MsgCommentSubmitted  = "评论提交成功，待审核"
	MsgCommentApproved   = "评论已通过"
	MsgArticlePublished  = "文章已发布"
	MsgArticleDrafted    = "草稿已保存"
	MsgDuplicateEmail    = "邮箱（User.email）已注册"
	MsgLoginFailed       = "邮箱或密码（User）错误"
	MsgUserNotFound      = "用户（User）不存在"
	MsgRoleForbidden     = "当前角色（UserRole）无权执行该操作"
	MsgArticleNotFound   = "文章（Article）不存在"
	MsgCategoryNotFound  = "分类（Category）不存在"
	MsgParentCategoryNotFound = "父分类（Category.parent_id）不存在"
	MsgTagNotFound       = "标签（Tag）不存在"
	MsgSlugConflict      = "别名（slug）已存在"
	MsgCommentNotFound   = "评论（Comment）不存在"
	MsgInternalError     = "服务内部错误"
	MsgParamInvalid      = "请求参数校验失败"
	MsgRateLimited       = "请求过于频繁，请稍后重试"
)
