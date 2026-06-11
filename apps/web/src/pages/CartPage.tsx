import type { CartView } from '@owncommerce/types';
import { EmptyState } from '@owncommerce/ui';
import { Button, Card, InputNumber, Table, Typography, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { formatPrice } from '../lib/format';

export default function CartPage() {
  const navigate = useNavigate();
  const [cart, setCart] = useState<CartView | null>(null);

  const load = () => storefrontApi.getCart().then(setCart);

  useEffect(() => {
    load();
  }, []);

  const updateQty = async (id: string, quantity: number) => {
    try {
      const res = await storefrontApi.updateCartItem(id, { quantity });
      setCart(res);
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal update');
    }
  };

  const remove = async (id: string) => {
    try {
      const res = await storefrontApi.removeCartItem(id);
      setCart(res);
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal hapus');
    }
  };

  if (!cart) return null;

  return (
    <div style={{ maxWidth: 800, margin: '0 auto', padding: '32px 16px' }}>
      <Typography.Title level={3}>Keranjang</Typography.Title>
      {cart.items.length === 0 ? (
        <Card bordered style={{ borderColor: '#F0F0F0', padding: 24 }}>
          <EmptyState
            title="Keranjang kosong"
            description="Yuk, tambahkan produk favorit Anda"
            action={
              <Button type="primary" onClick={() => navigate('/products')}>
                Lihat Katalog
              </Button>
            }
          />
        </Card>
      ) : (
        <>
          <Table
            rowKey="id"
            dataSource={cart.items}
            pagination={false}
            columns={[
              { title: 'Produk', dataIndex: 'product_name' },
              {
                title: 'Harga',
                dataIndex: 'unit_price',
                render: (v: number) => formatPrice(v),
              },
              {
                title: 'Qty',
                render: (_, r) => (
                  <InputNumber
                    min={1}
                    value={r.quantity}
                    onChange={(v) => updateQty(r.id, v ?? 1)}
                  />
                ),
              },
              {
                title: 'Subtotal',
                render: (_, r) => formatPrice(r.unit_price * r.quantity),
              },
              {
                title: '',
                render: (_, r) => (
                  <Button type="link" danger onClick={() => remove(r.id)}>
                    Hapus
                  </Button>
                ),
              },
            ]}
          />
          <Card bordered style={{ marginTop: 16, borderColor: '#F0F0F0' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Typography.Title level={4} style={{ margin: 0 }}>
                Total: {formatPrice(cart.total)}
              </Typography.Title>
              <Button type="primary" size="large" onClick={() => navigate('/checkout')}>
                Checkout
              </Button>
            </div>
          </Card>
        </>
      )}
    </div>
  );
}
