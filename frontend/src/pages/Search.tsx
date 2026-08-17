import { useEffect, useState } from 'react';
import { Card, Input, Pagination, Space } from 'antd';
import { useSearchParams } from 'react-router-dom';
import { searchArticles } from '../api/article';
import type { Article } from '../types';
import ArticleCard from '../components/common/ArticleCard';
import EmptyState from '../components/common/EmptyState';

export default function Search() {
  const [params, setParams] = useSearchParams();
  const [keyword, setKeyword] = useState(params.get('q') ?? '');
  const [articles, setArticles] = useState<Article[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);

  const doSearch = async (q: string, p: number) => {
    if (!q) { setArticles([]); setTotal(0); return; }
    const data = await searchArticles(q, { page: p, page_size: 8 });
    setArticles(data.list);
    setTotal(data.total);
  };

  useEffect(() => {
    const q = params.get('q') ?? '';
    setKeyword(q);
    doSearch(q, page).catch(() => undefined);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [params, page]);

  return (
    <Card size="small" title="全文搜索">
      <Space direction="vertical" size={16} style={{ width: '100%' }}>
        <Input.Search
          placeholder="输入关键词"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onSearch={(v) => { setPage(1); setParams({ q: v }); }}
          enterButton
        />
        {articles.length ? articles.map((a) => <ArticleCard key={a.id} article={a} />) : <EmptyState description={keyword ? '未找到相关文章' : '请输入关键词搜索'} />}
        {total > 0 && <Pagination current={page} pageSize={8} total={total} onChange={setPage} style={{ textAlign: 'right' }} />}
      </Space>
    </Card>
  );
}
