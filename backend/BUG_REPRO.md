# BUG_REPRO

## Bug 是什么
文章全文搜索只查了标题字段，没有把摘要和正文纳入匹配范围。

## 如何触发
在 backend 目录运行：

```bash
go test ./internal/service -run TestArticleService_SearchSummary -count=1
```

## 错误信息
```
total=0 len=0, want 1/1
```
