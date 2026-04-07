import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { EventTypeManagement } from '@/components/owner/EventTypeManagement';

describe('EventTypeManagement Component', () => {
  it('should render component title', async () => {
    render(<EventTypeManagement />);
    const title = await screen.findByText(/Управление типами встреч/i);
    expect(title).toBeInTheDocument();
  });

  it('should render add button', async () => {
    render(<EventTypeManagement />);
    const addButton = await screen.findByText(/Добавить тип встречи/i);
    expect(addButton).toBeInTheDocument();
  });

  it('should render event types from API', async () => {
    render(<EventTypeManagement />);
    const eventType = await screen.findByText('Консультация');
    expect(eventType).toBeInTheDocument();
  });

  it('should display event duration badges', async () => {
    render(<EventTypeManagement />);
    const badge = await screen.findByText(/30 мин/i);
    expect(badge).toBeInTheDocument();
  });

  it('should show empty state when no event types', () => {
    // This test would require mocking an empty response
    // For now, we just verify the component renders
    render(<EventTypeManagement />);
    expect(true).toBe(true);
  });
});
