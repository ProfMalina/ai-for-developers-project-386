import { useState, useEffect } from 'react';
import {
  Card,
  Button,
  Group,
  Text,
  Modal,
  Select,
  TextInput,
  NumberInput,
  Stack,
  Checkbox,
  LoadingOverlay,
  Alert,
} from '@mantine/core';
import { DatePickerInput } from '@mantine/dates';
import dayjs from 'dayjs';
import { notifications } from '@mantine/notifications';
import { IconClock, IconAlertCircle } from '@tabler/icons-react';
import { ownerApi } from '../../api/client';
import type { EventType } from '../../types/api';

const DAYS = [
  { value: '1', label: 'Пн' },
  { value: '2', label: 'Вт' },
  { value: '3', label: 'Ср' },
  { value: '4', label: 'Чт' },
  { value: '5', label: 'Пт' },
  { value: '6', label: 'Сб' },
  { value: '7', label: 'Вс' },
];

export function SlotGeneration() {
  const [eventTypes, setEventTypes] = useState<EventType[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpened, setModalOpened] = useState(false);

  const [selectedEventTypeId, setSelectedEventTypeId] = useState<string | null>(null);
  const [workingHoursStart, setWorkingHoursStart] = useState('09:00');
  const [workingHoursEnd, setWorkingHoursEnd] = useState('18:00');
  const [intervalMinutes, setIntervalMinutes] = useState(30);
  const [daysOfWeek, setDaysOfWeek] = useState<string[]>(['1', '2', '3', '4', '5']);
  const [dateFrom, setDateFrom] = useState<string | null>(dayjs().format('YYYY-MM-DD'));
  const [dateTo, setDateTo] = useState<string | null>(dayjs().add(30, 'day').format('YYYY-MM-DD'));

  useEffect(() => {
    loadEventTypes();
  }, []);

  const loadEventTypes = async () => {
    try {
      const result = await ownerApi.getEventTypes({ pageSize: 100 });
      setEventTypes(result.items);
      if (result.items.length > 0) {
        setSelectedEventTypeId(result.items[0].id);
      }
    } catch (error) {
      notifications.show({
        title: 'Ошибка',
        message: 'Не удалось загрузить типы встреч',
        color: 'red',
      });
    }
  };

  const handleSubmit = async () => {
    if (!selectedEventTypeId) {
      notifications.show({
        title: 'Ошибка',
        message: 'Выберите тип встречи',
        color: 'red',
      });
      return;
    }

    if (daysOfWeek.length === 0) {
      notifications.show({
        title: 'Ошибка',
        message: 'Выберите хотя бы один день недели',
        color: 'red',
      });
      return;
    }

    if (!dateFrom || !dateTo) {
      notifications.show({
        title: 'Ошибка',
        message: 'Укажите диапазон дат',
        color: 'red',
      });
      return;
    }

    setLoading(true);
    try {
      const result = await ownerApi.generateSlots(selectedEventTypeId, {
        workingHoursStart,
        workingHoursEnd,
        intervalMinutes,
        daysOfWeek: daysOfWeek.map(Number),
        dateFrom: dateFrom || undefined,
        dateTo: dateTo || undefined,
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
      });

      notifications.show({
        title: 'Успех',
        message: `Создано ${result.slotsCreated} слотов`,
        color: 'green',
      });

      setModalOpened(false);
    } catch (error: any) {
      notifications.show({
        title: 'Ошибка',
        message: error.message || 'Не удалось создать слоты',
        color: 'red',
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Card withBorder shadow="sm" radius="md">
        <Group justify="space-between" mb="md">
          <div>
            <Text fw={500} size="lg">
              Генерация расписания
            </Text>
            <Text size="sm" c="dimmed">
              Создайте доступные слоты для бронирования
            </Text>
          </div>
          <Button
            leftSection={<IconClock size={16} />}
            onClick={() => setModalOpened(true)}
            disabled={eventTypes.length === 0}
          >
            Создать слоты
          </Button>
        </Group>

        {eventTypes.length === 0 && (
          <Alert icon={<IconAlertCircle size={16} />} title="Нет типов встреч" color="yellow">
            Сначала создайте хотя бы один тип встречи
          </Alert>
        )}
      </Card>

      <Modal
        opened={modalOpened}
        onClose={() => setModalOpened(false)}
        title="Генерация слотов"
        size="md"
      >
        <LoadingOverlay visible={loading} />
        <Stack>
          <Select
            label="Тип встречи"
            placeholder="Выберите тип встречи"
            value={selectedEventTypeId}
            onChange={setSelectedEventTypeId}
            data={eventTypes.map((et) => ({ value: et.id, label: et.name }))}
            required
          />

          <Group grow>
            <TextInput
              label="Начало рабочего дня"
              value={workingHoursStart}
              onChange={(e) => setWorkingHoursStart(e.currentTarget.value)}
              placeholder="09:00"
              required
            />
            <TextInput
              label="Конец рабочего дня"
              value={workingHoursEnd}
              onChange={(e) => setWorkingHoursEnd(e.currentTarget.value)}
              placeholder="18:00"
              required
            />
          </Group>

          <NumberInput
            label="Длительность слота (минуты)"
            value={intervalMinutes}
            onChange={(val) => setIntervalMinutes(Number(val))}
            min={15}
            max={120}
            step={15}
            required
          />

          <Checkbox.Group
            label="Дни недели"
            value={daysOfWeek}
            onChange={setDaysOfWeek}
            required
          >
            <Group mt="xs">
              {DAYS.map((day) => (
                <Checkbox key={day.value} value={day.value} label={day.label} />
              ))}
            </Group>
          </Checkbox.Group>

          <Group grow>
            <DatePickerInput
              label="Дата начала"
              value={dateFrom}
              onChange={setDateFrom}
              placeholder="ДД.ММ.ГГГГ"
              minDate={dayjs().format('YYYY-MM-DD')}
              required
            />
            <DatePickerInput
              label="Дата окончания"
              value={dateTo}
              onChange={setDateTo}
              placeholder="ДД.ММ.ГГГГ"
              minDate={dateFrom || dayjs().format('YYYY-MM-DD')}
              required
            />
          </Group>

          <Button onClick={handleSubmit} fullWidth loading={loading}>
            Сгенерировать слоты
          </Button>
        </Stack>
      </Modal>
    </>
  );
}
