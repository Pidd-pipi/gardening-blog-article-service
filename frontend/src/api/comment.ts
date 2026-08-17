import request from '../utils/request';
import type { Comment, PageData } from '../types';

export function listComments(articleId: number): Promise<Comment[]> {
  return request.get('/comments', { params: { article_id: articleId } });
}

export function createComment(data: { article_id: number; parent_id?: number; nickname?: string; email?: string; content: string }): Promise<Comment> {
  return request.post('/comments', data);
}

export function listAdminComments(params: { status?: string; page?: number; page_size?: number }): Promise<PageData<Comment>> {
  return request.get('/admin/comments', { params });
}

export function updateCommentStatus(id: number, status: string): Promise<void> {
  return request.put(`/comments/${id}/status`, { status });
}
