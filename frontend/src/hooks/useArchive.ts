import { useCallback, useEffect, useState } from 'react';
import { getArchives, getArticlesByMonth } from '../api/stats';
import type { ArchiveMonth, Article } from '../types';

// 文章按月归档
export function useArchive() {
  const [months, setMonths] = useState<ArchiveMonth[]>([]);
  const [articles, setArticles] = useState<Article[]>([]);
  const [current, setCurrent] = useState('');

  const loadMonths = useCallback(async () => setMonths(await getArchives()), []);

  const loadMonth = useCallback(async (month: string) => {
    setCurrent(month);
    setArticles(await getArticlesByMonth(month));
  }, []);

  useEffect(() => { loadMonths().catch(() => undefined); }, [loadMonths]);

  return { months, articles, current, loadMonth };
}
