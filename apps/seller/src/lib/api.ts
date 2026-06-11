import { createMerchantApi } from '@owncommerce/sdk';
import { getAccessToken } from './auth';

const baseUrl = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

export const merchantApi = createMerchantApi(baseUrl, getAccessToken);
