import { useCallback, useEffect, useState } from 'react';
import { getSiteStats } from '../api/stats';
import type { SiteStats } from '../types';

// 站点文章/分类/标签/字数统计
export function useArticleStats() {
  const [stats, setStats] = useState<SiteStats | null>(null);
  const refresh = useCallback(async () => setStats(await getSiteStats()), []);
  useEffect(() => { refresh().catch(() => undefined); }, [refresh]);
  return { stats, refresh };
}
