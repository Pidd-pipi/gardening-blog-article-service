// 与 backend/internal/constants/comment.go 保持一致
export const CommentStatus = {
  PENDING: 'pending',
  APPROVED: 'approved',
  DELETED: 'deleted',
} as const;

export const CommentStatusLabels: Record<string, string> = {
  [CommentStatus.PENDING]: '待审核',
  [CommentStatus.APPROVED]: '已通过',
  [CommentStatus.DELETED]: '已删除',
};
