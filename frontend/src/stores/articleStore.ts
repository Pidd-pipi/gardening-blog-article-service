import { useCallback, useState } from 'react';
import { listArticles } from '../api/article';
import type { Article, PageData } from '../types';

export function useArticleStore() {
  const [articles, setArticles] = useState<Article[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);

  const fetchList = useCallback(async (params: { category_id?: number; tag_id?: number; page?: number; page_size?: number }): Promise<PageData<Article>> => {
    setLoading(true);
    try {
      const data = await listArticles(params);
      setArticles(data.list);
      setTotal(data.total);
      return data;
    } finally {
      setLoading(false);
    }
  }, []);

  return { articles, total, loading, fetchList };
}
