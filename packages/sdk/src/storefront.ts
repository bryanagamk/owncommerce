import type {
  CartView,
  Category,
  Customer,
  CustomerAuthResult,
  Order,
  PayResult,
  Product,
  Tenant,
} from '@owncommerce/types';
import { request, requestPaginated } from './http';

export interface StorefrontContext {
  tenantSlug: string;
  getToken: () => string | null;
  getCartSession: () => string;
}

export function createStorefrontApi(baseUrl: string, ctx: StorefrontContext) {
  const opts = () => ({
    tenantSlug: ctx.tenantSlug,
    token: ctx.getToken(),
    cartSession: ctx.getCartSession(),
  });

  return {
    home: () =>
      request<{ store: Tenant; featured_products: Product[] }>(
        baseUrl,
        '/v1/storefront/home',
        opts(),
      ),

    listCategories: () =>
      request<Category[]>(baseUrl, '/v1/storefront/categories', opts()),

    listProducts: (params?: { category_id?: string; q?: string; limit?: number; offset?: number }) => {
      const qs = new URLSearchParams();
      if (params?.category_id) qs.set('category_id', params.category_id);
      if (params?.q) qs.set('q', params.q);
      if (params?.limit) qs.set('limit', String(params.limit));
      if (params?.offset) qs.set('offset', String(params.offset));
      const query = qs.toString() ? `?${qs}` : '';
      return requestPaginated<Product[]>(baseUrl, `/v1/storefront/products${query}`, opts());
    },

    getProduct: (slug: string) =>
      request<Product>(baseUrl, `/v1/storefront/products/${slug}`, opts()),

    register: (body: { email: string; password: string; name: string; phone?: string }) =>
      request<CustomerAuthResult>(baseUrl, '/v1/storefront/auth/register', {
        method: 'POST',
        ...opts(),
        body,
      }),

    login: (body: { email: string; password: string }) =>
      request<CustomerAuthResult>(baseUrl, '/v1/storefront/auth/login', {
        method: 'POST',
        ...opts(),
        body,
      }),

    me: () => request<Customer>(baseUrl, '/v1/storefront/me', opts()),

    getCart: () => request<CartView>(baseUrl, '/v1/storefront/cart', opts()),

    addCartItem: (body: { product_id: string; quantity: number }) =>
      request<CartView>(baseUrl, '/v1/storefront/cart/items', {
        method: 'POST',
        ...opts(),
        body,
      }),

    updateCartItem: (id: string, body: { quantity: number }) =>
      request<CartView>(baseUrl, `/v1/storefront/cart/items/${id}`, {
        method: 'PATCH',
        ...opts(),
        body,
      }),

    removeCartItem: (id: string) =>
      request<CartView>(baseUrl, `/v1/storefront/cart/items/${id}`, {
        method: 'DELETE',
        ...opts(),
      }),

    checkout: (body: Record<string, string | number>) =>
      request<Order>(baseUrl, '/v1/storefront/checkout', {
        method: 'POST',
        ...opts(),
        body,
      }),

    payOrder: (orderId: string) =>
      request<PayResult>(baseUrl, `/v1/storefront/orders/${orderId}/pay`, {
        method: 'POST',
        ...opts(),
      }),

    getOrder: (orderId: string) =>
      request<Order>(baseUrl, `/v1/storefront/orders/${orderId}`, opts()),

    listOrders: () =>
      requestPaginated<Order[]>(baseUrl, '/v1/storefront/orders', opts()),
  };
}
