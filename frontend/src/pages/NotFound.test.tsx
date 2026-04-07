import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { NotFound } from '@/pages/NotFound';

describe('NotFound Page', () => {
  it('should render 404 message', () => {
    render(<NotFound />);
    const message = screen.getByText(/Страница не найдена/i);
    expect(message).toBeInTheDocument();
  });

  it('should render home link', () => {
    render(<NotFound />);
    const homeLink = screen.getByText(/На главную/i);
    expect(homeLink).toBeInTheDocument();
  });
});
