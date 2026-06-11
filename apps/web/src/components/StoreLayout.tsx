import { MenuOutlined, ShoppingCartOutlined, UserOutlined } from '@ant-design/icons';
import { Badge, Button, Drawer, Grid, Input, Layout, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { Outlet, Link, useNavigate } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { clearCustomerAuth, getCustomerName, isCustomerLoggedIn } from '../lib/auth';
import type { Tenant } from '@owncommerce/types';

const { Header, Content, Footer } = Layout;
const { useBreakpoint } = Grid;

export function StoreLayout() {
  const navigate = useNavigate();
  const screens = useBreakpoint();
  const isMobile = !screens.md;
  const [store, setStore] = useState<Tenant | null>(null);
  const [cartCount, setCartCount] = useState(0);
  const [menuOpen, setMenuOpen] = useState(false);
  const loggedIn = isCustomerLoggedIn();

  useEffect(() => {
    storefrontApi.home().then((res) => setStore(res.store));
    storefrontApi
      .getCart()
      .then((cart) => {
        const count = cart.items?.reduce((sum, i) => sum + i.quantity, 0) ?? 0;
        setCartCount(count);
      })
      .catch(() => {});
  }, []);

  const accountLinks = loggedIn ? (
    <>
      <Button type="link" block style={{ textAlign: 'left' }} onClick={() => { navigate('/account/profile'); setMenuOpen(false); }}>
        Profil
      </Button>
      <Button type="link" block style={{ textAlign: 'left' }} onClick={() => { navigate('/account/orders'); setMenuOpen(false); }}>
        Pesanan
      </Button>
      <Button
        type="link"
        block
        style={{ textAlign: 'left' }}
        onClick={() => {
          clearCustomerAuth();
          setMenuOpen(false);
          navigate('/');
        }}
      >
        Keluar
      </Button>
    </>
  ) : (
    <Button type="link" block style={{ textAlign: 'left' }} onClick={() => { navigate('/account/login'); setMenuOpen(false); }}>
      Masuk / Daftar
    </Button>
  );

  return (
    <Layout style={{ minHeight: '100vh', background: '#fff' }}>
      <Header
        style={{
          background: '#fff',
          borderBottom: '1px solid #F0F0F0',
          padding: isMobile ? '0 12px' : '0 24px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          height: 64,
          gap: 8,
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, minWidth: 0 }}>
          {isMobile && (
            <Button type="text" icon={<MenuOutlined />} onClick={() => setMenuOpen(true)} />
          )}
          <Link to="/" style={{ minWidth: 0 }}>
            <Typography.Text strong style={{ fontSize: isMobile ? 16 : 18 }} ellipsis>
              {store?.name ?? 'Toko'}
            </Typography.Text>
          </Link>
        </div>

        {!isMobile && (
          <Input.Search
            placeholder="Cari produk..."
            style={{ width: 220, maxWidth: '30vw' }}
            onSearch={(q) => navigate(`/products?q=${encodeURIComponent(q)}`)}
          />
        )}

        <div style={{ display: 'flex', alignItems: 'center', gap: isMobile ? 4 : 16 }}>
          <Badge count={cartCount}>
            <Button type="text" icon={<ShoppingCartOutlined />} onClick={() => navigate('/cart')} />
          </Badge>
          {!isMobile && loggedIn && (
            <Typography.Text style={{ maxWidth: 120 }} ellipsis>
              {getCustomerName()}
            </Typography.Text>
          )}
          {!isMobile && loggedIn && (
            <>
              <Button type="link" onClick={() => navigate('/account/profile')}>
                Profil
              </Button>
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
          )}
          {!isMobile && !loggedIn && (
            <Button type="text" icon={<UserOutlined />} onClick={() => navigate('/account/login')}>
              Masuk
            </Button>
          )}
        </div>
      </Header>

      <Drawer title={store?.name ?? 'Menu'} placement="left" open={menuOpen} onClose={() => setMenuOpen(false)}>
        <Input.Search
          placeholder="Cari produk..."
          style={{ marginBottom: 16 }}
          onSearch={(q) => {
            navigate(`/products?q=${encodeURIComponent(q)}`);
            setMenuOpen(false);
          }}
        />
        <Button type="link" block style={{ textAlign: 'left' }} onClick={() => { navigate('/products'); setMenuOpen(false); }}>
          Katalog
        </Button>
        <Button type="link" block style={{ textAlign: 'left' }} onClick={() => { navigate('/cart'); setMenuOpen(false); }}>
          Keranjang ({cartCount})
        </Button>
        {accountLinks}
      </Drawer>

      <Content style={{ background: '#FAFAFA', minHeight: 'calc(100vh - 128px)' }}>
        <Outlet context={{ setCartCount }} />
      </Content>
      <Footer
        style={{
          textAlign: 'center',
          background: '#fff',
          borderTop: '1px solid #F0F0F0',
          color: '#8C8C8C',
          padding: isMobile ? '16px 12px' : '24px 50px',
        }}
      >
        © {new Date().getFullYear()} {store?.name ?? 'OwnCommerce'}
      </Footer>
    </Layout>
  );
}
