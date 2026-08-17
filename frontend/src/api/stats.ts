import request from '../utils/request';
import type { ArchiveMonth, Article, SiteStats } from '../types';

export function getSiteStats(): Promise<SiteStats> {
  return request.get('/stats');
}

export function getArchives(): Promise<ArchiveMonth[]> {
  return request.get('/archives');
}

export function getArticlesByMonth(month: string): Promise<Article[]> {
  return request.get('/archives/articles', { params: { month } });
}
