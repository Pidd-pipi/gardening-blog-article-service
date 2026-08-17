import { Tag } from 'antd';
import { ArticleStatus, ArticleStatusLabels } from '../../constants/article';

interface Props {
  status?: string;
}

const colorMap: Record<string, string> = {
  [ArticleStatus.DRAFT]: 'default',
  [ArticleStatus.PUBLISHED]: 'green',
  [ArticleStatus.SCHEDULED]: 'blue',
};

// 文章详情/后台 共用状态徽标
export default function ArticleStatusBadge({ status }: Props) {
  return <Tag color={colorMap[status ?? ''] ?? 'default'}>{ArticleStatusLabels[status ?? ''] ?? status ?? '-'}</Tag>;
}
