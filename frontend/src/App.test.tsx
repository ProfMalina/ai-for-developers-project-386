import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { MemoryRouter } from 'react-router-dom';
import App from '@/App';

const renderApp = (route: string) => render(
  <MantineProvider>
    <MemoryRouter initialEntries={[route]}>
      <Notifications />
      <App />
    </MemoryRouter>
  </MantineProvider>
);

describe('App routing', () => {
  it('renders guest home route', async () => {
    renderApp('/');

    expect(await screen.findByText(/Забронировать встречу/i)).toBeInTheDocument();
  });

  it('renders owner dashboard route', async () => {
    renderApp('/owner');

    expect(await screen.findByText(/Панель управления/i)).toBeInTheDocument();
  });

  it('renders not found route for unknown paths', async () => {
    renderApp('/missing-route');

    expect(await screen.findByText(/Страница не найдена/i)).toBeInTheDocument();
  });
});
