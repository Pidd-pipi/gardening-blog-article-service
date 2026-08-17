import { useState } from 'react';
import { Button, Card, Form, Input, Tabs, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { login, register } from '../api/user';
import { useAuth } from '../hooks/useAuth';

export default function Login() {
  const [tab, setTab] = useState('login');
  const [loading, setLoading] = useState(false);
  const { setAuth } = useAuth();
  const navigate = useNavigate();

  async function onFinish(values: { email: string; password: string; username?: string; nickname?: string }) {
    setLoading(true);
    try {
      const result = tab === 'login'
        ? await login(values.email, values.password)
        : await register(values.username || 'user', values.email, values.password, values.nickname || '游客');
      setAuth(result.token, result.user);
      message.success('登录成功');
      navigate('/');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#f0f2f5' }}>
      <Card title="gbblog 个人博客登录" style={{ width: 400 }}>
        <Tabs activeKey={tab} onChange={setTab} items={[
          { key: 'login', label: '登录', children: (
            <Form layout="vertical" onFinish={onFinish}>
              <Form.Item name="email" label="邮箱" rules={[{ required: true, type: 'email' }]}><Input placeholder="admin@gbblog.com" /></Form.Item>
              <Form.Item name="password" label="密码" rules={[{ required: true }]}><Input.Password placeholder="admin123" /></Form.Item>
              <Button type="primary" htmlType="submit" block loading={loading}>登录</Button>
            </Form>
          ) },
          { key: 'register', label: '注册', children: (
            <Form layout="vertical" onFinish={onFinish}>
              <Form.Item name="username" label="用户名" rules={[{ required: true }]}><Input /></Form.Item>
              <Form.Item name="email" label="邮箱" rules={[{ required: true, type: 'email' }]}><Input /></Form.Item>
              <Form.Item name="nickname" label="昵称"><Input /></Form.Item>
              <Form.Item name="password" label="密码" rules={[{ required: true, min: 6 }]}><Input.Password /></Form.Item>
              <Button type="primary" htmlType="submit" block loading={loading}>注册并登录</Button>
            </Form>
          ) },
        ]} />
      </Card>
    </div>
  );
}
