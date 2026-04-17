import { beforeEach, describe, expect, it } from 'vitest';
import userEvent from '@testing-library/user-event';
import { http, HttpResponse } from 'msw';
import { render, screen } from '@/test/test-utils';
import { server } from '@/test/mocks';
import { BookingsList } from '@/components/owner/BookingsList';

const createBookingAtOffset = (daysOffset: number) => {
  const start = new Date();
  start.setUTCDate(start.getUTCDate() + daysOffset);
  start.setUTCHours(10, 0, 0, 0);

  const end = new Date(start);
  end.setUTCMinutes(end.getUTCMinutes() + 30);

  const createdAt = new Date();
  createdAt.setUTCHours(8, 0, 0, 0);

  return {
    id: `booking-${daysOffset}`,
    eventTypeId: 'event-type-1',
    startTime: start.toISOString(),
    endTime: end.toISOString(),
    guestName: 'Иван Иванов',
    guestEmail: 'ivan@example.com',
    createdAt: createdAt.toISOString(),
  };
};

describe('BookingsList', () => {
  beforeEach(() => {
    server.resetHandlers();
  });

  it('renders upcoming booking details', async () => {
    render(<BookingsList />);

    expect(await screen.findByText('Иван Иванов')).toBeInTheDocument();
    expect(screen.getByText('ivan@example.com')).toBeInTheDocument();
    expect(screen.getByText(/Предстоящая/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Отменить бронирование/i })).toBeEnabled();
  });

  it('shows empty state when there are no bookings', async () => {
    server.use(
      http.get('*/api/bookings', () =>
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

    render(<BookingsList />);

    expect(await screen.findByText(/Бронирований пока нет/i)).toBeInTheDocument();
  });

  it('disables cancellation for past bookings', async () => {
    server.use(
      http.get('*/api/bookings', () =>
        HttpResponse.json({
          items: [createBookingAtOffset(-1)],
          pagination: {
            page: 1,
            pageSize: 10,
            totalItems: 1,
            totalPages: 1,
            hasNext: false,
            hasPrev: false,
          },
        })
      )
    );

    render(<BookingsList />);

    expect(await screen.findByText(/Прошедшая/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Отменить бронирование/i })).toBeDisabled();
  });

  it('cancels booking after confirmation', async () => {
    const user = userEvent.setup();
    render(<BookingsList />);

    await screen.findByText('Иван Иванов');
    await user.click(screen.getByRole('button', { name: /Отменить бронирование/i }));
    await user.click(await screen.findByRole('button', { name: /Да, отменить/i }));

    expect(await screen.findByText(/Бронирование отменено/i)).toBeInTheDocument();
  });

  it('renders pagination when multiple pages exist', async () => {
    server.use(
      http.get('*/api/bookings', () =>
        HttpResponse.json({
          items: [createBookingAtOffset(1)],
          pagination: {
            page: 1,
            pageSize: 10,
            totalItems: 25,
            totalPages: 3,
            hasNext: true,
            hasPrev: false,
          },
        })
      )
    );

    render(<BookingsList />);

    expect(await screen.findByRole('button', { name: '2' })).toBeInTheDocument();
  });
});
