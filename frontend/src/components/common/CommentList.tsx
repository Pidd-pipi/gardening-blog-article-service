import { useEffect, useState } from 'react';
import { Avatar, Button, Divider, Form, Input, List, message } from 'antd';
import { useAuth } from '../../hooks/useAuth';
import { createComment, listComments } from '../../api/comment';
import type { Comment as CommentType } from '../../types';
import { formatDateTime } from '../../utils/dateFormat';

interface Props {
  articleId: number;
}

function CommentItem({ c, depth }: { c: CommentType; depth: number }) {
  const author = c.user?.nickname || c.user?.username || c.nickname || '匿名';
  return (
    <div style={{ marginLeft: depth * 24, marginBottom: 12 }}>
      <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
        <Avatar size="small" src={c.user?.avatar}>{author.slice(0, 1)}</Avatar>
        <b>{author}</b>
        <span style={{ color: '#999', fontSize: 12 }}>{formatDateTime(c.created_at)}</span>
      </div>
      <div style={{ marginTop: 4 }}>{c.content}</div>
      {c.replies?.map((r) => <CommentItem key={r.id} c={r} depth={depth + 1} />)}
    </div>
  );
}

export default function CommentList({ articleId }: Props) {
  const { user, isLoggedIn } = useAuth();
  const [comments, setComments] = useState<CommentType[]>([]);
  const [form] = Form.useForm();

  const load = async () => setComments(await listComments(articleId));
  useEffect(() => { load().catch(() => undefined); }, [articleId]);

  async function onFinish(values: { nickname?: string; email?: string; content: string }) {
    await createComment({
      article_id: articleId,
      nickname: isLoggedIn ? user?.nickname || user?.username : values.nickname,
      email: isLoggedIn ? user?.email : values.email,
      content: values.content,
    });
    message.success('评论已提交');
    form.resetFields();
    load();
  }

  return (
    <div>
      <List
        dataSource={comments}
        renderItem={(c) => <CommentItem c={c} depth={0} />}
        locale={{ emptyText: '暂无评论' }}
        split={false}
      />
      <Divider />
      <Form form={form} layout="vertical" onFinish={onFinish}>
        {!isLoggedIn && (
          <div style={{ display: 'flex', gap: 12 }}>
            <Form.Item name="nickname" label="昵称" style={{ flex: 1 }}><Input maxLength={50} /></Form.Item>
            <Form.Item name="email" label="邮箱" style={{ flex: 1 }}><Input maxLength={100} /></Form.Item>
          </div>
        )}
        <Form.Item name="content" rules={[{ required: true, message: '请输入评论内容' }, { max: 1000 }]}>
          <Input.TextArea rows={3} placeholder={isLoggedIn ? '登录用户评论将直接通过审核' : '评论提交后需博主审核'} />
        </Form.Item>
        <Button type="primary" htmlType="submit">提交评论</Button>
      </Form>
    </div>
  );
}
