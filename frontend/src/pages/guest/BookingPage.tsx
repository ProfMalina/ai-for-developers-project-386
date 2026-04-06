import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Container,
  Title,
  Text,
  Card,
  Button,
  Stack,
  Group,
  TextInput,
  LoadingOverlay,
  Badge,
  Divider,
  ThemeIcon,
  Stepper,
} from '@mantine/core';
import { DatePickerInput } from '@mantine/dates';
import { notifications } from '@mantine/notifications';
import { IconClock, IconCalendarEvent, IconUser, IconMail, IconCheck } from '@tabler/icons-react';
import dayjs from 'dayjs';
import 'dayjs/locale/ru';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import { guestApi } from '../../api/client';
import { ApiValidationError, ApiError } from '../../api/client';
import type { EventType, TimeSlot, CreateBookingRequest } from '../../types/api';

dayjs.extend(utc);
dayjs.extend(timezone);
dayjs.locale('ru');

export function BookingPage() {
  const { eventTypeId } = useParams<{ eventTypeId: string }>();
  const navigate = useNavigate();

  const [eventType, setEventType] = useState<EventType | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [availableSlots, setAvailableSlots] = useState<TimeSlot[]>([]);
  const [selectedSlot, setSelectedSlot] = useState<TimeSlot | null>(null);
  const [fetchingSlots, setFetchingSlots] = useState(false);
  const [activeStep, setActiveStep] = useState(0);

  // Form state
  const [guestName, setGuestName] = useState('');
  const [guestEmail, setGuestEmail] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const fetchEventType = useCallback(async () => {
    if (!eventTypeId) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      const data = await guestApi.getPublicEventType(eventTypeId);
      setEventType(data);
      document.title = `Бронирование: ${data.name}`;
    } catch {
      notifications.show({
        title: 'Ошибка',
        message: 'Тип встречи не найден',
        color: 'red',
      });
      navigate('/');
    } finally {
      setLoading(false);
    }
  }, [eventTypeId, navigate]);

  useEffect(() => {
    fetchEventType();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [fetchEventType]);

  const fetchAvailableSlots = async (date: Date) => {
    try {
      setFetchingSlots(true);
      const dateFrom = dayjs(date).startOf('day').toISOString();
      const dateTo = dayjs(date).endOf('day').toISOString();

      const response = await guestApi.getAvailableSlots({
        dateFrom,
        dateTo,
        pageSize: 100,
      });

      setAvailableSlots(response.items);
      if (response.items.length === 0) {
        notifications.show({
          title: 'Нет слотов',
          message: 'На выбранную дату нет доступных слотов',
          color: 'yellow',
        });
      }
    } catch {
      notifications.show({
        title: 'Ошибка',
        message: 'Не удалось загрузить доступные слоты',
        color: 'red',
      });
      setAvailableSlots([]);
    } finally {
      setFetchingSlots(false);
    }
  };

  const handleDateChange = (date: string | null) => {
    setSelectedDate(date);
    setSelectedSlot(null);
    if (date) {
      const dateObj = new Date(date);
      fetchAvailableSlots(dateObj);
      setActiveStep(1);
    }
  };

  const handleSlotSelect = (slot: TimeSlot) => {
    // Calculate end time based on event type duration
    if (eventType) {
      const start = dayjs(slot.startTime);
      const end = start.add(eventType.durationMinutes, 'minute');
      setSelectedSlot({
        ...slot,
        endTime: end.toISOString(),
      });
    } else {
      setSelectedSlot(slot);
    }
    setActiveStep(2);
  };

  const handleBooking = async () => {
    if (!selectedSlot || !guestName.trim() || !guestEmail.trim()) {
      notifications.show({
        title: 'Ошибка валидации',
        message: 'Пожалуйста, заполните все обязательные поля',
        color: 'red',
      });
      return;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(guestEmail)) {
      notifications.show({
        title: 'Ошибка валидации',
        message: 'Пожалуйста, введите корректный email',
        color: 'red',
      });
      return;
    }

    try {
      setSubmitting(true);
      const bookingData: CreateBookingRequest = {
        eventTypeId: eventTypeId!,
        slotId: selectedSlot.id,
        startTime: selectedSlot.startTime,
        guestName: guestName.trim(),
        guestEmail: guestEmail.trim(),
      };

      await guestApi.createBooking(bookingData);

      notifications.show({
        title: 'Успешно!',
        message: 'Ваша встреча успешно забронирована',
        color: 'green',
        icon: <IconCheck />,
      });

      navigate('/');
    } catch (error) {
      if (error instanceof ApiValidationError && error.data?.fieldErrors) {
        const messages = error.data.fieldErrors.map((e) => e.message).join(', ');
        notifications.show({
          title: 'Ошибка валидации',
          message: messages,
          color: 'red',
        });
      } else if (error instanceof ApiError && error.message) {
        notifications.show({
          title: 'Ошибка',
          message: error.message,
          color: 'red',
        });
      } else {
        notifications.show({
          title: 'Ошибка',
          message: 'Не удалось создать бронирование',
          color: 'red',
        });
      }
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <Container size="md" py="xl">
        <LoadingOverlay visible={loading} />
      </Container>
    );
  }

  if (!eventType) {
    return null;
  }

  return (
    <Container size="md" py="xl">
      <Stack gap="xl">
        <Stack gap="xs">
          <Title order={1} size={36}>
            Бронирование: {eventType.name}
          </Title>
          <Text c="dimmed" size="lg">
            {eventType.description}
          </Text>
          <Badge size="lg" leftSection={<IconClock size={12} />}>
            {eventType.durationMinutes} минут
          </Badge>
        </Stack>

        <Stepper active={activeStep} onStepClick={setActiveStep}>
          <Stepper.Step label="Шаг 1" description="Выберите дату" />
          <Stepper.Step label="Шаг 2" description="Выберите время" />
          <Stepper.Step label="Шаг 3" description="Ваши данные" />
        </Stepper>

        {/* Step 1: Date Selection */}
        <Card withBorder radius="md" p="xl">
          <Title order={4} mb="md">
            1. Выберите дату
          </Title>
          <DatePickerInput
            value={selectedDate}
            onChange={handleDateChange}
            placeholder="Выберите дату"
            label="Дата встречи"
            description="Доступны только будущие даты"
            minDate={dayjs().format('YYYY-MM-DD')}
            size="md"
            locale="ru"
          />
        </Card>

        {/* Step 2: Time Slot Selection */}
        {fetchingSlots ? (
          <Card withBorder radius="md" p="xl">
            <LoadingOverlay visible={fetchingSlots} />
            <Title order={4} mb="md">
              2. Выберите время
            </Title>
          </Card>
        ) : availableSlots.length > 0 ? (
          <Card withBorder radius="md" p="xl">
            <Title order={4} mb="md">
              2. Выберите время
            </Title>
            <Text size="sm" c="dimmed" mb="md">
              Доступные слоты на {dayjs(selectedDate).locale('ru').format('D MMMM YYYY')}
            </Text>
            <Group gap="xs" wrap="wrap">
              {availableSlots.map((slot) => (
                <Button
                  key={slot.id}
                  variant={selectedSlot?.id === slot.id ? 'filled' : 'outline'}
                  onClick={() => handleSlotSelect(slot)}
                  size="md"
                  leftSection={<IconClock size={16} />}
                >
                  {dayjs.utc(slot.startTime).local().format('HH:mm')}
                </Button>
              ))}
            </Group>
          </Card>
        ) : selectedDate ? (
          <Card withBorder radius="md" p="xl">
            <Title order={4} mb="md">
              2. Выберите время
            </Title>
            <Text c="dimmed" ta="center" py="xl">
              На выбранную дату нет доступных слотов
            </Text>
          </Card>
        ) : null}

        {/* Step 3: Guest Details */}
        {selectedSlot && (
          <Card withBorder radius="md" p="xl">
            <Title order={4} mb="md">
              3. Ваши данные
            </Title>
            <Stack gap="md">
              <TextInput
                label="Ваше имя"
                placeholder="Иван Иванов"
                value={guestName}
                onChange={(e) => setGuestName(e.target.value)}
                required
                size="md"
                leftSection={<ThemeIcon size={20} radius="xl" variant="light"><IconUser size={14} /></ThemeIcon>}
              />
              <TextInput
                label="Email"
                placeholder="ivan@example.com"
                value={guestEmail}
                onChange={(e) => setGuestEmail(e.target.value)}
                required
                size="md"
                leftSection={<ThemeIcon size={20} radius="xl" variant="light"><IconMail size={14} /></ThemeIcon>}
              />
              <Divider my="md" />
              <Card withBorder p="md" radius="md" bg="var(--mantine-color-gray-0)">
                <Title order={5} mb="sm">
                  Детали бронирования:
                </Title>
                <Stack gap="xs">
                  <Group justify="space-between">
                    <Text c="dimmed">Тип встречи:</Text>
                    <Text fw={500}>{eventType.name}</Text>
                  </Group>
                  <Group justify="space-between">
                    <Text c="dimmed">Дата:</Text>
                    <Text fw={500}>
                      {dayjs(selectedSlot.startTime).locale('ru').format('D MMMM YYYY')}
                    </Text>
                  </Group>
                  <Group justify="space-between">
                    <Text c="dimmed">Время:</Text>
                    <Text fw={500}>
                      {dayjs(selectedSlot.startTime).format('HH:mm')} -{' '}
                      {dayjs(selectedSlot.endTime).format('HH:mm')}
                    </Text>
                  </Group>
                  <Group justify="space-between">
                    <Text c="dimmed">Длительность:</Text>
                    <Text fw={500}>{eventType.durationMinutes} минут</Text>
                  </Group>
                </Stack>
              </Card>
              <Button
                onClick={handleBooking}
                loading={submitting}
                size="lg"
                fullWidth
                leftSection={<IconCalendarEvent size={18} />}
              >
                Подтвердить бронирование
              </Button>
            </Stack>
          </Card>
        )}
      </Stack>
    </Container>
  );
}
