import type { ThemeConfig } from 'antd';

export const ownCommerceTheme: ThemeConfig = {
  token: {
    colorPrimary: '#1677FF',
    colorBgContainer: '#FFFFFF',
    colorBgLayout: '#FAFAFA',
    colorBorder: '#F0F0F0',
    colorText: '#262626',
    colorTextSecondary: '#8C8C8C',
    borderRadius: 8,
    fontFamily: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
  },
  components: {
    Layout: {
      siderBg: '#FFFFFF',
      headerBg: '#FFFFFF',
      bodyBg: '#FAFAFA',
      triggerBg: '#FFFFFF',
    },
    Menu: {
      itemBg: 'transparent',
      itemSelectedBg: '#F5F5F5',
      itemHoverBg: '#FAFAFA',
      itemColor: '#595959',
      itemSelectedColor: '#262626',
    },
    Card: {
      paddingLG: 24,
    },
    Table: {
      headerBg: '#FAFAFA',
      borderColor: '#F0F0F0',
    },
  },
};
