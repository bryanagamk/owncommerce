const ACCESS_KEY = 'owncommerce_customer_access';
const REFRESH_KEY = 'owncommerce_customer_refresh';
const NAME_KEY = 'owncommerce_customer_name';

export function saveCustomerAuth(access: string, refresh: string, name: string) {
  localStorage.setItem(ACCESS_KEY, access);
  localStorage.setItem(REFRESH_KEY, refresh);
  localStorage.setItem(NAME_KEY, name);
}

export function getCustomerToken() {
  return localStorage.getItem(ACCESS_KEY);
}

export function getCustomerName() {
  return localStorage.getItem(NAME_KEY);
}

export function clearCustomerAuth() {
  localStorage.removeItem(ACCESS_KEY);
  localStorage.removeItem(REFRESH_KEY);
  localStorage.removeItem(NAME_KEY);
}

export function isCustomerLoggedIn() {
  return !!getCustomerToken();
}
