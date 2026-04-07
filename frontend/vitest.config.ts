import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html', 'lcov'],
      thresholds: {
        lines: 45,
        branches: 30,
        functions: 40,
        statements: 45,
      },
      include: [
        'src/**/*.ts',
        'src/**/*.tsx',
      ],
      exclude: [
        'src/main.tsx',
        'src/test/**',
        'src/**/*.d.ts',
        'src/types/**',
        'src/**/*.test.ts',
        'src/**/*.test.tsx',
      ],
    },
    include: ['src/**/*.test.ts', 'src/**/*.test.tsx'],
    alias: {
      '@test': resolve(__dirname, './src/test'),
      '@': resolve(__dirname, './src'),
    },
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    },
  },
});
