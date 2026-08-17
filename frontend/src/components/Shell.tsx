import { Layout, Menu, Button, Space, Avatar, Dropdown, Input } from 'antd';
import { HomeOutlined, FolderOutlined, TagsOutlined, ClockCircleOutlined, SearchOutlined, UserOutlined, LoginOutlined, EditOutlined } from '@ant-design/icons';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { useAuth } from '../hooks/useAuth';

const { Header, Content, Sider } = Layout;

export default function Shell() {
  const { user, isAdmin, logout } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const [keyword, setKeyword] = useState('');

  const menuItems = [
    { key: '/', icon: <HomeOutlined />, label: '首页' },
    { key: '/categories/tech', icon: <FolderOutlined />, label: '技术分类' },
    { key: '/archives', icon: <ClockCircleOutlined />, label: '归档' },
    { key: '/search', icon: <SearchOutlined />, label: '搜索' },
  ];
  if (isAdmin) {
    menuItems.push({ key: '/admin', icon: <EditOutlined />, label: '后台管理' });
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider breakpoint="lg" collapsedWidth={64}>
        <div style={{ color: '#fff', padding: 16, fontWeight: 700, fontSize: 17, cursor: 'pointer' }} onClick={() => navigate('/')}>gbblog</div>
        <Menu theme="dark" mode="inline" selectedKeys={[location.pathname]} items={menuItems} onClick={({ key }) => navigate(key)} />
        <div style={{ padding: 16 }}>
          <Input.Search
            placeholder="搜索文章"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onSearch={(v) => navigate(`/search?q=${encodeURIComponent(v)}`)}
          />
        </div>
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div style={{ fontWeight: 600 }}>功能完善的个人博客系统</div>
          <Space>
            {isAdmin && <Button size="small" type="primary" onClick={() => navigate('/admin')}>写文章</Button>}
            {user ? (
              <Dropdown menu={{ items: [{ key: 'logout', label: '退出登录' }], onClick: ({ key }) => { if (key === 'logout') { logout(); navigate('/'); } } }}>
                <Avatar size="small" src={user.avatar || undefined} icon={<UserOutlined />} style={{ cursor: 'pointer' }}>
                  {!user.avatar && (user.nickname || user.username || 'U').slice(0, 1)}
                </Avatar>
              </Dropdown>
            ) : (
              <Button size="small" icon={<LoginOutlined />} onClick={() => navigate('/login')}>登录</Button>
            )}
          </Space>
        </Header>
        <Content style={{ margin: 24 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
