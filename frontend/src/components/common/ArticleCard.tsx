import { Card, Space, Tag, Typography } from 'antd';
import { EyeOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import type { Article } from '../../types';
import { formatDate, formatWordCount } from '../../utils/dateFormat';

const { Text } = Typography;

interface Props {
  article: Article;
}

// 首页/分类/标签/搜索 共用文章卡片
export default function ArticleCard({ article }: Props) {
  const navigate = useNavigate();
  return (
    <Card
      hoverable
      onClick={() => navigate(`/articles/${article.slug}`)}
      title={article.title}
      extra={article.is_top ? <Tag color="gold">置顶</Tag> : undefined}
      style={{ marginBottom: 12 }}
    >
      <Space direction="vertical" size={6} style={{ width: '100%' }}>
        <Text type="secondary">{article.summary || '（无摘要）'}</Text>
        <Space wrap size={4}>
          {article.category && <Tag color="blue">{article.category.name}</Tag>}
          {(article.tags ?? []).map((t) => <Tag key={t.id}>{t.name}</Tag>)}
        </Space>
        <Space size={16}>
          <span>{formatDate(article.published_at)}</span>
          <span><EyeOutlined /> {article.view_count}</span>
          <span>{formatWordCount(article.word_count)}</span>
        </Space>
      </Space>
    </Card>
  );
}
