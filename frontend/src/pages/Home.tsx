import { useEffect, useState } from 'react';
import { Card, Col, Pagination, Row, Space, Spin } from 'antd';
import { useArticleStore } from '../stores/articleStore';
import ArticleCard from '../components/common/ArticleCard';
import CategoryTree from '../components/common/CategoryTree';
import StatsCards from '../components/common/StatsCards';
import { topArticles } from '../api/article';
import type { Article } from '../types';
import EmptyState from '../components/common/EmptyState';

export default function Home() {
  const { articles, total, loading, fetchList } = useArticleStore();
  const [page, setPage] = useState(1);
  const [top, setTop] = useState<Article[]>([]);

  useEffect(() => {
    fetchList({ page, page_size: 5 }).catch(() => undefined);
    topArticles().then(setTop).catch(() => undefined);
  }, [fetchList, page]);

  return (
    <Spin spinning={loading}>
      <Space direction="vertical" size={16} style={{ width: '100%' }}>
        <StatsCards />
        <Row gutter={16}>
          <Col span={17}>
            <Card size="small" title="最新文章">
              {articles.length ? articles.map((a) => <ArticleCard key={a.id} article={a} />) : <EmptyState />}
              <Pagination current={page} pageSize={5} total={total} onChange={setPage} style={{ textAlign: 'right' }} />
            </Card>
          </Col>
          <Col span={7}>
            <Card size="small" title="分类导航" style={{ marginBottom: 16 }}><CategoryTree /></Card>
            <Card size="small" title="热门文章">
              {top.length ? top.map((a) => (
                <div key={a.id} style={{ padding: '6px 0' }}>
                  <a href={`#/articles/${a.slug}`} onClick={(e) => { e.preventDefault(); location.href = `/articles/${a.slug}`; }}>{a.title}</a>
                  <span style={{ marginLeft: 8, color: '#999' }}>{a.view_count} 阅读</span>
                </div>
              )) : <EmptyState description="暂无热门文章" />}
            </Card>
          </Col>
        </Row>
      </Space>
    </Spin>
  );
}
