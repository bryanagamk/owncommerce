import { PageHeader } from '@owncommerce/ui';
import { Card, Col, Row, Statistic, Table, Tag } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { merchantApi } from '../lib/api';
import { formatPrice } from '../lib/format';

export default function DashboardPage() {
  const navigate = useNavigate();
  const [stats, setStats] = useState({ products: 0, orders: 0, lowStock: 0 });
  const [recentOrders, setRecentOrders] = useState<unknown[]>([]);

  useEffect(() => {
    Promise.all([
      merchantApi.listProducts({ limit: 100 }),
      merchantApi.listOrders({ limit: 5 }),
    ]).then(([productsRes, ordersRes]) => {
      const products = productsRes.data;
      setStats({
        products: products.filter((p) => p.status === 'active').length,
        orders: ordersRes.meta.total,
        lowStock: products.filter(
          (p) => p.inventory && p.inventory.quantity <= p.inventory.low_stock_threshold,
        ).length,
      });
      setRecentOrders(ordersRes.data);
    });
  }, []);

  return (
    <div>
      <PageHeader title="Dashboard" subtitle="Ringkasan toko Anda" />
      <div
        style={{
          background: 'linear-gradient(135deg, #1a1a2e 0%, #16213e 100%)',
          borderRadius: 12,
          padding: '32px 40px',
          marginBottom: 24,
          color: '#fff',
        }}
      >
        <h2 style={{ margin: 0, fontSize: 24, fontWeight: 600 }}>Selamat datang di OwnCommerce</h2>
        <p style={{ margin: '8px 0 0', opacity: 0.8 }}>
          Kelola produk, pesanan, dan pengaturan toko dari satu dashboard.
        </p>
      </div>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={8}>
          <Card bordered style={{ borderColor: '#F0F0F0' }}>
            <Statistic title="Produk Aktif" value={stats.products} />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card bordered style={{ borderColor: '#F0F0F0' }}>
            <Statistic title="Total Pesanan" value={stats.orders} />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card bordered style={{ borderColor: '#F0F0F0' }}>
            <Statistic title="Stok Menipis" value={stats.lowStock} valueStyle={{ color: stats.lowStock > 0 ? '#cf1322' : undefined }} />
          </Card>
        </Col>
      </Row>
      <Card title="Pesanan Terbaru" bordered style={{ borderColor: '#F0F0F0' }}>
        <Table
          rowKey="id"
          dataSource={recentOrders as { id: string; order_number: string; status: string; total: number; created_at: string }[]}
          pagination={false}
          onRow={(record) => ({ onClick: () => navigate(`/orders/${record.id}`), style: { cursor: 'pointer' } })}
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
      </Card>
    </div>
  );
}
