import { Empty, Tree } from 'antd';
import type { DataNode } from 'antd/es/tree';
import { useNavigate } from 'react-router-dom';
import { useCategoryStore } from '../../stores/categoryStore';
import type { Category } from '../../types';

function toNodes(cats: Category[]): DataNode[] {
  return cats.map((c) => ({
    key: c.slug,
    title: `${c.name}${c.children?.length ? `（${c.children.length}）` : ''}`,
    children: c.children?.length ? toNodes(c.children) : undefined,
  }));
}

// 分类树导航（分类页/后台共用）
export default function CategoryTree() {
  const { tree, refresh } = useCategoryStore();
  const navigate = useNavigate();
  if (!tree.length) {
    refresh().catch(() => undefined);
  }
  if (!tree.length) return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无分类" />;
  return (
    <Tree
      blockNode
      defaultExpandAll
      treeData={toNodes(tree)}
      onSelect={(keys) => { if (keys.length) navigate(`/categories/${keys[0]}`); }}
    />
  );
}
