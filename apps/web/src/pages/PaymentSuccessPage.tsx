import { Button, Result } from 'antd';
import { Link } from 'react-router-dom';

export default function PaymentSuccessPage() {
  return (
    <div style={{ padding: 48 }}>
      <Result
        status="success"
        title="Pembayaran Berhasil"
        subTitle="Terima kasih! Pesanan Anda sedang diproses."
        extra={[
          <Link to="/account/orders" key="orders">
            <Button type="primary">Lihat Pesanan</Button>
          </Link>,
          <Link to="/" key="home">
            <Button>Kembali Belanja</Button>
          </Link>,
        ]}
      />
    </div>
  );
}
