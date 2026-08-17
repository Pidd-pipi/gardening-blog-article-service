# BUG_REPRO

## Bug 是什么
分类树组装时父分类匹配条件写反，导致子分类不会挂到父分类下。

## 如何触发
在 backend 目录运行：

```bash
go test ./internal/service -run TestCategoryTree_AttachesChildren -count=1
```

## 错误信息
```
root children = 0, want 1
```
