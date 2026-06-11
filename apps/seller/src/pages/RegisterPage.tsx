import { AuthLayout } from '@owncommerce/ui';
import { Button, Form, Input, Typography, message } from 'antd';
import { Link, useNavigate } from 'react-router-dom';
import { merchantApi } from '../lib/api';
import { saveAuth } from '../lib/auth';

export default function RegisterPage() {
  const navigate = useNavigate();
  const [form] = Form.useForm();

  const onFinish = async (values: {
    email: string;
    password: string;
    name: string;
    store_name: string;
    slug: string;
  }) => {
    try {
      const result = await merchantApi.register(values);
      saveAuth(result.tokens.access_token, result.tokens.refresh_token, result.user.name);
      message.success('Registrasi berhasil');
      navigate('/');
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Registrasi gagal');
    }
  };

  return (
    <AuthLayout title="Daftar Merchant" subtitle="Buat toko online Anda">
      <Form form={form} layout="vertical" onFinish={onFinish} size="large">
        <Form.Item name="name" label="Nama Anda" rules={[{ required: true }]}>
          <Input placeholder="Budi" />
        </Form.Item>
        <Form.Item name="store_name" label="Nama Toko" rules={[{ required: true }]}>
          <Input placeholder="Toko Bunga" />
        </Form.Item>
        <Form.Item
          name="slug"
          label="Subdomain"
          rules={[{ required: true }, { pattern: /^[a-z0-9-]+$/, message: 'Huruf kecil, angka, strip' }]}
          extra="Contoh: tokobunga → tokobunga.localhost"
        >
          <Input placeholder="tokobunga" />
        </Form.Item>
        <Form.Item name="email" label="Email" rules={[{ required: true, type: 'email' }]}>
          <Input placeholder="owner@toko.com" />
        </Form.Item>
        <Form.Item name="password" label="Password" rules={[{ required: true, min: 8 }]}>
          <Input.Password placeholder="Min. 8 karakter" />
        </Form.Item>
        <Button type="primary" htmlType="submit" block>
          Daftar
        </Button>
      </Form>
      <Typography.Paragraph style={{ textAlign: 'center', marginTop: 16, marginBottom: 0 }}>
        Sudah punya akun? <Link to="/login">Masuk</Link>
      </Typography.Paragraph>
    </AuthLayout>
  );
}
