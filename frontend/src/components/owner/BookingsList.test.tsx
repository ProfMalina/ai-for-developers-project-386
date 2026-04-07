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
});
