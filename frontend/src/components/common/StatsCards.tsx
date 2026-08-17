import { Card, Col, Row, Statistic } from 'antd';
import { FileTextOutlined, FolderOutlined, TagsOutlined, CommentOutlined, ReadOutlined } from '@ant-design/icons';
import { useArticleStats } from '../../hooks/useArticleStats';
import { formatWordCount } from '../../utils/dateFormat';

// 首页站点统计卡片
export default function StatsCards() {
  const { stats } = useArticleStats();
  return (
    <Row gutter={12}>
      <Col span={5}><Card size="small"><Statistic title="文章" value={stats?.article_count ?? 0} prefix={<FileTextOutlined />} /></Card></Col>
      <Col span={5}><Card size="small"><Statistic title="分类" value={stats?.category_count ?? 0} prefix={<FolderOutlined />} /></Card></Col>
      <Col span={5}><Card size="small"><Statistic title="标签" value={stats?.tag_count ?? 0} prefix={<TagsOutlined />} /></Card></Col>
      <Col span={5}><Card size="small"><Statistic title="评论" value={stats?.comment_count ?? 0} prefix={<CommentOutlined />} /></Card></Col>
      <Col span={4}><Card size="small"><Statistic title="总字数" value={formatWordCount(stats?.word_count ?? 0)} prefix={<ReadOutlined />} /></Card></Col>
    </Row>
  );
}
