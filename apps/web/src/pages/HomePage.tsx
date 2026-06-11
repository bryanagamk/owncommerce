import type { Product } from '@owncommerce/types';
import { Button, Card, Col, Row, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { fileUrl, formatPrice } from '../lib/format';

export default function HomePage() {
  const [storeName, setStoreName] = useState('');
  const [description, setDescription] = useState('');
  const [featured, setFeatured] = useState<Product[]>([]);

  useEffect(() => {
    storefrontApi.home().then((res) => {
      setStoreName(res.store.name);
      setDescription(res.store.description ?? '');
      setFeatured(res.featured_products ?? []);
    });
  }, []);

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto', padding: '24px 16px' }}>
      <div
        style={{
          background: 'linear-gradient(135deg, #1a1a2e 0%, #16213e 100%)',
          borderRadius: 12,
          padding: '48px 40px',
          marginBottom: 40,
          color: '#fff',
        }}
      >
        <Typography.Title level={2} style={{ color: '#fff', margin: 0 }}>
          {storeName}
        </Typography.Title>
        <Typography.Paragraph style={{ color: 'rgba(255,255,255,0.8)', marginTop: 8, maxWidth: 480 }}>
          {description || 'Selamat berbelanja di toko kami.'}
        </Typography.Paragraph>
        <Link to="/products">
          <Button type="default" style={{ marginTop: 16 }}>
            Lihat Katalog
          </Button>
        </Link>
      </div>
      <Typography.Title level={4}>Produk Unggulan</Typography.Title>
      <Row gutter={[16, 16]}>
        {featured.map((p) => (
          <Col xs={12} sm={8} md={6} key={p.id}>
            <Link to={`/products/${p.slug}`}>
              <Card
                hoverable
                bordered
                style={{ borderColor: '#F0F0F0' }}
                cover={
                  p.images?.[0]?.url ? (
                    <img src={fileUrl(p.images[0].url)} alt={p.name} style={{ height: 160, objectFit: 'cover' }} />
                  ) : (
                    <div style={{ height: 160, background: '#F5F5F5' }} />
                  )
                }
              >
                <Card.Meta title={p.name} description={formatPrice(p.price)} />
              </Card>
            </Link>
          </Col>
        ))}
      </Row>
    </div>
  );
}
