import { useCallback, useEffect, useState } from 'react';
import { Button, Card, Form, Input, InputNumber, Modal, Select, Space, Table, Tabs, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { createArticle, deleteArticle, listAdminArticles, publishArticle, updateArticle } from '../api/article';
import { createCategory, deleteCategory, getCategoryTree, updateCategory } from '../api/category';
import { createTag, deleteTag, listTags, updateTag } from '../api/tag';
import { listAdminComments, updateCommentStatus } from '../api/comment';
import type { Article, Category, Comment, Tag as TagType } from '../types';
import ArticleStatusBadge from '../components/common/ArticleStatusBadge';
import MarkdownEditor from '../components/common/MarkdownEditor';
import EmptyState from '../components/common/EmptyState';
import { ArticleStatusLabels, ArticleStatus } from '../constants/article';
import { CommentStatusLabels } from '../constants/comment';
import { formatDateTime } from '../utils/dateFormat';
import { slugify } from '../utils/slug';

export default function Admin() {
  const [tab, setTab] = useState('articles');
  const [articles, setArticles] = useState<Article[]>([]);
  const [cats, setCats] = useState<Category[]>([]);
  const [tags, setTags] = useState<TagType[]>([]);
  const [comments, setComments] = useState<Comment[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [editorOpen, setEditorOpen] = useState(false);
  const [editing, setEditing] = useState<Article | null>(null);
  const [form] = Form.useForm();
  const [catForm] = Form.useForm();
  const [tagForm] = Form.useForm();

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const [a, c, t, cm] = await Promise.all([
        listAdminArticles({ page: 1, page_size: 50 }),
        getCategoryTree(),
        listTags(),
        listAdminComments({ page: 1, page_size: 50 }),
      ]);
      setArticles(a.list);
      setTotal(a.total);
      setCats(c);
      setTags(t);
      setComments(cm.list);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { load(); }, [load]);

  function flatCats(cs: Category[], prefix = ''): { value: number; label: string }[] {
    return cs.flatMap((c) => [
      { value: c.id, label: `${prefix}${c.name}` },
      ...flatCats(c.children ?? [], `${prefix}${c.name} / `),
    ]);
  }

  async function onSave(values: any) {
    const payload = {
      title: values.title,
      slug: values.slug || slugify(values.title) || `article-${Date.now()}`,
      category_id: values.category_id || null,
      summary: values.summary,
      content_markdown: values.content_markdown || '',
      status: values.status || ArticleStatus.DRAFT,
      is_top: values.is_top,
      tag_ids: values.tag_ids || [],
    };
    if (editing) await updateArticle(editing.id, payload);
    else await createArticle(payload);
    message.success('保存成功');
    setEditorOpen(false);
    form.resetFields();
    load();
  }

  async function onPublish(a: Article) {
    await publishArticle(a.id);
    message.success('已发布');
    load();
  }

  async function onDeleteArticle(a: Article) {
    await deleteArticle(a.id);
    message.success('已删除');
    load();
  }

  async function onCatSave(values: any) {
    await createCategory(values);
    message.success('分类已创建');
    catForm.resetFields();
    load();
  }

  async function onTagSave(values: any) {
    await createTag(values);
    message.success('标签已创建');
    tagForm.resetFields();
    load();
  }

  async function onCommentStatus(c: Comment, status: string) {
    await updateCommentStatus(c.id, status);
    message.success('评论状态已更新');
    load();
  }

  const articleColumns: ColumnsType<Article> = [
    { title: '标题', dataIndex: 'title' },
    { title: '别名', dataIndex: 'slug' },
    { title: '状态', dataIndex: 'status', render: (v) => <ArticleStatusBadge status={v} /> },
    { title: '置顶', dataIndex: 'is_top', render: (v) => (v ? <Tag color="gold">置顶</Tag> : '-') },
    { title: '阅读', dataIndex: 'view_count' },
    { title: '更新', dataIndex: 'updated_at', render: (v) => formatDateTime(v) },
    { title: '操作', render: (_, r) => (
      <Space>
        {r.status !== ArticleStatus.PUBLISHED && <a onClick={() => onPublish(r)}>发布</a>}
        <a onClick={() => { setEditing(r); form.setFieldsValue({ ...r, tag_ids: (r.tags ?? []).map((t) => t.id), category_id: r.category_id }); setEditorOpen(true); }}>编辑</a>
        <a style={{ color: '#ff4d4f' }} onClick={() => Modal.confirm({ title: '确认删除？', onOk: () => onDeleteArticle(r) })}>删除</a>
      </Space>
    ) },
  ];

  const commentColumns: ColumnsType<Comment> = [
    { title: '文章', render: (_, r) => r.article?.title ?? r.article_id },
    { title: '昵称', render: (_, r) => r.user?.nickname || r.user?.username || r.nickname },
    { title: '内容', dataIndex: 'content' },
    { title: '状态', dataIndex: 'status', render: (v) => CommentStatusLabels[v] ?? v },
    { title: '时间', dataIndex: 'created_at', render: (v) => formatDateTime(v) },
    { title: '操作', render: (_, r) => (
      <Space>
        {r.status !== 'approved' && <a onClick={() => onCommentStatus(r, 'approved')}>通过</a>}
        {r.status !== 'deleted' && <a style={{ color: '#ff4d4f' }} onClick={() => onCommentStatus(r, 'deleted')}>删除</a>}
      </Space>
    ) },
  ];

  return (
    <Card size="small">
      <Tabs activeKey={tab} onChange={setTab} items={[
        { key: 'articles', label: '文章管理', children: (
          <Space direction="vertical" size={12} style={{ width: '100%' }}>
            <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); form.resetFields(); setEditorOpen(true); }}>写文章</Button>
            <Table rowKey="id" size="small" columns={articleColumns} dataSource={articles} loading={loading} pagination={false} locale={{ emptyText: <EmptyState /> }} />
          </Space>
        ) },
        { key: 'categories', label: '分类管理', children: (
          <Space direction="vertical" size={12} style={{ width: '100%' }}>
            <Space>
              <Form form={catForm} layout="inline" onFinish={onCatSave}>
                <Form.Item name="name" rules={[{ required: true }]}><Input placeholder="分类名" /></Form.Item>
                <Form.Item name="slug" rules={[{ required: true }]}><Input placeholder="slug" /></Form.Item>
                <Form.Item name="parent_id"><Select allowClear placeholder="父分类" style={{ width: 140 }} options={flatCats(cats)} /></Form.Item>
                <Button type="primary" htmlType="submit">新增分类</Button>
              </Form>
            </Space>
            <Table size="small" rowKey="id" pagination={false} dataSource={cats} locale={{ emptyText: <EmptyState /> }}
              columns={[
                { title: '名称', dataIndex: 'name' },
                { title: 'slug', dataIndex: 'slug' },
                { title: '子分类', render: (_, r) => (r.children?.length ?? 0) },
                { title: '操作', render: (_, r) => <a style={{ color: '#ff4d4f' }} onClick={() => Modal.confirm({ title: '确认删除？', onOk: async () => { await deleteCategory(r.id); message.success('已删除'); load(); } })}>删除</a> },
              ]} />
          </Space>
        ) },
        { key: 'tags', label: '标签管理', children: (
          <Space direction="vertical" size={12} style={{ width: '100%' }}>
            <Form form={tagForm} layout="inline" onFinish={onTagSave}>
              <Form.Item name="name" rules={[{ required: true }]}><Input placeholder="标签名" /></Form.Item>
              <Form.Item name="slug" rules={[{ required: true }]}><Input placeholder="slug" /></Form.Item>
              <Button type="primary" htmlType="submit">新增标签</Button>
            </Form>
            <Space wrap>
              {tags.map((t) => (
                <Tag key={t.id} closable onClose={() => Modal.confirm({ title: '确认删除？', onOk: async () => { await deleteTag(t.id); message.success('已删除'); load(); } })}>{t.name}</Tag>
              ))}
            </Space>
          </Space>
        ) },
        { key: 'comments', label: '评论审核', children: (
          <Table rowKey="id" size="small" columns={commentColumns} dataSource={comments} pagination={false} locale={{ emptyText: <EmptyState /> }} />
        ) },
      ]} />

      <Modal open={editorOpen} title={editing ? '编辑文章' : '写文章'} width={820} onCancel={() => setEditorOpen(false)} footer={null} destroyOnClose>
        <Form form={form} layout="vertical" onFinish={onSave}>
          <Form.Item name="title" label="标题" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="slug" label="别名（slug）"><Input placeholder="留空自动生成" /></Form.Item>
          <Space size={12} style={{ display: 'flex' }}>
            <Form.Item name="category_id" label="分类"><Select allowClear style={{ width: 180 }} options={flatCats(cats)} /></Form.Item>
            <Form.Item name="tag_ids" label="标签"><Select mode="multiple" style={{ width: 260 }} options={tags.map((t) => ({ value: t.id, label: t.name }))} /></Form.Item>
            <Form.Item name="status" label="状态" initialValue={ArticleStatus.DRAFT}>
              <Select style={{ width: 140 }} options={Object.entries(ArticleStatusLabels).map(([value, label]) => ({ value, label }))} />
            </Form.Item>
          </Space>
          <Form.Item name="summary" label="摘要"><Input.TextArea rows={2} maxLength={500} /></Form.Item>
          <Form.Item name="content_markdown" label="正文（Markdown）"><MarkdownEditor /></Form.Item>
          <Button type="primary" htmlType="submit" block>保存</Button>
        </Form>
      </Modal>
    </Card>
  );
}
