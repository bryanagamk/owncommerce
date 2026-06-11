import { ShoppingCartOutlined, UserOutlined } from '@ant-design/icons';
import { Badge, Button, Input, Layout, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { Outlet, Link, useNavigate } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { clearCustomerAuth, getCustomerName, isCustomerLoggedIn } from '../lib/auth';
import type { Tenant } from '@owncommerce/types';

const { Header, Content, Footer } = Layout;

export function StoreLayout() {
  const navigate = useNavigate();
  const [store, setStore] = useState<Tenant | null>(null);
  const [cartCount, setCartCount] = useState(0);
  const loggedIn = isCustomerLoggedIn();

  useEffect(() => {
    storefrontApi.home().then((res) => setStore(res.store));
    storefrontApi.getCart().then((cart) => {
      const count = cart.items?.reduce((sum, i) => sum + i.quantity, 0) ?? 0;
      setCartCount(count);
    }).catch(() => {});
  }, []);

  return (
    <Layout style={{ minHeight: '100vh', background: '#fff' }}>
      <Header
        style={{
          background: '#fff',
          borderBottom: '1px solid #F0F0F0',
          padding: '0 24px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          height: 64,
        }}
      >
        <Link to="/">
          <Typography.Text strong style={{ fontSize: 18 }}>
            {store?.name ?? 'Toko'}
          </Typography.Text>
        </Link>
        <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
          <Input.Search
            placeholder="Cari produk..."
            style={{ width: 220 }}
            onSearch={(q) => navigate(`/products?q=${encodeURIComponent(q)}`)}
          />
          <Badge count={cartCount}>
            <Button type="text" icon={<ShoppingCartOutlined />} onClick={() => navigate('/cart')} />
          </Badge>
          {loggedIn ? (
            <>
              <Typography.Text>{getCustomerName()}</Typography.Text>
              <Button type="link" onClick={() => navigate('/account/orders')}>
                Pesanan
              </Button>
              <Button
                type="link"
                onClick={() => {
                  clearCustomerAuth();
                  navigate('/');
                }}
              >
                Keluar
              </Button>
            </>
          ) : (
            <Button type="text" icon={<UserOutlined />} onClick={() => navigate('/account/login')}>
              Masuk
            </Button>
          )}
        </div>
      </Header>
      <Content style={{ background: '#FAFAFA', minHeight: 'calc(100vh - 128px)' }}>
        <Outlet context={{ setCartCount }} />
      </Content>
      <Footer style={{ textAlign: 'center', background: '#fff', borderTop: '1px solid #F0F0F0', color: '#8C8C8C' }}>
        © {new Date().getFullYear()} {store?.name ?? 'OwnCommerce'}
      </Footer>
    </Layout>
  );
}
