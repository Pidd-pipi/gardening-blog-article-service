# BUG_REPRO

## Bug 是什么
文章标签替换时直接逐个插入 ArticleTag，没有去重；TagIDs 里出现重复项时违反复合主键唯一约束，事务整体回滚。

## 如何触发
在 backend 目录运行：

```bash
go test ./internal/service -run TestArticleService_UpdateDeduplicatesTagIDs -count=1
```

测试会对文章更新两次相同的 tag id。

## 错误信息
```
UNIQUE constraint failed: article_tags.article_id, article_tags.tag_id
```
