import { AuthLayout } from '@owncommerce/ui';
import { Button, Form, Input, Typography, message } from 'antd';
import { Link, useNavigate } from 'react-router-dom';
import { merchantApi } from '../lib/api';
import { completeOnboarding, saveAuth } from '../lib/auth';

export default function LoginPage() {
  const navigate = useNavigate();
  const [form] = Form.useForm();

  const onFinish = async (values: { email: string; password: string }) => {
    try {
      const result = await merchantApi.login(values);
      saveAuth(result.tokens.access_token, result.tokens.refresh_token, result.user.name);
      completeOnboarding();
      message.success('Login berhasil');
      navigate('/');
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Login gagal');
    }
  };

  return (
    <AuthLayout title="OwnCommerce Seller" subtitle="Masuk ke dashboard merchant">
      <Form form={form} layout="vertical" onFinish={onFinish} size="large">
        <Form.Item name="email" label="Email" rules={[{ required: true, type: 'email' }]}>
          <Input placeholder="owner@toko.com" />
        </Form.Item>
        <Form.Item name="password" label="Password" rules={[{ required: true }]}>
          <Input.Password placeholder="Password" />
        </Form.Item>
        <Button type="primary" htmlType="submit" block>
          Masuk
        </Button>
      </Form>
      <Typography.Paragraph style={{ textAlign: 'center', marginTop: 16, marginBottom: 0 }}>
        Belum punya akun? <Link to="/register">Daftar</Link>
      </Typography.Paragraph>
    </AuthLayout>
  );
}
