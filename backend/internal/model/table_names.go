package model

// 显式指定表名，确保与 database/init.sql 中的建表名一致。

func (User) TableName() string         { return "users" }
func (Category) TableName() string     { return "categories" }
func (Tag) TableName() string          { return "tags" }
func (Article) TableName() string      { return "articles" }
func (Comment) TableName() string      { return "comments" }
func (ArticleTag) TableName() string   { return "article_tags" }
