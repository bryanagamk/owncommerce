import {
  AppstoreOutlined,
  HomeOutlined,
  SettingOutlined,
  ShoppingOutlined,
  TagsOutlined,
} from '@ant-design/icons';
import { AppLayout } from '@owncommerce/ui';
import { useNavigate } from 'react-router-dom';
import { clearAuth, getUserName } from '../lib/auth';

const navItems = [
  { key: 'dashboard', icon: <HomeOutlined />, label: 'Dashboard', path: '/' },
  { key: 'products', icon: <AppstoreOutlined />, label: 'Produk', path: '/products' },
  { key: 'categories', icon: <TagsOutlined />, label: 'Kategori', path: '/categories' },
  { key: 'orders', icon: <ShoppingOutlined />, label: 'Pesanan', path: '/orders' },
  { key: 'settings', icon: <SettingOutlined />, label: 'Pengaturan Toko', path: '/settings' },
];

export function SellerLayout() {
  const navigate = useNavigate();

  return (
    <AppLayout
      brand="OwnCommerce"
      navItems={navItems}
      userName={getUserName()}
      onLogout={() => {
        clearAuth();
        navigate('/login');
      }}
    />
  );
}
