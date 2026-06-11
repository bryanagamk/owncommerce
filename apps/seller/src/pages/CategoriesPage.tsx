import { PageHeader } from '@owncommerce/ui';
import type { Category } from '@owncommerce/types';
import { Button, Form, Input, Modal, Switch, Table, message } from 'antd';
import { useEffect, useState } from 'react';
import { merchantApi } from '../lib/api';

export default function CategoriesPage() {
  const [items, setItems] = useState<Category[]>([]);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<Category | null>(null);
  const [form] = Form.useForm();

  const load = () => merchantApi.listCategories().then(setItems);

  useEffect(() => {
    load();
  }, []);

  const openModal = (cat?: Category) => {
    setEditing(cat ?? null);
    form.setFieldsValue(cat ?? { is_active: true, sort_order: 0 });
    setOpen(true);
  };

  const onSave = async (values: Partial<Category>) => {
    try {
      if (editing) {
        await merchantApi.updateCategory(editing.id, values);
      } else {
        await merchantApi.createCategory(values);
      }
      message.success('Kategori disimpan');
      setOpen(false);
      load();
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal menyimpan');
    }
  };

  const onDelete = async (id: string) => {
    try {
      await merchantApi.deleteCategory(id);
      message.success('Kategori dihapus');
      load();
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal menghapus');
    }
  };

  return (
    <div>
      <PageHeader
        title="Kategori"
        subtitle="Kelola kategori produk"
        extra={
          <Button type="primary" onClick={() => openModal()}>
            Tambah Kategori
          </Button>
        }
      />
      <Table
        rowKey="id"
        dataSource={items}
        bordered
        style={{ borderColor: '#F0F0F0' }}
        columns={[
          { title: 'Nama', dataIndex: 'name' },
          { title: 'Slug', dataIndex: 'slug' },
          {
            title: 'Aktif',
            dataIndex: 'is_active',
            render: (v: boolean) => (v ? 'Ya' : 'Tidak'),
          },
          {
            title: 'Aksi',
            render: (_, record) => (
              <>
                <Button type="link" onClick={() => openModal(record)}>
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
      <Modal
        title={editing ? 'Edit Kategori' : 'Tambah Kategori'}
        open={open}
        onCancel={() => setOpen(false)}
        onOk={() => form.submit()}
        okText="Simpan"
      >
        <Form form={form} layout="vertical" onFinish={onSave}>
          <Form.Item name="name" label="Nama" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="slug" label="Slug">
            <Input placeholder="otomatis dari nama jika kosong" />
          </Form.Item>
          <Form.Item name="description" label="Deskripsi">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="is_active" label="Aktif" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
