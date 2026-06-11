import { Button, Result } from 'antd';
import { Link } from 'react-router-dom';

export default function PaymentFailedPage() {
  return (
    <div style={{ padding: 48 }}>
      <Result
        status="error"
        title="Pembayaran Gagal"
        subTitle="Silakan coba lagi atau hubungi toko."
        extra={
          <Link to="/cart">
            <Button type="primary">Kembali ke Keranjang</Button>
          </Link>
        }
      />
    </div>
  );
}
