import { Card, Col, Row, Tag } from 'antd';
import { useArchive } from '../hooks/useArchive';
import TimelineList from '../components/common/TimelineList';
import EmptyState from '../components/common/EmptyState';

export default function Archives() {
  const { months, articles, current, loadMonth } = useArchive();

  return (
    <Row gutter={16}>
      <Col span={7}>
        <Card size="small" title="归档月份">
          {months.length ? months.map((m) => (
            <div key={m.month} style={{ padding: '6px 0' }}>
              <Tag color={current === m.month ? 'blue' : 'default'} style={{ cursor: 'pointer' }} onClick={() => loadMonth(m.month)}>
                {m.month}（{m.count}）
              </Tag>
            </div>
          )) : <EmptyState description="暂无归档" />}
        </Card>
      </Col>
      <Col span={17}>
        <Card size="small" title={current ? `${current} 归档文章` : '按年月时间线'}>
          <TimelineList articles={articles} emptyText={current ? '该月暂无文章' : '请选择归档月份'} />
        </Card>
      </Col>
    </Row>
  );
}
