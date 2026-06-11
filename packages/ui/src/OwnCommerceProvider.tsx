import { ConfigProvider } from 'antd';
import idID from 'antd/locale/id_ID';
import type { ReactNode } from 'react';
import { ownCommerceTheme } from './theme';

export function OwnCommerceProvider({ children }: { children: ReactNode }) {
  return (
    <ConfigProvider theme={ownCommerceTheme} locale={idID}>
      {children}
    </ConfigProvider>
  );
}
