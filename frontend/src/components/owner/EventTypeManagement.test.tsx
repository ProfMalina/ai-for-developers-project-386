import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import userEvent from '@testing-library/user-event';
import { http, HttpResponse } from 'msw';
import { render, screen } from '@/test/test-utils';
import { server } from '@/test/mocks';
import { EventTypeManagement } from '@/components/owner/EventTypeManagement';

vi.mock('@mantine/core', async () => {
  const actual = await vi.importActual<typeof import('@mantine/core')>('@mantine/core');
  return {
    ...actual,
    Modal: ({ opened, title, children }: { opened: boolean; title: string; children: ReactNode }) => (
      opened ? <div role="dialog" aria-label={title}>{children}</div> : null
    ),
  };
});

describe('EventTypeManagement', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('renders fetched event types', async () => {
    render(<EventTypeManagement />);

    expect(await screen.findByText('Консультация')).toBeInTheDocument();
    expect(screen.getByText('Встреча')).toBeInTheDocument();
    expect(screen.getByText(/Индивидуальная консультация по проекту/i)).toBeInTheDocument();
  });

  it('shows empty state when there are no event types', async () => {
    server.use(
      http.get('*/api/event-types', () =>
        HttpResponse.json({
          items: [],
          pagination: {
            page: 1,
            pageSize: 10,
            totalItems: 0,
            totalPages: 1,
            hasNext: false,
            hasPrev: false,
          },
        })
      )
    );

    render(<EventTypeManagement />);

    expect(await screen.findByText(/Типы встреч еще не созданы/i)).toBeInTheDocument();
  });

  it('validates required fields when creating an event type', async () => {
    const user = userEvent.setup();
    render(<EventTypeManagement />);

    await screen.findByText('Консультация');
    await user.click(screen.getByRole('button', { name: /Добавить тип встречи/i }));
    await user.click(screen.getByRole('button', { name: /^Создать$/i }));

    expect(await screen.findByText(/Все поля обязательны для заполнения/i)).toBeInTheDocument();
  });

  it('creates a new event type successfully', async () => {
    const user = userEvent.setup();
    render(<EventTypeManagement />);

    await screen.findByText('Консультация');
    await user.click(screen.getByRole('button', { name: /Добавить тип встречи/i }));
    await user.type(screen.getByLabelText(/Название/i), 'Стратегическая сессия');
    await user.type(screen.getByLabelText(/Описание/i), 'Разбор целей и следующих шагов');
    await user.clear(screen.getByLabelText(/Длительность/i));
    await user.type(screen.getByLabelText(/Длительность/i), '45');
    await user.click(screen.getByRole('button', { name: /^Создать$/i }));

    expect(await screen.findByText(/Тип встречи создан/i)).toBeInTheDocument();
  });

  it('deletes an event type after confirmation', async () => {
    const user = userEvent.setup();
    vi.spyOn(window, 'confirm').mockReturnValue(true);

    render(<EventTypeManagement />);

    await screen.findByText('Консультация');
    await user.click(screen.getByLabelText(/Удалить тип встречи Консультация/i));

    expect(window.confirm).toHaveBeenCalled();
    expect(await screen.findByText(/Тип встречи удален/i)).toBeInTheDocument();
  });

  it('renders pagination when API returns multiple pages', async () => {
    server.use(
      http.get('*/api/event-types', () =>
        HttpResponse.json({
          items: [
            { id: 'event-type-1', name: 'Консультация', description: 'Индивидуальная консультация по проекту', durationMinutes: 30 },
          ],
          pagination: {
            page: 1,
            pageSize: 10,
            totalItems: 30,
            totalPages: 3,
            hasNext: true,
            hasPrev: false,
          },
        })
      )
    );

    render(<EventTypeManagement />);

    expect(await screen.findByRole('button', { name: '2' })).toBeInTheDocument();
  });
});
