/**
 * Test fixtures and mock data for E2E tests
 */

export const testEventTypes = [
  {
    id: 'test-consultation',
    name: 'Консультация',
    description: 'Индивидуальная консультация по вашему вопросу',
    durationMinutes: 30,
  },
  {
    id: 'test-meeting',
    name: 'Встреча',
    description: 'Групповая встреча для обсуждения проекта',
    durationMinutes: 60,
  },
];

export const testBooking = {
  guestName: 'Тестовый Пользователь',
  guestEmail: 'test@example.com',
};

export const testEventTypeNew = {
  name: 'Мастер-класс',
  description: 'Практическое занятие для группы',
  durationMinutes: 90,
};

export const testEventTypeUpdated = {
  name: 'Мастер-класс PRO',
  description: 'Продвинутое практическое занятие',
  durationMinutes: 120,
};
