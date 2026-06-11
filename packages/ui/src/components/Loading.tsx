import { Spin } from 'antd';

export function Loading({ tip = 'Memuat...' }: { tip?: string }) {
  return (
    <div style={{ display: 'flex', justifyContent: 'center', padding: 64 }}>
      <Spin size="large" tip={tip} />
    </div>
  );
}
