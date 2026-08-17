import { useEffect, useState } from 'react';
import { Card, Col, Pagination, Row, Tag } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { listTags } from '../api/tag';
import { useArticleStore } from '../stores/articleStore';
import ArticleCard from '../components/common/ArticleCard';
import EmptyState from '../components/common/EmptyState';

export default function TagArticles() {
  const { slug } = useParams();
  const { articles, total, loading, fetchList } = useArticleStore();
  const [page, setPage] = useState(1);
  const [tagId, setTagId] = useState<number | undefined>();

  useEffect(() => {
    listTags().then((tags) => {
      const t = tags.find((x) => x.slug === slug);
      setTagId(t?.id);
    }).catch(() => undefined);
  }, [slug]);

  useEffect(() => {
    if (!tagId) return;
    fetchList({ tag_id: tagId, page, page_size: 8 }).catch(() => undefined);
  }, [fetchList, page, tagId]);

  return (
    <Row gutter={16}>
      <Col span={17}>
        <Card size="small" title={`标签：${slug}`} loading={loading}>
          {articles.length ? articles.map((a) => <ArticleCard key={a.id} article={a} />) : <EmptyState description="该标签暂无文章" />}
          <Pagination current={page} pageSize={8} total={total} onChange={setPage} style={{ textAlign: 'right' }} />
        </Card>
      </Col>
      <Col span={7}>
        <Card size="small" title="标签云"><TagsCloud /></Card>
      </Col>
    </Row>
  );
}

function TagsCloud() {
  const [tags, setTags] = useState<{ id: number; name: string; slug: string }[]>([]);
  const navigate = useNavigate();
  useEffect(() => { listTags().then(setTags).catch(() => undefined); }, []);
  return (
    <div>
      {tags.map((t) => (
        <Tag key={t.id} color="blue" style={{ cursor: 'pointer' }} onClick={() => navigate(`/tags/${t.slug}`)}>
          {t.name}
        </Tag>
      ))}
    </div>
  );
}
