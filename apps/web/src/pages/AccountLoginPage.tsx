import { AuthLayout } from '@owncommerce/ui';
import { Button, Form, Input, Typography, message } from 'antd';
import { Link, useNavigate } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { saveCustomerAuth } from '../lib/auth';

export default function AccountLoginPage() {
  const navigate = useNavigate();
  const [form] = Form.useForm();

  const onFinish = async (values: { email: string; password: string }) => {
    try {
      const res = await storefrontApi.login(values);
      saveCustomerAuth(res.tokens.access_token, res.tokens.refresh_token, res.customer.name);
      message.success('Login berhasil');
      navigate('/');
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Login gagal');
    }
  };

  return (
    <AuthLayout title="Masuk" subtitle="Akun pelanggan">
      <Form form={form} layout="vertical" onFinish={onFinish} size="large">
        <Form.Item name="email" label="Email" rules={[{ required: true, type: 'email' }]}>
          <Input />
        </Form.Item>
        <Form.Item name="password" label="Password" rules={[{ required: true }]}>
          <Input.Password />
        </Form.Item>
        <Button type="primary" htmlType="submit" block>
          Masuk
        </Button>
      </Form>
      <Typography.Paragraph style={{ textAlign: 'center', marginTop: 16 }}>
        Belum punya akun? <Link to="/account/register">Daftar</Link>
      </Typography.Paragraph>
    </AuthLayout>
  );
}
