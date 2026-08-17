import { useEffect, useState } from 'react';
import { Card, Col, Row, Space, Spin, Tag, Typography } from 'antd';
import { EyeOutlined } from '@ant-design/icons';
import { useParams } from 'react-router-dom';
import { getArticle } from '../api/article';
import type { Article } from '../types';
import ArticleStatusBadge from '../components/common/ArticleStatusBadge';
import CommentList from '../components/common/CommentList';
import CategoryTree from '../components/common/CategoryTree';
import { formatDateTime } from '../utils/dateFormat';

const { Title } = Typography;

export default function ArticleDetail() {
  const { slug } = useParams();
  const [article, setArticle] = useState<Article | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!slug) return;
    setLoading(true);
    getArticle(slug)
      .then(setArticle)
      .catch(() => undefined)
      .finally(() => setLoading(false));
  }, [slug]);

  return (
    <Spin spinning={loading}>
      {article && (
        <Row gutter={16}>
          <Col span={17}>
            <Card size="small">
              <Space direction="vertical" size={12} style={{ width: '100%' }}>
                <Title level={2} style={{ marginBottom: 0 }}>{article.title}</Title>
                <Space wrap>
                  <ArticleStatusBadge status={article.status} />
                  <span>{formatDateTime(article.published_at)}</span>
                  <span><EyeOutlined /> {article.view_count}</span>
                  {article.category && <Tag color="blue">{article.category.name}</Tag>}
                  {(article.tags ?? []).map((t) => <Tag key={t.id}>{t.name}</Tag>)}
                </Space>
                <div className="markdown-body" dangerouslySetInnerHTML={{ __html: article.content_html }} />
              </Space>
            </Card>
            <Card size="small" title={`评论（${article.id}）`} style={{ marginTop: 16 }}>
              <CommentList articleId={article.id} />
            </Card>
          </Col>
          <Col span={7}>
            <Card size="small" title="分类导航"><CategoryTree /></Card>
          </Col>
        </Row>
      )}
    </Spin>
  );
}
