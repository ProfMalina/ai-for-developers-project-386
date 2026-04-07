import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { OwnerDashboard } from '@/pages/owner/OwnerDashboard';

describe('OwnerDashboard Page', () => {
  it('should render dashboard title', () => {
    render(<OwnerDashboard />);
    const title = screen.getByText(/Панель управления/i);
    expect(title).toBeInTheDocument();
  });

  it('should render tabs', () => {
    render(<OwnerDashboard />);
    const eventsTab = screen.getByText(/Типы встреч/i);
    const slotsTab = screen.getByText(/Расписание/i);
    const bookingsTabs = screen.getAllByText(/Бронирования/i);

    expect(eventsTab).toBeInTheDocument();
    expect(slotsTab).toBeInTheDocument();
    expect(bookingsTabs.length).toBeGreaterThan(0);
  });

  it('should show events panel by default', () => {
    render(<OwnerDashboard />);
    const eventsPanel = screen.getByText(/Управление типами встреч/i);
    expect(eventsPanel).toBeInTheDocument();
  });
});
