export interface ApiResponse<T = unknown> {
  success: boolean;
  message?: string;
  data?: T;
  error?: unknown;
}

export interface PaginatedMeta {
  total: number;
  limit: number;
  offset: number;
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: T;
  meta: PaginatedMeta;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  phone?: string;
}

export interface Tenant {
  id: string;
  name: string;
  slug: string;
  status: string;
  description?: string;
  logo_url?: string;
  contact_email?: string;
  contact_phone?: string;
  address?: string;
  city?: string;
  province?: string;
  postal_code?: string;
}

export interface TenantDomain {
  id: string;
  domain: string;
  type: string;
  is_primary: boolean;
}

export interface AuthResult {
  user: User;
  tenant?: Tenant;
  tokens: TokenPair;
  roles?: string[];
  permissions?: string[];
}

export interface MeResponse {
  user: User;
  tenant_id?: string;
  roles: string[];
  permissions: string[];
}

export interface Category {
  id: string;
  tenant_id: string;
  name: string;
  slug: string;
  description?: string;
  is_active: boolean;
  sort_order: number;
}

export interface ProductImage {
  id: string;
  url: string;
  path: string;
  sort_order: number;
}

export interface Inventory {
  id: string;
  quantity: number;
  low_stock_threshold: number;
}

export interface Product {
  id: string;
  tenant_id: string;
  category_id?: string;
  name: string;
  slug: string;
  description?: string;
  sku?: string;
  price: number;
  status: string;
  is_featured: boolean;
  images?: ProductImage[];
  inventory?: Inventory;
}

export interface Customer {
  id: string;
  email: string;
  name: string;
  phone?: string;
}

export interface CustomerAuthResult {
  customer: Customer;
  tokens: TokenPair;
}

export interface CartItemView {
  id: string;
  product_id: string;
  quantity: number;
  unit_price: number;
  product_name: string;
  product_slug: string;
  image_url?: string;
}

export interface CartView {
  Cart: { id: string };
  items: CartItemView[];
  total: number;
}

export interface OrderItem {
  id: string;
  product_id: string;
  product_name: string;
  sku?: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
}

export interface Order {
  id: string;
  order_number: string;
  status: string;
  payment_status: string;
  subtotal: number;
  shipping_cost: number;
  total: number;
  recipient_name: string;
  recipient_phone: string;
  shipping_address: string;
  shipping_city: string;
  shipping_province: string;
  shipping_postal_code: string;
  customer_email?: string;
  items?: OrderItem[];
  created_at: string;
}

export interface PayResult {
  snap_token: string;
  redirect_url?: string;
  payment_id: string;
  order_id: string;
}

export interface Address {
  id: string;
  label?: string;
  recipient_name: string;
  phone: string;
  address_line: string;
  city: string;
  province: string;
  postal_code: string;
  is_default: boolean;
}
