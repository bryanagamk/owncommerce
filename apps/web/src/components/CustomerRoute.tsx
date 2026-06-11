import { Navigate, Outlet } from 'react-router-dom';
import { isCustomerLoggedIn } from '../lib/auth';

export function CustomerRoute() {
  if (!isCustomerLoggedIn()) {
    return <Navigate to="/account/login" replace />;
  }
  return <Outlet />;
}
