import request from '../utils/request';
import type { Article, PageData } from '../types';

export function listArticles(params: { category_id?: number; tag_id?: number; keyword?: string; page?: number; page_size?: number }): Promise<PageData<Article>> {
  return request.get('/articles', { params });
}

export function getArticle(slug: string): Promise<Article> {
  return request.get(`/articles/${slug}`);
}

export function listAdminArticles(params: { status?: string; page?: number; page_size?: number }): Promise<PageData<Article>> {
  return request.get('/admin/articles', { params });
}

export function getAdminArticle(id: number): Promise<Article> {
  return request.get(`/admin/articles/${id}`);
}

export function createArticle(data: Partial<Article> & { tag_ids?: number[] }): Promise<Article> {
  return request.post('/articles', data);
}

export function updateArticle(id: number, data: Partial<Article> & { tag_ids?: number[] }): Promise<Article> {
  return request.put(`/articles/${id}`, data);
}

export function deleteArticle(id: number): Promise<void> {
  return request.delete(`/articles/${id}`);
}

export function publishArticle(id: number): Promise<Article> {
  return request.post(`/articles/${id}/publish`);
}

export function topArticles(): Promise<Article[]> {
  return request.get('/articles/top');
}

export function searchArticles(q: string, params: { page?: number; page_size?: number } = {}): Promise<PageData<Article>> {
  return request.get('/search', { params: { q, ...params } });
}
