import { PageHeader } from '@owncommerce/ui';
import type { Order } from '@owncommerce/types';
import { Select, Table, Tag } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { merchantApi } from '../lib/api';
import { formatPrice } from '../lib/format';

const statusColors: Record<string, string> = {
  pending_payment: 'orange',
  paid: 'blue',
  processing: 'cyan',
  shipped: 'purple',
  completed: 'green',
  cancelled: 'red',
};

export default function OrdersPage() {
  const navigate = useNavigate();
  const [items, setItems] = useState<Order[]>([]);
  const [status, setStatus] = useState('');

  useEffect(() => {
    merchantApi.listOrders({ status: status || undefined, limit: 50 }).then((res) => setItems(res.data));
  }, [status]);

  return (
    <div>
      <PageHeader title="Pesanan" subtitle="Kelola pesanan masuk" />
      <Select
        placeholder="Filter status"
        allowClear
        style={{ width: 200, marginBottom: 16 }}
        onChange={(v) => setStatus(v ?? '')}
        options={[
          'pending_payment',
          'paid',
          'processing',
          'shipped',
          'completed',
          'cancelled',
        ].map((s) => ({ value: s, label: s }))}
      />
      <Table
        rowKey="id"
        dataSource={items}
        bordered
        onRow={(r) => ({ onClick: () => navigate(`/orders/${r.id}`), style: { cursor: 'pointer' } })}
        columns={[
          { title: 'No. Order', dataIndex: 'order_number' },
          { title: 'Penerima', dataIndex: 'recipient_name' },
          {
            title: 'Status',
            dataIndex: 'status',
            render: (s: string) => <Tag color={statusColors[s]}>{s}</Tag>,
          },
          {
            title: 'Total',
            dataIndex: 'total',
            render: (v: number) => formatPrice(v),
          },
          {
            title: 'Tanggal',
            dataIndex: 'created_at',
            render: (v: string) => new Date(v).toLocaleString('id-ID'),
          },
        ]}
      />
    </div>
  );
}
