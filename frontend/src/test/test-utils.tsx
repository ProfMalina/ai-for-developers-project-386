/* eslint-disable react-refresh/only-export-components */
import React from 'react';
import { render } from '@testing-library/react';
import type { RenderOptions } from '@testing-library/react';
import { MantineProvider } from '@mantine/core';
import { BrowserRouter } from 'react-router-dom';
import { Notifications } from '@mantine/notifications';

// Custom wrapper with providers
function AllTheProviders({ children }: { children: React.ReactNode }) {
  return (
    <MantineProvider>
      <BrowserRouter>
        <Notifications />
        {children}
      </BrowserRouter>
    </MantineProvider>
  );
}

// Custom render function
const customRender = (
  ui: React.ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>
) => render(ui, { wrapper: AllTheProviders, ...options });

// Re-export everything
export * from '@testing-library/react';

// Override render method
export { customRender as render };
