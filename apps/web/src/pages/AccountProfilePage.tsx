import type { Address, Customer } from '@owncommerce/types';
import { EmptyState } from '@owncommerce/ui';
import { Button, Card, Form, Input, Modal, Switch, Table, Typography, message } from 'antd';
import { useEffect, useState } from 'react';
import { storefrontApi } from '../lib/api';

type AddressFormValues = Omit<Address, 'id'>;

export default function AccountProfilePage() {
  const [profileForm] = Form.useForm();
  const [addressForm] = Form.useForm<AddressFormValues>();
  const [customer, setCustomer] = useState<Customer | null>(null);
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);

  const load = async () => {
    const [me, addr] = await Promise.all([storefrontApi.me(), storefrontApi.listAddresses()]);
    setCustomer(me);
    profileForm.setFieldsValue({ name: me.name, phone: me.phone, email: me.email });
    setAddresses(addr);
  };

  useEffect(() => {
    load().catch((e) => message.error(e instanceof Error ? e.message : 'Gagal memuat profil'));
  }, [profileForm]);

  const saveProfile = async (values: { name: string; phone?: string }) => {
    try {
      const updated = await storefrontApi.updateMe(values);
      setCustomer(updated);
      message.success('Profil diperbarui');
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal menyimpan profil');
    }
  };

  const openAddressModal = (address?: Address) => {
    setEditingId(address?.id ?? null);
    addressForm.setFieldsValue(
      address ?? {
        label: '',
        recipient_name: '',
        phone: '',
        address_line: '',
        city: '',
        province: '',
        postal_code: '',
        is_default: false,
      },
    );
    setModalOpen(true);
  };

  const saveAddress = async (values: AddressFormValues) => {
    try {
      if (editingId) {
        await storefrontApi.updateAddress(editingId, values);
        message.success('Alamat diperbarui');
      } else {
        await storefrontApi.createAddress(values);
        message.success('Alamat ditambahkan');
      }
      setModalOpen(false);
      setEditingId(null);
      addressForm.resetFields();
      const addr = await storefrontApi.listAddresses();
      setAddresses(addr);
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal menyimpan alamat');
    }
  };

  const removeAddress = async (id: string) => {
    try {
      await storefrontApi.deleteAddress(id);
      message.success('Alamat dihapus');
      setAddresses((prev) => prev.filter((a) => a.id !== id));
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal menghapus alamat');
    }
  };

  return (
    <div style={{ maxWidth: 800, margin: '0 auto', padding: '24px 16px' }}>
      <Typography.Title level={3}>Profil Saya</Typography.Title>

      <Card title="Data Akun" bordered style={{ borderColor: '#F0F0F0', marginBottom: 24 }}>
        <Form form={profileForm} layout="vertical" onFinish={saveProfile}>
          <Form.Item name="email" label="Email">
            <Input disabled />
          </Form.Item>
          <Form.Item name="name" label="Nama" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="Telepon">
            <Input />
          </Form.Item>
          <Button type="primary" htmlType="submit">
            Simpan Profil
          </Button>
        </Form>
      </Card>

      <Card
        title="Alamat Pengiriman"
        bordered
        style={{ borderColor: '#F0F0F0' }}
        extra={
          <Button type="primary" onClick={() => openAddressModal()}>
            Tambah Alamat
          </Button>
        }
      >
        {addresses.length === 0 ? (
          <EmptyState
            title="Belum ada alamat"
            description="Tambahkan alamat untuk mempercepat checkout"
            action={
              <Button type="primary" onClick={() => openAddressModal()}>
                Tambah Alamat
              </Button>
            }
          />
        ) : (
          <Table
            rowKey="id"
            dataSource={addresses}
            pagination={false}
            scroll={{ x: 600 }}
            columns={[
              { title: 'Label', dataIndex: 'label', render: (v: string) => v || '-' },
              { title: 'Penerima', dataIndex: 'recipient_name' },
              { title: 'Kota', dataIndex: 'city' },
              {
                title: 'Default',
                dataIndex: 'is_default',
                render: (v: boolean) => (v ? 'Ya' : '-'),
              },
              {
                title: '',
                render: (_, r) => (
                  <div style={{ display: 'flex', gap: 8 }}>
                    <Button type="link" onClick={() => openAddressModal(r)}>
                      Edit
                    </Button>
                    <Button type="link" danger onClick={() => removeAddress(r.id)}>
                      Hapus
                    </Button>
                  </div>
                ),
              },
            ]}
          />
        )}
      </Card>

      <Modal
        title={editingId ? 'Edit Alamat' : 'Tambah Alamat'}
        open={modalOpen}
        onCancel={() => {
          setModalOpen(false);
          setEditingId(null);
        }}
        footer={null}
        destroyOnClose
      >
        <Form form={addressForm} layout="vertical" onFinish={saveAddress}>
          <Form.Item name="label" label="Label">
            <Input placeholder="Rumah, Kantor, dll." />
          </Form.Item>
          <Form.Item name="recipient_name" label="Nama Penerima" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="Telepon" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="address_line" label="Alamat" rules={[{ required: true }]}>
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="city" label="Kota" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="province" label="Provinsi" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="postal_code" label="Kode Pos" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="is_default" label="Jadikan default" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Button type="primary" htmlType="submit" block>
            Simpan
          </Button>
        </Form>
      </Modal>

      {customer && (
        <Typography.Paragraph type="secondary" style={{ marginTop: 16 }}>
          Masuk sebagai {customer.email}
        </Typography.Paragraph>
      )}
    </div>
  );
}
