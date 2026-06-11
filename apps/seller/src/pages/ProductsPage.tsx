import { PageHeader } from '@owncommerce/ui';
import type { Product } from '@owncommerce/types';
import { Button, Input, Select, Table, Tag, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { merchantApi } from '../lib/api';
import { fileUrl, formatPrice } from '../lib/format';

export default function ProductsPage() {
  const navigate = useNavigate();
  const [items, setItems] = useState<Product[]>([]);
  const [status, setStatus] = useState<string>('');
  const [search, setSearch] = useState('');

  const load = () =>
    merchantApi
      .listProducts({ status: status || undefined, q: search || undefined, limit: 50 })
      .then((res) => setItems(res.data));

  useEffect(() => {
    load();
  }, [status]);

  const onDelete = async (id: string) => {
    try {
      await merchantApi.deleteProduct(id);
      message.success('Produk dihapus');
      load();
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal menghapus');
    }
  };

  return (
    <div>
      <PageHeader
        title="Produk"
        subtitle="Kelola katalog produk toko"
        extra={
          <Button type="primary" onClick={() => navigate('/products/new')}>
            Tambah Produk
          </Button>
        }
      />
      <div style={{ display: 'flex', gap: 12, marginBottom: 16 }}>
        <Input.Search
          placeholder="Cari produk..."
          allowClear
          onSearch={(v) => {
            setSearch(v);
            merchantApi.listProducts({ q: v || undefined, limit: 50 }).then((res) => setItems(res.data));
          }}
          style={{ maxWidth: 280 }}
        />
        <Select
          placeholder="Status"
          allowClear
          style={{ width: 160 }}
          value={status || undefined}
          onChange={(v) => setStatus(v ?? '')}
          options={[
            { value: 'active', label: 'Active' },
            { value: 'draft', label: 'Draft' },
            { value: 'archived', label: 'Archived' },
          ]}
        />
      </div>
      <Table
        rowKey="id"
        dataSource={items}
        bordered
        columns={[
          {
            title: 'Gambar',
            render: (_, r) =>
              r.images?.[0]?.url ? (
                <img src={fileUrl(r.images[0].url)} alt="" style={{ width: 48, height: 48, objectFit: 'cover', borderRadius: 4 }} />
              ) : (
                '-'
              ),
          },
          { title: 'Nama', dataIndex: 'name' },
          { title: 'SKU', dataIndex: 'sku' },
          {
            title: 'Harga',
            dataIndex: 'price',
            render: (v: number) => formatPrice(v),
          },
          {
            title: 'Stok',
            render: (_, r) => r.inventory?.quantity ?? 0,
          },
          {
            title: 'Status',
            dataIndex: 'status',
            render: (s: string) => <Tag color={s === 'active' ? 'green' : 'default'}>{s}</Tag>,
          },
          {
            title: 'Aksi',
            render: (_, record) => (
              <>
                <Button type="link" onClick={() => navigate(`/products/${record.id}`)}>
                  Edit
                </Button>
                <Button type="link" danger onClick={() => onDelete(record.id)}>
                  Hapus
                </Button>
              </>
            ),
          },
        ]}
      />
    </div>
  );
}
