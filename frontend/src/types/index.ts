export interface User {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar: string;
  bio: string;
  social_links: string;
  role: string;
}

export interface Category {
  id: number;
  parent_id?: number | null;
  name: string;
  slug: string;
  description: string;
  sort_order: number;
  children?: Category[];
}

export interface Tag {
  id: number;
  name: string;
  slug: string;
}

export interface Article {
  id: number;
  user_id: number;
  category_id?: number | null;
  title: string;
  slug: string;
  summary: string;
  content_markdown: string;
  content_html: string;
  cover: string;
  status: string;
  is_top: boolean;
  published_at?: string | null;
  view_count: number;
  word_count: number;
  created_at: string;
  updated_at: string;
  user?: User;
  category?: Category | null;
  tags?: Tag[];
}

export interface Comment {
  id: number;
  article_id: number;
  parent_id?: number | null;
  user_id?: number | null;
  nickname: string;
  email: string;
  content: string;
  status: string;
  created_at: string;
  user?: User | null;
  article?: Article | null;
  replies?: Comment[];
}

export interface SiteStats {
  article_count: number;
  category_count: number;
  tag_count: number;
  comment_count: number;
  word_count: number;
}

export interface ArchiveMonth {
  month: string;
  count: number;
}

export interface PageData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}
