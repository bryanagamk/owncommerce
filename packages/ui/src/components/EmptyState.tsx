import { Empty, type EmptyProps } from 'antd';
import type { ReactNode } from 'react';

interface EmptyStateProps {
  title?: string;
  description?: string;
  action?: ReactNode;
  image?: EmptyProps['image'];
}

export function EmptyState({
  title = 'Belum ada data',
  description,
  action,
  image = Empty.PRESENTED_IMAGE_SIMPLE,
}: EmptyStateProps) {
  return (
    <Empty
      image={image}
      description={
        <div>
          <div style={{ fontWeight: 500, color: '#262626' }}>{title}</div>
          {description && <div style={{ color: '#8C8C8C', marginTop: 4 }}>{description}</div>}
        </div>
      }
    >
      {action}
    </Empty>
  );
}
