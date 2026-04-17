import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import userEvent from '@testing-library/user-event';
import { http, HttpResponse } from 'msw';
import { fireEvent } from '@testing-library/react';
import { render, screen } from '@/test/test-utils';
import { server } from '@/test/mocks';
import { SlotGeneration } from '@/components/owner/SlotGeneration';

vi.mock('@mantine/core', async () => {
  const actual = await vi.importActual<typeof import('@mantine/core')>('@mantine/core');
  return {
    ...actual,
    Modal: ({ opened, title, children }: { opened: boolean; title: string; children: ReactNode }) => (
      opened ? <div role="dialog" aria-label={title}>{children}</div> : null
    ),
  };
});

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

describe('SlotGeneration', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('opens modal with generation form', async () => {
    const user = userEvent.setup();
    render(<SlotGeneration />);

    await user.click(screen.getByRole('button', { name: /Создать слоты/i }));

    expect(await screen.findByRole('dialog')).toBeInTheDocument();
    expect(screen.getByLabelText(/Начало рабочего дня/i)).toHaveValue('09:00');
    expect((screen.getByLabelText(/Дата начала/i) as HTMLInputElement).value).not.toBe('');
  });

  it('requires at least one weekday', async () => {
    const user = userEvent.setup();
    render(<SlotGeneration />);

    await user.click(screen.getByRole('button', { name: /Создать слоты/i }));
    await user.click(screen.getByRole('checkbox', { name: 'Пн' }));
    await user.click(screen.getByRole('checkbox', { name: 'Вт' }));
    await user.click(screen.getByRole('checkbox', { name: 'Ср' }));
    await user.click(screen.getByRole('checkbox', { name: 'Чт' }));
    await user.click(screen.getByRole('checkbox', { name: 'Пт' }));
    await user.click(screen.getByRole('button', { name: /Сгенерировать слоты/i }));

    expect(await screen.findByText(/Выберите хотя бы один день недели/i)).toBeInTheDocument();
  });

  it('requires a date range before submission', async () => {
    const user = userEvent.setup();
    render(<SlotGeneration />);

    await user.click(screen.getByRole('button', { name: /Создать слоты/i }));
    fireEvent.change(screen.getByLabelText(/Дата начала/i), { target: { value: '' } });
    fireEvent.change(screen.getByLabelText(/Дата окончания/i), { target: { value: '' } });
    await user.click(screen.getByRole('button', { name: /Сгенерировать слоты/i }));

    expect(await screen.findByText(/Укажите диапазон дат/i)).toBeInTheDocument();
  });

  it('shows success message after generating slots', async () => {
    const user = userEvent.setup();
    render(<SlotGeneration />);

    await user.click(screen.getByRole('button', { name: /Создать слоты/i }));
    await user.click(screen.getByRole('button', { name: /Сгенерировать слоты/i }));

    expect(await screen.findByText(/Создано 10 слотов/i)).toBeInTheDocument();
  });

  it('shows API error when generation fails', async () => {
    server.use(
      http.post('*/api/slots/generate', () =>
        HttpResponse.json({ message: 'Генерация недоступна' }, { status: 500 })
      )
    );

    const user = userEvent.setup();
    render(<SlotGeneration />);

    await user.click(screen.getByRole('button', { name: /Создать слоты/i }));
    await user.click(screen.getByRole('button', { name: /Сгенерировать слоты/i }));

    expect(await screen.findByText(/Генерация недоступна/i)).toBeInTheDocument();
  });
});
