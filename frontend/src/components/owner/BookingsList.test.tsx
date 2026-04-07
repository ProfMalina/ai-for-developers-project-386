import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { BookingsList } from '@/components/owner/BookingsList';

describe('BookingsList Component', () => {
  it('should render component title', async () => {
    render(<BookingsList />);
    const title = await screen.findByText(/Список бронирований/i);
    expect(title).toBeInTheDocument();
  });

  it('should render bookings from API', async () => {
    render(<BookingsList />);
    const guestName = await screen.findByText('Иван Иванов');
    expect(guestName).toBeInTheDocument();
  });

  it('should display booking email', async () => {
    render(<BookingsList />);
    const email = await screen.findByText('ivan@example.com');
    expect(email).toBeInTheDocument();
  });

  it('should display cancel button for upcoming bookings', async () => {
    render(<BookingsList />);
    const cancelButton = await screen.findByText(/Отменить бронирование/i);
    expect(cancelButton).toBeInTheDocument();
  });

  it('should display booking time information', async () => {
    render(<BookingsList />);
    await screen.findByText('Иван Иванов');

    const timeLabel = screen.getByText(/Время:/i);
    expect(timeLabel).toBeInTheDocument();
  });

  it('should display event type ID', async () => {
    render(<BookingsList />);
    await screen.findByText('Иван Иванов');

    const idLabel = screen.getByText(/ID типа встречи:/i);
    expect(idLabel).toBeInTheDocument();
  });

  it('should display booking creation date', async () => {
    render(<BookingsList />);
    await screen.findByText('Иван Иванов');

    const createdLabel = screen.getByText(/Создано:/i);
    expect(createdLabel).toBeInTheDocument();
  });

  it('should display status badge', async () => {
    render(<BookingsList />);
    await screen.findByText('Иван Иванов');

    // Badge should show either "Предстоящая" or "Прошедшая"
    const badges = screen.getAllByText(/Предстоящая|Прошедшая/i);
    expect(badges.length).toBeGreaterThan(0);
  });
});
