# BUG_REPRO

## Bug 是什么
前台评论列表加载根评论后没有继续加载子回复，导致 Replies 为空。

## 如何触发
在 backend 目录运行：

```bash
go test ./internal/service -run TestCommentService_ListByArticleLoadsReplies -count=1
```

## 错误信息
```
roots=1 replies=0, want 1/1
```
