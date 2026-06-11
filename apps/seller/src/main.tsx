import { OwnCommerceProvider } from '@owncommerce/ui';
import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <OwnCommerceProvider>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </OwnCommerceProvider>
  </React.StrictMode>,
);
