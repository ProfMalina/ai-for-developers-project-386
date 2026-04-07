import '@testing-library/jest-dom';
import { setupMockServer } from './mocks';

// Setup MSW server
setupMockServer();

// Polyfill for matchMedia (needed by Mantine)
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => false,
  }),
});
