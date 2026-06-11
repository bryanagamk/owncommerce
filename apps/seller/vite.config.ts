import react from '@vitejs/plugin-react';
import path from 'path';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@owncommerce/ui': path.resolve(__dirname, '../../packages/ui/src'),
      '@owncommerce/sdk': path.resolve(__dirname, '../../packages/sdk/src'),
      '@owncommerce/types': path.resolve(__dirname, '../../packages/types/src'),
    },
  },
  server: {
    port: 5173,
    fs: { allow: ['../..'] },
    proxy: {
      '/v1': { target: 'http://localhost:8080', changeOrigin: true },
      '/files': { target: 'http://localhost:8080', changeOrigin: true },
    },
  },
});
