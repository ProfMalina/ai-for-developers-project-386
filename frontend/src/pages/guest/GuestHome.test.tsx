import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { GuestHome } from '@/pages/guest/GuestHome';

describe('GuestHome Page', () => {
  it('should render page title', async () => {
    render(<GuestHome />);
    const title = await screen.findByText(/Забронировать встречу/i);
    expect(title).toBeInTheDocument();
  });

  it('should render description text', async () => {
    render(<GuestHome />);
    const description = await screen.findByText(/Выберите тип встречи и удобное для вас время/i);
    expect(description).toBeInTheDocument();
  });

  it('should render event type cards from API', async () => {
    render(<GuestHome />);
    const consultationCard = await screen.findByText('Консультация');
    expect(consultationCard).toBeInTheDocument();
  });

  it('should render booking buttons for each event type', async () => {
    render(<GuestHome />);
    const bookButtons = await screen.findAllByText(/Забронировать/i);
    expect(bookButtons.length).toBeGreaterThan(0);
  });

  it('should display event duration in badges', async () => {
    render(<GuestHome />);
    const durationBadge = await screen.findByText(/30 мин/i);
    expect(durationBadge).toBeInTheDocument();
  });
});
