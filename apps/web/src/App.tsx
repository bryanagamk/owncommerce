import { Route, Routes } from 'react-router-dom';
import { StoreLayout } from './components/StoreLayout';
import AccountLoginPage from './pages/AccountLoginPage';
import AccountOrdersPage from './pages/AccountOrdersPage';
import AccountRegisterPage from './pages/AccountRegisterPage';
import CartPage from './pages/CartPage';
import CheckoutPage from './pages/CheckoutPage';
import HomePage from './pages/HomePage';
import PaymentFailedPage from './pages/PaymentFailedPage';
import PaymentPage from './pages/PaymentPage';
import PaymentSuccessPage from './pages/PaymentSuccessPage';
import ProductDetailPage from './pages/ProductDetailPage';
import ProductsPage from './pages/ProductsPage';

export default function App() {
  return (
    <Routes>
      <Route element={<StoreLayout />}>
        <Route index element={<HomePage />} />
        <Route path="products" element={<ProductsPage />} />
        <Route path="products/:slug" element={<ProductDetailPage />} />
        <Route path="cart" element={<CartPage />} />
        <Route path="checkout" element={<CheckoutPage />} />
        <Route path="payment/:orderId" element={<PaymentPage />} />
        <Route path="payment/success" element={<PaymentSuccessPage />} />
        <Route path="payment/failed" element={<PaymentFailedPage />} />
        <Route path="account/login" element={<AccountLoginPage />} />
        <Route path="account/register" element={<AccountRegisterPage />} />
        <Route path="account/orders" element={<AccountOrdersPage />} />
      </Route>
    </Routes>
  );
}
