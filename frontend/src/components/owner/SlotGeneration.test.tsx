import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { SlotGeneration } from '@/components/owner/SlotGeneration';

describe('SlotGeneration Component', () => {
  it('should render component title', () => {
    render(<SlotGeneration />);
    const title = screen.getByText(/Генерация расписания/i);
    expect(title).toBeInTheDocument();
  });

  it('should render create button', () => {
    render(<SlotGeneration />);
    const button = screen.getByText(/Создать слоты/i);
    expect(button).toBeInTheDocument();
  });

  it('should render description text', () => {
    render(<SlotGeneration />);
    const description = screen.getByText(/Создайте доступные слоты для бронирования/i);
    expect(description).toBeInTheDocument();
  });
});
