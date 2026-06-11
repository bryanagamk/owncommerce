import type { Order } from '@owncommerce/types';
import { Table, Tag, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { formatPrice } from '../lib/format';

export default function AccountOrdersPage() {
  const navigate = useNavigate();
  const [orders, setOrders] = useState<Order[]>([]);

  useEffect(() => {
    storefrontApi.listOrders().then((res) => setOrders(res.data));
  }, []);

  return (
    <div style={{ maxWidth: 900, margin: '0 auto', padding: '32px 16px' }}>
      <Typography.Title level={3}>Pesanan Saya</Typography.Title>
      <Table
        rowKey="id"
        dataSource={orders}
        bordered
        onRow={(r) => ({ onClick: () => navigate(`/payment/${r.id}`), style: { cursor: 'pointer' } })}
        columns={[
          { title: 'No. Order', dataIndex: 'order_number' },
          {
            title: 'Status',
            dataIndex: 'status',
            render: (s: string) => <Tag>{s}</Tag>,
          },
          {
            title: 'Total',
            dataIndex: 'total',
            render: (v: number) => formatPrice(v),
          },
          {
            title: 'Tanggal',
            dataIndex: 'created_at',
            render: (v: string) => new Date(v).toLocaleDateString('id-ID'),
          },
        ]}
      />
    </div>
  );
}
