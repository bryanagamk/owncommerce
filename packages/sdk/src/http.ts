import type { ApiResponse } from '@owncommerce/types';

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export interface RequestOptions {
  method?: string;
  body?: unknown;
  token?: string | null;
  tenantSlug?: string;
  cartSession?: string;
  formData?: FormData;
}

export async function request<T>(
  baseUrl: string,
  path: string,
  options: RequestOptions = {},
): Promise<T> {
  const headers: Record<string, string> = {};

  if (!options.formData) {
    headers['Content-Type'] = 'application/json';
  }

  if (options.token) {
    headers['Authorization'] = `Bearer ${options.token}`;
  }
  if (options.tenantSlug) {
    headers['X-Tenant-Slug'] = options.tenantSlug;
  }
  if (options.cartSession) {
    headers['X-Cart-Session'] = options.cartSession;
  }

  const res = await fetch(`${baseUrl}${path}`, {
    method: options.method ?? 'GET',
    headers,
    body: options.formData ?? (options.body ? JSON.stringify(options.body) : undefined),
  });

  const json = (await res.json()) as ApiResponse<T>;

  if (!res.ok || !json.success) {
    throw new ApiError(json.message ?? 'Request failed', res.status);
  }

  return json.data as T;
}

export async function requestPaginated<T>(
  baseUrl: string,
  path: string,
  options: RequestOptions = {},
): Promise<{ data: T; meta: { total: number; limit: number; offset: number } }> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };

  if (options.token) {
    headers['Authorization'] = `Bearer ${options.token}`;
  }
  if (options.tenantSlug) {
    headers['X-Tenant-Slug'] = options.tenantSlug;
  }
  if (options.cartSession) {
    headers['X-Cart-Session'] = options.cartSession;
  }

  const res = await fetch(`${baseUrl}${path}`, {
    method: options.method ?? 'GET',
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  const json = await res.json();

  if (!res.ok || !json.success) {
    throw new ApiError(json.message ?? 'Request failed', res.status);
  }

  return { data: json.data as T, meta: json.meta };
}
