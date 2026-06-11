import { Card, Typography } from 'antd';
import type { ReactNode } from 'react';

interface AuthLayoutProps {
  title: string;
  subtitle?: string;
  children: ReactNode;
}

export function AuthLayout({ title, subtitle, children }: AuthLayoutProps) {
  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: '#FAFAFA',
        padding: 24,
      }}
    >
      <Card
        style={{
          width: '100%',
          maxWidth: 420,
          border: '1px solid #F0F0F0',
          boxShadow: '0 1px 2px rgba(0,0,0,0.03)',
        }}
        bordered
      >
        <div style={{ marginBottom: 32, textAlign: 'center' }}>
          <Typography.Title level={3} style={{ marginBottom: 8 }}>
            {title}
          </Typography.Title>
          {subtitle && <Typography.Text type="secondary">{subtitle}</Typography.Text>}
        </div>
        {children}
      </Card>
    </div>
  );
}
