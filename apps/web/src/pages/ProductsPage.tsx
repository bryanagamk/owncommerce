import type { Category, Product } from '@owncommerce/types';
import { EmptyState } from '@owncommerce/ui';
import { Card, Col, Row, Select, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { fileUrl, formatPrice } from '../lib/format';

export default function ProductsPage() {
  const [params, setParams] = useSearchParams();
  const q = params.get('q') ?? '';
  const categoryId = params.get('category_id') ?? '';
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);

  useEffect(() => {
    storefrontApi.listCategories().then(setCategories).catch(() => {});
  }, []);

  useEffect(() => {
    storefrontApi
      .listProducts({
        q: q || undefined,
        category_id: categoryId || undefined,
        limit: 48,
      })
      .then((res) => setProducts(res.data));
  }, [q, categoryId]);

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto', padding: '24px 16px' }}>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: 12, alignItems: 'center', justifyContent: 'space-between' }}>
        <Typography.Title level={3} style={{ margin: 0 }}>
          {q ? `Hasil: "${q}"` : 'Semua Produk'}
        </Typography.Title>
        <Select
          placeholder="Filter kategori"
          allowClear
          style={{ minWidth: 200 }}
          value={categoryId || undefined}
          onChange={(v) => {
            const next = new URLSearchParams(params);
            if (v) next.set('category_id', v);
            else next.delete('category_id');
            setParams(next);
          }}
          options={categories.map((c) => ({ value: c.id, label: c.name }))}
        />
      </div>

      {products.length === 0 ? (
        <div style={{ marginTop: 48 }}>
          <EmptyState title="Produk tidak ditemukan" description="Coba kata kunci atau kategori lain" />
        </div>
      ) : (
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
      )}
    </div>
  );
}
