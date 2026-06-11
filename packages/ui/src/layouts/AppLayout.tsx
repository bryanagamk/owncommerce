import {
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons';
import { Avatar, Dropdown, Layout, Menu, Typography } from 'antd';
import type { MenuProps } from 'antd';
import { useState, type ReactNode } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';

const { Header, Sider, Content } = Layout;

export interface NavItem {
  key: string;
  icon?: ReactNode;
  label: string;
  path: string;
}

interface AppLayoutProps {
  brand: string;
  navItems: NavItem[];
  userName?: string;
  onLogout?: () => void;
}

export function AppLayout({ brand, navItems, userName, onLogout }: AppLayoutProps) {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  const selectedKey =
    navItems.find((item) => location.pathname === item.path || location.pathname.startsWith(`${item.path}/`))
      ?.key ?? navItems[0]?.key;

  const menuItems: MenuProps['items'] = navItems.map((item) => ({
    key: item.key,
    icon: item.icon,
    label: item.label,
    onClick: () => navigate(item.path),
  }));

  const userMenu: MenuProps['items'] = [
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: 'Keluar',
      onClick: onLogout,
    },
  ];

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        collapsedWidth={0}
        breakpoint="lg"
        onBreakpoint={(broken) => setCollapsed(broken)}
        width={240}
        style={{
          borderRight: '1px solid #F0F0F0',
          background: '#FFFFFF',
        }}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            padding: collapsed ? '0 16px' : '0 24px',
            borderBottom: '1px solid #F0F0F0',
          }}
        >
          <Typography.Text strong style={{ fontSize: collapsed ? 14 : 18, whiteSpace: 'nowrap' }}>
            {collapsed ? brand.charAt(0) : brand}
          </Typography.Text>
        </div>
        <Menu
          mode="inline"
          selectedKeys={selectedKey ? [selectedKey] : []}
          items={menuItems}
          style={{ border: 'none', marginTop: 8 }}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            padding: '0 24px',
            background: '#FFFFFF',
            borderBottom: '1px solid #F0F0F0',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
          }}
        >
          <button
            type="button"
            onClick={() => setCollapsed(!collapsed)}
            style={{
              border: 'none',
              background: 'transparent',
              cursor: 'pointer',
              fontSize: 18,
              color: '#595959',
            }}
          >
            {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          </button>
          <Dropdown menu={{ items: userMenu }} placement="bottomRight">
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}>
              <Avatar style={{ backgroundColor: '#1677FF' }}>{userName?.charAt(0) ?? 'U'}</Avatar>
              {!collapsed && <Typography.Text>{userName ?? 'User'}</Typography.Text>}
            </div>
          </Dropdown>
        </Header>
        <Content style={{ margin: '16px 12px', minHeight: 280 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
