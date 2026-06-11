import type { Product } from '@owncommerce/types';
import { Button, Col, InputNumber, Row, Typography, message } from 'antd';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { fileUrl, formatPrice } from '../lib/format';

export default function ProductDetailPage() {
  const { slug } = useParams();
  const [product, setProduct] = useState<Product | null>(null);
  const [qty, setQty] = useState(1);

  useEffect(() => {
    if (slug) storefrontApi.getProduct(slug).then(setProduct);
  }, [slug]);

  const addToCart = async () => {
    if (!product) return;
    try {
      await storefrontApi.addCartItem({ product_id: product.id, quantity: qty });
      message.success('Ditambahkan ke keranjang');
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal menambah ke keranjang');
    }
  };

  if (!product) return null;

  return (
    <div style={{ maxWidth: 960, margin: '0 auto', padding: '32px 16px' }}>
      <Row gutter={[32, 32]}>
        <Col xs={24} md={12}>
          {product.images?.[0]?.url ? (
            <img
              src={fileUrl(product.images[0].url)}
              alt={product.name}
              style={{ width: '100%', borderRadius: 12, border: '1px solid #F0F0F0' }}
            />
          ) : (
            <div style={{ aspectRatio: '1', background: '#F5F5F5', borderRadius: 12 }} />
          )}
        </Col>
        <Col xs={24} md={12}>
          <Typography.Title level={2}>{product.name}</Typography.Title>
          <Typography.Title level={3} style={{ color: '#1677FF', marginTop: 0 }}>
            {formatPrice(product.price)}
          </Typography.Title>
          <Typography.Paragraph type="secondary">{product.description}</Typography.Paragraph>
          <Typography.Text>Stok: {product.inventory?.quantity ?? 0}</Typography.Text>
          <div style={{ marginTop: 24, display: 'flex', gap: 12, alignItems: 'center' }}>
            <InputNumber min={1} max={product.inventory?.quantity ?? 99} value={qty} onChange={(v) => setQty(v ?? 1)} />
            <Button type="primary" size="large" onClick={addToCart}>
              Tambah ke Keranjang
            </Button>
          </div>
        </Col>
      </Row>
    </div>
  );
}
