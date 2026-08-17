import { useEffect, useState } from 'react';
import { Card, Col, Pagination, Row } from 'antd';
import { useParams } from 'react-router-dom';
import { getCategoryTree } from '../api/category';
import { useArticleStore } from '../stores/articleStore';
import ArticleCard from '../components/common/ArticleCard';
import CategoryTree from '../components/common/CategoryTree';
import EmptyState from '../components/common/EmptyState';
import type { Category } from '../types';

function collectIds(cats: Category[], slug: string, out: number[]): boolean {
  for (const c of cats) {
    if (c.slug === slug) {
      out.push(c.id);
      collectChildren(c.children ?? [], out);
      return true;
    }
    if (collectIds(c.children ?? [], slug, out)) return true;
  }
  return false;
}

function collectChildren(cats: Category[], out: number[]) {
  for (const c of cats) {
    out.push(c.id);
    collectChildren(c.children ?? [], out);
  }
}

export default function CategoryArticles() {
  const { slug } = useParams();
  const { articles, total, loading, fetchList } = useArticleStore();
  const [page, setPage] = useState(1);
  const [catIds, setCatIds] = useState<number[]>([]);

  useEffect(() => {
    getCategoryTree().then((tree) => {
      const ids: number[] = [];
      collectIds(tree, slug ?? '', ids);
      setCatIds(ids);
    }).catch(() => undefined);
  }, [slug]);

  useEffect(() => {
    if (!catIds.length) return;
    fetchList({ category_id: catIds[0], page, page_size: 8 }).catch(() => undefined);
  }, [fetchList, page, catIds]);

  return (
    <Row gutter={16}>
      <Col span={17}>
        <Card size="small" title={`分类：${slug}`} loading={loading}>
          {articles.length ? articles.map((a) => <ArticleCard key={a.id} article={a} />) : <EmptyState description="该分类暂无文章" />}
          <Pagination current={page} pageSize={8} total={total} onChange={setPage} style={{ textAlign: 'right' }} />
        </Card>
      </Col>
      <Col span={7}>
        <Card size="small" title="分类导航"><CategoryTree /></Card>
      </Col>
    </Row>
  );
}
