import { createStorefrontApi } from '@owncommerce/sdk';
import { getCustomerToken } from './auth';
import { getCartSession } from './cart';
import { TENANT_SLUG } from './tenant';

const baseUrl = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

export const storefrontApi = createStorefrontApi(baseUrl, {
  tenantSlug: TENANT_SLUG,
  getToken: getCustomerToken,
  getCartSession,
});
