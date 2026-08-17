# BUG_REPRO

## Bug 是什么
按标签筛选文章时，列表查询带上了标签过滤，但 total 统计查询没有带标签过滤，导致总数失真。

## 如何触发
在 backend 目录运行：

```bash
go test ./internal/service -run TestArticleService_ListPublishedByTagCount -count=1
```

## 错误信息
```
total=2 len=1, want total=1 len=1
```
