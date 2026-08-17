import { Input, Tabs } from 'antd';
import { renderMarkdown } from '../../utils/markdown';

interface Props {
  value?: string;
  onChange?: (value: string) => void;
}

// Markdown 编辑器（写作 + 预览）
export default function MarkdownEditor({ value = '', onChange }: Props) {
  return (
    <Tabs
      items={[
        {
          key: 'write',
          label: '编辑',
          children: <Input.TextArea rows={16} value={value} onChange={(e) => onChange?.(e.target.value)} placeholder="支持 Markdown 语法" />,
        },
        {
          key: 'preview',
          label: '预览',
          children: <div className="markdown-preview" dangerouslySetInnerHTML={{ __html: renderMarkdown(value) }} />,
        },
      ]}
    />
  );
}
