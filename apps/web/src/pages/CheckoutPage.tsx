import type { Address } from '@owncommerce/types';
import { Button, Card, Form, Input, InputNumber, Select, Typography, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { isCustomerLoggedIn } from '../lib/auth';

export default function CheckoutPage() {
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [addresses, setAddresses] = useState<Address[]>([]);

  useEffect(() => {
    if (!isCustomerLoggedIn()) return;
    Promise.all([storefrontApi.me(), storefrontApi.listAddresses()])
      .then(([me, addr]) => {
        setAddresses(addr);
        const defaultAddr = addr.find((a) => a.is_default) ?? addr[0];
        form.setFieldsValue({
          customer_email: me.email,
          recipient_name: defaultAddr?.recipient_name,
          recipient_phone: defaultAddr?.phone,
          shipping_address: defaultAddr?.address_line,
          shipping_city: defaultAddr?.city,
          shipping_province: defaultAddr?.province,
          shipping_postal_code: defaultAddr?.postal_code,
          shipping_cost: 15000,
        });
      })
      .catch(() => {
        form.setFieldsValue({ shipping_cost: 15000 });
      });
  }, [form]);

  const applyAddress = (id: string) => {
    const addr = addresses.find((a) => a.id === id);
    if (!addr) return;
    form.setFieldsValue({
      recipient_name: addr.recipient_name,
      recipient_phone: addr.phone,
      shipping_address: addr.address_line,
      shipping_city: addr.city,
      shipping_province: addr.province,
      shipping_postal_code: addr.postal_code,
    });
  };

  const onFinish = async (values: Record<string, string | number>) => {
    try {
      const order = await storefrontApi.checkout(values);
      message.success('Pesanan dibuat');
      navigate(`/payment/${order.id}`);
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Checkout gagal');
    }
  };

  return (
    <div style={{ maxWidth: 560, margin: '0 auto', padding: '24px 16px' }}>
      <Typography.Title level={3}>Checkout</Typography.Title>
      <Card bordered style={{ borderColor: '#F0F0F0' }}>
        <Form form={form} layout="vertical" onFinish={onFinish} initialValues={{ shipping_cost: 15000 }}>
          {addresses.length > 0 && (
            <Form.Item label="Gunakan alamat tersimpan">
              <Select
                placeholder="Pilih alamat"
                allowClear
                onChange={applyAddress}
                options={addresses.map((a) => ({
                  value: a.id,
                  label: `${a.label || a.recipient_name} — ${a.city}`,
                }))}
              />
            </Form.Item>
          )}
          <Form.Item name="recipient_name" label="Nama Penerima" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="recipient_phone" label="Telepon" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="customer_email" label="Email" rules={[{ type: 'email' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="shipping_address" label="Alamat" rules={[{ required: true }]}>
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="shipping_city" label="Kota" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="shipping_province" label="Provinsi" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="shipping_postal_code" label="Kode Pos" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="shipping_cost" label="Ongkir (IDR)">
            <InputNumber style={{ width: '100%' }} min={0} />
          </Form.Item>
          <Form.Item name="customer_note" label="Catatan">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Button type="primary" htmlType="submit" block size="large">
            Buat Pesanan & Bayar
          </Button>
        </Form>
      </Card>
    </div>
  );
}
