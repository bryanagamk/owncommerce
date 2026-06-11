import { AuthLayout } from '@owncommerce/ui';
import { Button, Form, Input, Typography, message } from 'antd';
import { Link, useNavigate } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { saveCustomerAuth } from '../lib/auth';

export default function AccountRegisterPage() {
  const navigate = useNavigate();
  const [form] = Form.useForm();

  const onFinish = async (values: { email: string; password: string; name: string; phone?: string }) => {
    try {
      const res = await storefrontApi.register(values);
      saveCustomerAuth(res.tokens.access_token, res.tokens.refresh_token, res.customer.name);
      message.success('Registrasi berhasil');
      navigate('/');
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Registrasi gagal');
    }
  };

  return (
    <AuthLayout title="Daftar" subtitle="Buat akun pelanggan">
      <Form form={form} layout="vertical" onFinish={onFinish} size="large">
        <Form.Item name="name" label="Nama" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="email" label="Email" rules={[{ required: true, type: 'email' }]}>
          <Input />
        </Form.Item>
        <Form.Item name="phone" label="Telepon">
          <Input />
        </Form.Item>
        <Form.Item name="password" label="Password" rules={[{ required: true, min: 8 }]}>
          <Input.Password />
        </Form.Item>
        <Button type="primary" htmlType="submit" block>
          Daftar
        </Button>
      </Form>
      <Typography.Paragraph style={{ textAlign: 'center', marginTop: 16 }}>
        Sudah punya akun? <Link to="/account/login">Masuk</Link>
      </Typography.Paragraph>
    </AuthLayout>
  );
}
