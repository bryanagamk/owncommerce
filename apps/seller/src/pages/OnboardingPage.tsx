import { PageHeader } from '@owncommerce/ui';
import { Button, Card, Form, Input, InputNumber, Steps, Typography, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { merchantApi } from '../lib/api';
import { completeOnboarding } from '../lib/auth';

export default function OnboardingPage() {
  const navigate = useNavigate();
  const [step, setStep] = useState(0);
  const [storeForm] = Form.useForm();
  const [categoryForm] = Form.useForm();
  const [productForm] = Form.useForm();
  const [tenantSlug, setTenantSlug] = useState('');

  useEffect(() => {
    merchantApi.getTenant().then((res) => {
      storeForm.setFieldsValue(res.tenant);
      setTenantSlug(res.tenant.slug);
    });
  }, [storeForm]);

  const saveStore = async () => {
    try {
      const values = await storeForm.validateFields();
      await merchantApi.updateStore(values);
      message.success('Profil toko disimpan');
      setStep(1);
    } catch (e) {
      if (e && typeof e === 'object' && 'errorFields' in e) return;
      message.error(e instanceof Error ? e.message : 'Gagal menyimpan');
    }
  };

  const saveCategory = async () => {
    try {
      const values = await categoryForm.validateFields();
      await merchantApi.createCategory({ ...values, is_active: true, sort_order: 0 });
      message.success('Kategori ditambahkan');
      setStep(2);
    } catch (e) {
      if (e && typeof e === 'object' && 'errorFields' in e) return;
      message.error(e instanceof Error ? e.message : 'Gagal menyimpan kategori');
    }
  };

  const saveProduct = async () => {
    try {
      const values = await productForm.validateFields();
      const product = await merchantApi.createProduct({
        ...values,
        status: 'active',
        is_featured: true,
      });
      if (values.initial_stock != null) {
        await merchantApi.updateInventory(product.id, { quantity: values.initial_stock });
      }
      message.success('Produk pertama ditambahkan');
      finish();
    } catch (e) {
      if (e && typeof e === 'object' && 'errorFields' in e) return;
      message.error(e instanceof Error ? e.message : 'Gagal menyimpan produk');
    }
  };

  const finish = () => {
    completeOnboarding();
    navigate('/');
  };

  const skipToDashboard = () => {
    completeOnboarding();
    navigate('/');
  };

  return (
    <div style={{ maxWidth: 720, margin: '0 auto', padding: '32px 16px' }}>
      <PageHeader
        title="Setup Toko"
        subtitle="Lengkapi toko Anda dalam beberapa langkah"
      />
      <Steps
        current={step}
        style={{ marginBottom: 32 }}
        items={[
          { title: 'Profil Toko' },
          { title: 'Kategori' },
          { title: 'Produk Pertama' },
        ]}
        responsive
      />

      {step === 0 && (
        <Card bordered style={{ borderColor: '#F0F0F0' }}>
          <Typography.Paragraph type="secondary">
            Informasi ini akan tampil di storefront pelanggan Anda.
          </Typography.Paragraph>
          <Form form={storeForm} layout="vertical">
            <Form.Item name="name" label="Nama Toko" rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Form.Item name="description" label="Deskripsi" rules={[{ required: true }]}>
              <Input.TextArea rows={3} placeholder="Ceritakan tentang toko Anda" />
            </Form.Item>
            <Form.Item name="contact_phone" label="Telepon Kontak">
              <Input />
            </Form.Item>
            <Form.Item name="city" label="Kota">
              <Input />
            </Form.Item>
            <Button type="primary" onClick={saveStore}>
              Lanjut
            </Button>
          </Form>
        </Card>
      )}

      {step === 1 && (
        <Card bordered style={{ borderColor: '#F0F0F0' }}>
          <Typography.Paragraph type="secondary">
            Buat kategori pertama (opsional — bisa dilewati).
          </Typography.Paragraph>
          <Form form={categoryForm} layout="vertical">
            <Form.Item name="name" label="Nama Kategori" rules={[{ required: true }]}>
              <Input placeholder="Bunga Segar" />
            </Form.Item>
            <Form.Item
              name="slug"
              label="Slug"
              rules={[{ required: true }, { pattern: /^[a-z0-9-]+$/, message: 'Huruf kecil, angka, strip' }]}
            >
              <Input placeholder="bunga-segar" />
            </Form.Item>
            <div style={{ display: 'flex', gap: 8 }}>
              <Button type="primary" onClick={saveCategory}>
                Lanjut
              </Button>
              <Button onClick={() => setStep(2)}>Lewati</Button>
            </div>
          </Form>
        </Card>
      )}

      {step === 2 && (
        <Card bordered style={{ borderColor: '#F0F0F0' }}>
          <Typography.Paragraph type="secondary">
            Tambahkan produk pertama agar toko siap dijual (opsional).
          </Typography.Paragraph>
          <Form form={productForm} layout="vertical" initialValues={{ price: 50000, initial_stock: 10 }}>
            <Form.Item name="name" label="Nama Produk" rules={[{ required: true }]}>
              <Input />
            </Form.Item>
            <Form.Item name="price" label="Harga (IDR)" rules={[{ required: true }]}>
              <InputNumber style={{ width: '100%' }} min={0} />
            </Form.Item>
            <Form.Item name="initial_stock" label="Stok Awal">
              <InputNumber style={{ width: '100%' }} min={0} />
            </Form.Item>
            <Form.Item name="description" label="Deskripsi">
              <Input.TextArea rows={2} />
            </Form.Item>
            <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
              <Button type="primary" onClick={saveProduct}>
                Selesai & ke Dashboard
              </Button>
              <Button onClick={finish}>Lewati ke Dashboard</Button>
            </div>
          </Form>
        </Card>
      )}

      {tenantSlug && (
        <Card style={{ marginTop: 24, borderColor: '#F0F0F0' }} bordered>
          <Typography.Text type="secondary">URL toko Anda (dev):</Typography.Text>
          <Typography.Paragraph copyable style={{ marginBottom: 0 }}>
            http://localhost:5174 — tenant: {tenantSlug}
          </Typography.Paragraph>
        </Card>
      )}

      <div style={{ marginTop: 16, textAlign: 'center' }}>
        <Button type="link" onClick={skipToDashboard}>
          Lewati setup, ke dashboard
        </Button>
      </div>
    </div>
  );
}
