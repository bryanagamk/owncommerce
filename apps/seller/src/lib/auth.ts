const ACCESS_KEY = 'owncommerce_seller_access';
const REFRESH_KEY = 'owncommerce_seller_refresh';
const USER_KEY = 'owncommerce_seller_user';

export function saveAuth(access: string, refresh: string, userName: string) {
  localStorage.setItem(ACCESS_KEY, access);
  localStorage.setItem(REFRESH_KEY, refresh);
  localStorage.setItem(USER_KEY, userName);
}

export function getAccessToken() {
  return localStorage.getItem(ACCESS_KEY);
}

export function getUserName() {
  return localStorage.getItem(USER_KEY) ?? 'Merchant';
}

export function clearAuth() {
  localStorage.removeItem(ACCESS_KEY);
  localStorage.removeItem(REFRESH_KEY);
  localStorage.removeItem(USER_KEY);
}

export function isAuthenticated() {
  return !!getAccessToken();
}
