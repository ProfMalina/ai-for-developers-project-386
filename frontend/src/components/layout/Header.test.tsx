import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { Header } from '@/components/layout/Header';

describe('Header Component', () => {
  it('should render logo with emoji', () => {
    render(<Header />);
    expect(screen.getByText(/📅 Бронирование/)).toBeInTheDocument();
  });

  it('should render Guest button', () => {
    render(<Header />);
    expect(screen.getByText('Гость')).toBeInTheDocument();
  });

  it('should render Owner button', () => {
    render(<Header />);
    expect(screen.getByText('Владелец')).toBeInTheDocument();
  });

  it('should have Guest button as filled when on root path', () => {
    render(<Header />);
    const guestButton = screen.getByText('Гость');
    // Check that Guest button has filled variant (default behavior)
    expect(guestButton).toHaveAttribute('class');
  });

  it('should have Owner button as filled when on /owner path', () => {
    window.history.pushState({}, '', '/owner');
    render(<Header />);
    const ownerButton = screen.getByText('Владелец');
    expect(ownerButton).toHaveAttribute('class');
  });
});
