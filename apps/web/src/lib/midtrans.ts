declare global {
  interface Window {
    snap?: {
      pay: (token: string, options?: {
        onSuccess?: (result: unknown) => void;
        onPending?: (result: unknown) => void;
        onError?: (result: unknown) => void;
        onClose?: () => void;
      }) => void;
    };
  }
}

export function openSnapPayment(
  token: string,
  callbacks: {
    onSuccess?: () => void;
    onPending?: () => void;
    onError?: () => void;
    onClose?: () => void;
  },
) {
  if (!window.snap) {
    throw new Error('Midtrans Snap belum dimuat. Cek VITE_MIDTRANS_CLIENT_KEY.');
  }
  window.snap.pay(token, {
    onSuccess: () => callbacks.onSuccess?.(),
    onPending: () => callbacks.onPending?.(),
    onError: () => callbacks.onError?.(),
    onClose: () => callbacks.onClose?.(),
  });
}
