const CART_SESSION_KEY = 'owncommerce_cart_session';

export function getCartSession(): string {
  let session = localStorage.getItem(CART_SESSION_KEY);
  if (!session) {
    session = crypto.randomUUID();
    localStorage.setItem(CART_SESSION_KEY, session);
  }
  return session;
}
