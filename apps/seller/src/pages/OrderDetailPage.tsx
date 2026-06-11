import { PageHeader } from '@owncommerce/ui';
import type { Order } from '@owncommerce/types';
import { Button, Card, Descriptions, Select, Table, Tag, message } from 'antd';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { merchantApi } from '../lib/api';
import { formatPrice } from '../lib/format';

export default function OrderDetailPage() {
  const { id } = useParams();
  const [order, setOrder] = useState<Order | null>(null);
  const [nextStatus, setNextStatus] = useState('processing');

  const load = () => {
    if (id) merchantApi.getOrder(id).then(setOrder);
  };

  useEffect(() => {
    load();
  }, [id]);

  const updateStatus = async () => {
    if (!id) return;
    try {
      await merchantApi.updateOrderStatus(id, { status: nextStatus });
      message.success('Status diperbarui');
      load();
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal update');
    }
  };

  if (!order) return null;

  return (
    <div>
      <PageHeader title={`Order ${order.order_number}`} subtitle={order.status} />
      <Card bordered style={{ borderColor: '#F0F0F0', marginBottom: 24 }}>
        <Descriptions column={2}>
          <Descriptions.Item label="Status">
            <Tag>{order.status}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="Pembayaran">{order.payment_status}</Descriptions.Item>
          <Descriptions.Item label="Penerima">{order.recipient_name}</Descriptions.Item>
          <Descriptions.Item label="Telepon">{order.recipient_phone}</Descriptions.Item>
          <Descriptions.Item label="Alamat" span={2}>
            {order.shipping_address}, {order.shipping_city}, {order.shipping_province} {order.shipping_postal_code}
          </Descriptions.Item>
          <Descriptions.Item label="Subtotal">{formatPrice(order.subtotal)}</Descriptions.Item>
          <Descriptions.Item label="Ongkir">{formatPrice(order.shipping_cost)}</Descriptions.Item>
          <Descriptions.Item label="Total">
            <strong>{formatPrice(order.total)}</strong>
          </Descriptions.Item>
        </Descriptions>
      </Card>
      <Card title="Item Pesanan" bordered style={{ borderColor: '#F0F0F0', marginBottom: 24 }}>
        <Table
          rowKey="id"
          dataSource={order.items ?? []}
          pagination={false}
          columns={[
            { title: 'Produk', dataIndex: 'product_name' },
            { title: 'Qty', dataIndex: 'quantity' },
            {
              title: 'Harga',
              dataIndex: 'unit_price',
              render: (v: number) => formatPrice(v),
            },
            {
              title: 'Subtotal',
              dataIndex: 'subtotal',
              render: (v: number) => formatPrice(v),
            },
          ]}
        />
      </Card>
      {order.status !== 'cancelled' && order.status !== 'completed' && (
        <Card title="Update Status" bordered style={{ borderColor: '#F0F0F0', maxWidth: 400 }}>
          <Select
            style={{ width: '100%', marginBottom: 12 }}
            value={nextStatus}
            onChange={setNextStatus}
            options={['processing', 'shipped', 'completed', 'cancelled'].map((s) => ({
              value: s,
              label: s,
            }))}
          />
          <Button type="primary" onClick={updateStatus}>
            Update Status
          </Button>
        </Card>
      )}
    </div>
  );
}
