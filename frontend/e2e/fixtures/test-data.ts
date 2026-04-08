/**
 * Test fixtures and mock data for E2E tests
 * Note: When using real backend, update these to match actual data
 */

export const testEventTypes = [
  {
    id: 'fc856f83-9dee-436a-bf46-8b8746b664d6',
    name: 'Встреча на 15 минут',
    description: 'Четверть часовая встреча',
    durationMinutes: 15,
  },
  {
    id: '3fdba5e8-9273-4017-955f-a22be3d377c6',
    name: 'Встреча на 30 минут',
    description: 'Получасовая встреча',
    durationMinutes: 30,
  },
];

// For mocked tests
export const mockedEventTypes = [
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
