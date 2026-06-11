const apiBase = (import.meta.env.VITE_API_URL ?? 'http://localhost:8080').replace(/\/$/, '');

export function fileUrl(url?: string) {
  if (!url) return '';
  if (url.startsWith('http://') || url.startsWith('https://')) return url;
  return `${apiBase}${url.startsWith('/') ? url : `/${url}`}`;
}

export function formatPrice(amount: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount);
}
