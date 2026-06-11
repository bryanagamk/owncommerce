import { Loading } from '@owncommerce/ui';
import { Button, Card, Typography, message } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { storefrontApi } from '../lib/api';
import { openSnapPayment } from '../lib/midtrans';

export default function PaymentPage() {
  const { orderId } = useParams();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!orderId) return;
    storefrontApi
      .payOrder(orderId)
      .then((result) => {
        setLoading(false);
        openSnapPayment(result.snap_token, {
          onSuccess: () => navigate('/payment/success'),
          onPending: () => navigate('/payment/success'),
          onError: () => navigate('/payment/failed'),
          onClose: () => message.info('Pembayaran ditutup'),
        });
      })
      .catch((e) => {
        setLoading(false);
        message.error(e instanceof Error ? e.message : 'Gagal memuat pembayaran');
      });
  }, [orderId, navigate]);

  if (loading) return <Loading tip="Menyiapkan pembayaran..." />;

  return (
    <div style={{ maxWidth: 480, margin: '64px auto', padding: 16 }}>
      <Card bordered style={{ borderColor: '#F0F0F0', textAlign: 'center' }}>
        <Typography.Paragraph>Jendela pembayaran Midtrans seharusnya terbuka.</Typography.Paragraph>
        <Button type="primary" onClick={() => orderId && navigate(`/payment/${orderId}`)}>
          Buka Ulang Pembayaran
        </Button>
      </Card>
    </div>
  );
}
