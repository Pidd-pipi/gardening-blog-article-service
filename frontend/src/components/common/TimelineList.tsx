import { Empty, Timeline } from 'antd';
import { useNavigate } from 'react-router-dom';
import type { Article } from '../../types';
import { formatDate } from '../../utils/dateFormat';

interface Props {
  articles: Article[];
  emptyText?: string;
}

// 归档时间线
export default function TimelineList({ articles, emptyText = '暂无文章' }: Props) {
  const navigate = useNavigate();
  if (!articles.length) return <Empty description={emptyText} />;
  return (
    <Timeline
      items={articles.map((a) => ({
        color: a.is_top ? 'gold' : 'blue',
        children: (
          <div>
            <a onClick={() => navigate(`/articles/${a.slug}`)}>{a.title}</a>
            <span style={{ marginLeft: 12, color: '#999' }}>{formatDate(a.published_at)} · {a.view_count} 阅读</span>
          </div>
        ),
      }))}
    />
  );
}
