import request from '../utils/request';
import type { Tag } from '../types';

export function listTags(): Promise<Tag[]> {
  return request.get('/tags');
}

export function createTag(data: Partial<Tag>): Promise<Tag> {
  return request.post('/tags', data);
}

export function updateTag(id: number, data: Partial<Tag>): Promise<Tag> {
  return request.put(`/tags/${id}`, data);
}

export function deleteTag(id: number): Promise<void> {
  return request.delete(`/tags/${id}`);
}
