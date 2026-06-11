import type { Product } from '@owncommerce/types';
import { Card, Col, Row, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { fileUrl, formatPrice } from '../lib/format';

export default function ProductsPage() {
  const [params] = useSearchParams();
  const q = params.get('q') ?? '';
  const [products, setProducts] = useState<Product[]>([]);

  useEffect(() => {
    storefrontApi.listProducts({ q: q || undefined, limit: 48 }).then((res) => setProducts(res.data));
  }, [q]);

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto', padding: '24px 16px' }}>
      <Typography.Title level={3}>{q ? `Hasil: "${q}"` : 'Semua Produk'}</Typography.Title>
      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        {products.map((p) => (
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
