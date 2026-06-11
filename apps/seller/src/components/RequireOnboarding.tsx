import { Navigate, Outlet } from 'react-router-dom';
import { isOnboardingComplete } from '../lib/auth';

export function RequireOnboarding() {
  if (!isOnboardingComplete()) {
    return <Navigate to="/onboarding" replace />;
  }
  return <Outlet />;
}
