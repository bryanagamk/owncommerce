import { PageHeader } from '@owncommerce/ui';
import type { Category, Product } from '@owncommerce/types';
import { Button, Card, Form, Input, InputNumber, Select, Switch, Upload, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { merchantApi } from '../lib/api';
import { fileUrl } from '../lib/format';

export default function ProductFormPage() {
  const { id } = useParams();
  const isNew = !id || id === 'new';
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [categories, setCategories] = useState<Category[]>([]);
  const [product, setProduct] = useState<Product | null>(null);

  useEffect(() => {
    merchantApi.listCategories().then(setCategories);
    if (!isNew && id) {
      merchantApi.getProduct(id).then((p) => {
        setProduct(p);
        form.setFieldsValue({
          ...p,
          quantity: p.inventory?.quantity,
          low_stock_threshold: p.inventory?.low_stock_threshold,
        });
      });
    }
  }, [id, isNew, form]);

  const onFinish = async (values: Record<string, unknown>) => {
    try {
      const body = {
        name: values.name,
        slug: values.slug,
        description: values.description,
        sku: values.sku,
        price: values.price,
        status: values.status,
        is_featured: values.is_featured,
        category_id: values.category_id,
        quantity: values.quantity,
        low_stock_threshold: values.low_stock_threshold,
      };
      if (isNew) {
        const created = await merchantApi.createProduct(body);
        message.success('Produk dibuat');
        navigate(`/products/${created.id}`);
      } else if (id) {
        await merchantApi.updateProduct(id, body);
        if (values.quantity !== undefined) {
          await merchantApi.updateInventory(id, {
            quantity: values.quantity as number,
            low_stock_threshold: values.low_stock_threshold as number,
          });
        }
        message.success('Produk diperbarui');
      }
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Gagal menyimpan');
    }
  };

  const onUpload = async (file: File) => {
    if (!id || isNew) {
      message.warning('Simpan produk terlebih dahulu');
      return false;
    }
    try {
      await merchantApi.uploadImage(id, file);
      message.success('Gambar diupload');
      const p = await merchantApi.getProduct(id);
      setProduct(p);
    } catch (e) {
      message.error(e instanceof Error ? e.message : 'Upload gagal');
    }
    return false;
  };

  return (
    <div>
      <PageHeader
        title={isNew ? 'Tambah Produk' : 'Edit Produk'}
        subtitle={isNew ? 'Buat produk baru' : product?.name}
      />
      <Card bordered style={{ borderColor: '#F0F0F0', maxWidth: 720 }}>
        <Form form={form} layout="vertical" onFinish={onFinish} initialValues={{ status: 'draft', is_featured: false, quantity: 0, low_stock_threshold: 5 }}>
          <Form.Item name="name" label="Nama Produk" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="slug" label="Slug">
            <Input placeholder="otomatis jika kosong" />
          </Form.Item>
          <Form.Item name="category_id" label="Kategori">
            <Select
              allowClear
              options={categories.map((c) => ({ value: c.id, label: c.name }))}
            />
          </Form.Item>
          <Form.Item name="sku" label="SKU">
            <Input />
          </Form.Item>
          <Form.Item name="price" label="Harga (IDR)" rules={[{ required: true }]}>
            <InputNumber style={{ width: '100%' }} min={0} />
          </Form.Item>
          <Form.Item name="description" label="Deskripsi">
            <Input.TextArea rows={4} />
          </Form.Item>
          <Form.Item name="status" label="Status">
            <Select
              options={[
                { value: 'draft', label: 'Draft' },
                { value: 'active', label: 'Active' },
                { value: 'archived', label: 'Archived' },
              ]}
            />
          </Form.Item>
          <Form.Item name="is_featured" label="Featured" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="quantity" label="Stok">
            <InputNumber style={{ width: '100%' }} min={0} />
          </Form.Item>
          <Form.Item name="low_stock_threshold" label="Ambang Stok Menipis">
            <InputNumber style={{ width: '100%' }} min={0} />
          </Form.Item>
          <Button type="primary" htmlType="submit">
            Simpan
          </Button>
        </Form>
      </Card>
      {!isNew && (
        <Card title="Gambar Produk" style={{ marginTop: 24, maxWidth: 720, borderColor: '#F0F0F0' }} bordered>
          <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap', marginBottom: 16 }}>
            {product?.images?.map((img) => (
              <img key={img.id} src={fileUrl(img.url)} alt="" style={{ width: 96, height: 96, objectFit: 'cover', borderRadius: 8, border: '1px solid #F0F0F0' }} />
            ))}
          </div>
          <Upload beforeUpload={onUpload} showUploadList={false} accept="image/*">
            <Button>Upload Gambar</Button>
          </Upload>
        </Card>
      )}
    </div>
  );
}
