import { OwnCommerceProvider } from '@owncommerce/ui';
import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import './index.css';

const clientKey = import.meta.env.VITE_MIDTRANS_CLIENT_KEY;
const snapScript = document.getElementById('midtrans-snap');
if (snapScript && clientKey) {
  snapScript.setAttribute('data-client-key', clientKey);
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <OwnCommerceProvider>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </OwnCommerceProvider>
  </React.StrictMode>,
);
