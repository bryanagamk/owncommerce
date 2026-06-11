import { PageHeader } from '@owncommerce/ui';
import type { Tenant } from '@owncommerce/types';
import { Button, Card, Form, Input, message } from 'antd';
import { useEffect, useState } from 'react';
import { merchantApi } from '../lib/api';

export default function SettingsPage() {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    merchantApi.getTenant().then((res) => {
      form.setFieldsValue(res.tenant);
      setLoading(false);
    });
  }, [form]);

  const onFinish = async (values: Partial<Tenant>) => {
    try {
      await merchantApi.updateStore(values as Record<string, string>);
      message.success('Pengaturan toko disimpan');
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal menyimpan');
    }
  };

  return (
    <div>
      <PageHeader title="Pengaturan Toko" subtitle="Informasi dan kontak toko Anda" />
      <Card bordered style={{ borderColor: '#F0F0F0', maxWidth: 640 }}>
        <Form form={form} layout="vertical" onFinish={onFinish} disabled={loading}>
          <Form.Item name="name" label="Nama Toko" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="Deskripsi">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item name="contact_email" label="Email Kontak">
            <Input type="email" />
          </Form.Item>
          <Form.Item name="contact_phone" label="Telepon">
            <Input />
          </Form.Item>
          <Form.Item name="address" label="Alamat">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="city" label="Kota">
            <Input />
          </Form.Item>
          <Form.Item name="province" label="Provinsi">
            <Input />
          </Form.Item>
          <Form.Item name="postal_code" label="Kode Pos">
            <Input />
          </Form.Item>
          <Button type="primary" htmlType="submit">
            Simpan
          </Button>
        </Form>
      </Card>
    </div>
  );
}
