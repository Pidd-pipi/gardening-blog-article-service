// 与 backend/internal/constants/article.go 保持一致
export const ArticleStatus = {
  DRAFT: 'draft',
  PUBLISHED: 'published',
  SCHEDULED: 'scheduled',
} as const;

export const ArticleStatusLabels: Record<string, string> = {
  [ArticleStatus.DRAFT]: '草稿',
  [ArticleStatus.PUBLISHED]: '已发布',
  [ArticleStatus.SCHEDULED]: '定时发布',
};
