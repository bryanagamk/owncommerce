import { Button, Card, Form, Input, InputNumber, Typography, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { storefrontApi } from '../lib/api';

export default function CheckoutPage() {
  const navigate = useNavigate();
  const [form] = Form.useForm();

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
    <div style={{ maxWidth: 560, margin: '0 auto', padding: '32px 16px' }}>
      <Typography.Title level={3}>Checkout</Typography.Title>
      <Card bordered style={{ borderColor: '#F0F0F0' }}>
        <Form form={form} layout="vertical" onFinish={onFinish} initialValues={{ shipping_cost: 15000 }}>
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
