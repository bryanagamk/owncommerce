import type {
  AuthResult,
  Category,
  MeResponse,
  Order,
  Product,
  Tenant,
} from '@owncommerce/types';
import { request, requestPaginated } from './http';

export function createMerchantApi(baseUrl: string, getToken: () => string | null) {
  const token = () => getToken();

  return {
    register: (body: {
      email: string;
      password: string;
      name: string;
      store_name: string;
      slug: string;
    }) => request<AuthResult>(baseUrl, '/v1/auth/register', { method: 'POST', body }),

    login: (body: { email: string; password: string }) =>
      request<AuthResult>(baseUrl, '/v1/auth/login', { method: 'POST', body }),

    me: () => request<MeResponse>(baseUrl, '/v1/me', { token: token() }),

    getTenant: () =>
      request<{ tenant: Tenant; domains: unknown[] }>(baseUrl, '/v1/tenants/current', {
        token: token(),
      }),

    updateStore: (body: Record<string, string | undefined>) =>
      request<Tenant>(baseUrl, '/v1/merchant/store/settings', {
        method: 'PATCH',
        token: token(),
        body,
      }),

    listCategories: () =>
      request<Category[]>(baseUrl, '/v1/merchant/categories', { token: token() }),

    createCategory: (body: Partial<Category>) =>
      request<Category>(baseUrl, '/v1/merchant/categories', {
        method: 'POST',
        token: token(),
        body,
      }),

    updateCategory: (id: string, body: Partial<Category>) =>
      request<Category>(baseUrl, `/v1/merchant/categories/${id}`, {
        method: 'PATCH',
        token: token(),
        body,
      }),

    deleteCategory: (id: string) =>
      request<void>(baseUrl, `/v1/merchant/categories/${id}`, {
        method: 'DELETE',
        token: token(),
      }),

    listProducts: (params?: { status?: string; q?: string; limit?: number; offset?: number }) => {
      const qs = new URLSearchParams();
      if (params?.status) qs.set('status', params.status);
      if (params?.q) qs.set('q', params.q);
      if (params?.limit) qs.set('limit', String(params.limit));
      if (params?.offset) qs.set('offset', String(params.offset));
      const query = qs.toString() ? `?${qs}` : '';
      return requestPaginated<Product[]>(baseUrl, `/v1/merchant/products${query}`, {
        token: token(),
      });
    },

    getProduct: (id: string) =>
      request<Product>(baseUrl, `/v1/merchant/products/${id}`, { token: token() }),

    createProduct: (body: Record<string, unknown>) =>
      request<Product>(baseUrl, '/v1/merchant/products', {
        method: 'POST',
        token: token(),
        body,
      }),

    updateProduct: (id: string, body: Record<string, unknown>) =>
      request<Product>(baseUrl, `/v1/merchant/products/${id}`, {
        method: 'PATCH',
        token: token(),
        body,
      }),

    deleteProduct: (id: string) =>
      request<void>(baseUrl, `/v1/merchant/products/${id}`, {
        method: 'DELETE',
        token: token(),
      }),

    updateInventory: (id: string, body: { quantity?: number; low_stock_threshold?: number }) =>
      request<unknown>(baseUrl, `/v1/merchant/products/${id}/inventory`, {
        method: 'PATCH',
        token: token(),
        body,
      }),

    uploadImage: async (id: string, file: File) => {
      const formData = new FormData();
      formData.append('image', file);
      const headers: Record<string, string> = {};
      const t = token();
      if (t) headers['Authorization'] = `Bearer ${t}`;
      const res = await fetch(`${baseUrl}/v1/merchant/products/${id}/images`, {
        method: 'POST',
        headers,
        body: formData,
      });
      const json = await res.json();
      if (!res.ok || !json.success) throw new Error(json.message ?? 'Upload failed');
      return json.data;
    },

    listOrders: (params?: { status?: string; limit?: number; offset?: number }) => {
      const qs = new URLSearchParams();
      if (params?.status) qs.set('status', params.status);
      if (params?.limit) qs.set('limit', String(params.limit));
      if (params?.offset) qs.set('offset', String(params.offset));
      const query = qs.toString() ? `?${qs}` : '';
      return requestPaginated<Order[]>(baseUrl, `/v1/merchant/orders${query}`, {
        token: token(),
      });
    },

    getOrder: (id: string) =>
      request<Order>(baseUrl, `/v1/merchant/orders/${id}`, { token: token() }),

    updateOrderStatus: (id: string, body: { status: string; note?: string }) =>
      request<Order>(baseUrl, `/v1/merchant/orders/${id}/status`, {
        method: 'PATCH',
        token: token(),
        body,
      }),
  };
}
