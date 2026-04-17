import { beforeEach, describe, expect, it, vi } from 'vitest';
import userEvent from '@testing-library/user-event';
import { fireEvent } from '@testing-library/react';
import { http, HttpResponse } from 'msw';
import { render, screen, waitFor } from '@/test/test-utils';
import { server } from '@/test/mocks';
import { BookingPage } from '@/pages/guest/BookingPage';

const mockNavigate = vi.fn();

vi.mock('@mantine/dates', () => ({
  DatePickerInput: ({
    label,
    value,
    onChange,
    placeholder,
  }: {
    label: string;
    value: string | null;
    onChange: (value: string | null) => void;
    placeholder?: string;
  }) => (
    <input
      aria-label={label}
      placeholder={placeholder}
      value={value ?? ''}
      onChange={(event) => onChange(event.currentTarget.value || null)}
    />
  ),
}));

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom');
  return {
    ...actual,
    useParams: () => ({ eventTypeId: 'event-type-1' }),
    useNavigate: () => mockNavigate,
  };
});

describe('BookingPage', () => {
  const findTimeSlotButton = async () => {
    let slotButton: HTMLElement | undefined;

    await waitFor(() => {
      slotButton = screen.getAllByRole('button').find((button) => /\d{2}:\d{2}/.test(button.textContent ?? ''));
      expect(slotButton).toBeDefined();
    });

    return slotButton!;
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders event type details from API', async () => {
    render(<BookingPage />);

    expect(await screen.findByText(/Бронирование: Консультация/i)).toBeInTheDocument();
    expect(screen.getByText(/Индивидуальная консультация по проекту/i)).toBeInTheDocument();
    expect(screen.getAllByText(/30 минут/i).length).toBeGreaterThan(0);
  });

  it('navigates to home when event type cannot be loaded', async () => {
    server.use(
      http.get('*/api/public/event-types/:id', () =>
        HttpResponse.json({ message: 'Тип встречи не найден' }, { status: 404 })
      )
    );

    render(<BookingPage />);

    await waitFor(() => expect(mockNavigate).toHaveBeenCalledWith('/'));
    expect(await screen.findByText(/Тип встречи не найден/i)).toBeInTheDocument();
  });

  it('shows empty state when selected date has no future slots', async () => {
    server.use(
      http.get('*/api/public/slots', () =>
        HttpResponse.json({
          items: [],
          pagination: {
            page: 1,
            pageSize: 10,
            totalItems: 0,
            totalPages: 0,
            hasNext: false,
            hasPrev: false,
          },
        })
      )
    );

    render(<BookingPage />);

    await screen.findByText(/Бронирование: Консультация/i);
    fireEvent.change(screen.getByLabelText(/Дата встречи/i), { target: { value: '2026-05-01' } });

    expect((await screen.findAllByText(/На выбранную дату нет доступных слотов/i)).length).toBeGreaterThan(0);
  });

  it('validates required guest fields before booking', async () => {
    const user = userEvent.setup();
    render(<BookingPage />);

    await screen.findByText(/Бронирование: Консультация/i);
    fireEvent.change(screen.getByLabelText(/Дата встречи/i), { target: { value: '2026-05-01' } });
    await user.click(await findTimeSlotButton());
    await user.click(screen.getByRole('button', { name: /Подтвердить бронирование/i }));

    expect(await screen.findByText(/Пожалуйста, заполните все обязательные поля/i)).toBeInTheDocument();
  });

  it('validates email format before booking', async () => {
    const user = userEvent.setup();
    render(<BookingPage />);

    await screen.findByText(/Бронирование: Консультация/i);
    fireEvent.change(screen.getByLabelText(/Дата встречи/i), { target: { value: '2026-05-01' } });
    await user.click(await findTimeSlotButton());
    await user.type(screen.getByLabelText(/Ваше имя/i), 'Иван Иванов');
    await user.type(screen.getByPlaceholderText('ivan@example.com'), 'wrong-email');
    await user.click(screen.getByRole('button', { name: /Подтвердить бронирование/i }));

    expect(await screen.findByText(/Пожалуйста, введите корректный email/i)).toBeInTheDocument();
  });

  it('creates booking successfully and redirects to home', async () => {
    const user = userEvent.setup();
    render(<BookingPage />);

    await screen.findByText(/Бронирование: Консультация/i);
    fireEvent.change(screen.getByLabelText(/Дата встречи/i), { target: { value: '2026-05-01' } });
    await user.click(await findTimeSlotButton());
    await user.type(screen.getByLabelText(/Ваше имя/i), 'Иван Иванов');
    await user.type(screen.getByPlaceholderText('ivan@example.com'), 'ivan@example.com');
    await user.click(screen.getByRole('button', { name: /Подтвердить бронирование/i }));

    expect(await screen.findByText(/Ваша встреча успешно забронирована/i)).toBeInTheDocument();
    await waitFor(() => expect(mockNavigate).toHaveBeenCalledWith('/'));
  });
});
