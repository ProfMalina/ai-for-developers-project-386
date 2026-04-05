import { useState, useEffect, useCallback } from 'react';
import {
  Card,
  Text,
  Group,
  Badge,
  Stack,
  LoadingOverlay,
  Pagination,
  Button,
  Modal,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import dayjs from 'dayjs';
import 'dayjs/locale/ru';
import { ownerApi } from '../../api/client';
import type { Booking } from '../../types/api';

dayjs.locale('ru');

export function BookingsList() {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [cancelModal, setCancelModal] = useState<{ opened: boolean; bookingId: string | null }>({
    opened: false,
    bookingId: null,
  });

  const fetchBookings = useCallback(async () => {
    try {
      setLoading(true);
      const response = await ownerApi.getAllBookings({
        page,
        pageSize: 10,
        sortBy: 'startTime',
        sortOrder: 'asc',
      });
      setBookings(response.items);
      setTotalPages(response.pagination.totalPages);
    } catch {
      notifications.show({
        title: 'Ошибка',
        message: 'Не удалось загрузить бронирования',
        color: 'red',
      });
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    fetchBookings();
  }, [fetchBookings]);

  const handleCancelBooking = async () => {
    if (!cancelModal.bookingId) return;

    try {
      await ownerApi.cancelBooking(cancelModal.bookingId);
      notifications.show({
        title: 'Успешно',
        message: 'Бронирование отменено',
        color: 'green',
      });
      setCancelModal({ opened: false, bookingId: null });
      fetchBookings();
    } catch {
      notifications.show({
        title: 'Ошибка',
        message: 'Не удалось отменить бронирование',
        color: 'red',
      });
    }
  };

  if (bookings.length === 0 && !loading) {
    return (
      <Card withBorder radius="md">
        <Text c="dimmed" ta="center" py="xl">
          Бронирований пока нет
        </Text>
      </Card>
    );
  }

  return (
    <div>
      <Text size="lg" fw={500} mb="md">
        Список бронирований
      </Text>

      <Card withBorder pos="relative" radius="md">
        <LoadingOverlay visible={loading} />
        <Stack gap="md">
          {bookings.map((booking) => (
            <Card key={booking.id} withBorder padding="md" radius="md">
              <Stack gap="xs">
                <Group justify="space-between">
                  <Text fw={600} size="lg">
                    {booking.guestName}
                  </Text>
                  <Badge
                    color={dayjs(booking.startTime).isBefore(dayjs()) ? 'gray' : 'green'}
                    variant="light"
                  >
                    {dayjs(booking.startTime).isBefore(dayjs()) ? 'Прошедшая' : 'Предстоящая'}
                  </Badge>
                </Group>

                <Stack gap="xs">
                  <Group justify="space-between">
                    <Text size="sm" c="dimmed">Email:</Text>
                    <Text size="sm">{booking.guestEmail}</Text>
                  </Group>
                  <Group justify="space-between">
                    <Text size="sm" c="dimmed">Время:</Text>
                    <Text size="sm">
                      {dayjs(booking.startTime).locale('ru').format('D MMMM YYYY, HH:mm')} -{' '}
                      {dayjs(booking.endTime).format('HH:mm')}
                    </Text>
                  </Group>
                  <Group justify="space-between">
                    <Text size="sm" c="dimmed">ID типа встречи:</Text>
                    <Text size="sm" ff="monospace">{booking.eventTypeId}</Text>
                  </Group>
                  <Text size="xs" c="dimmed">
                    Создано: {dayjs(booking.createdAt).locale('ru').format('D MMMM YYYY, HH:mm')}
                  </Text>
                </Stack>

                <Group justify="flex-end" mt="xs">
                  <Button
                    color="red"
                    variant="light"
                    onClick={() => setCancelModal({ opened: true, bookingId: booking.id })}
                    disabled={dayjs(booking.startTime).isBefore(dayjs())}
                  >
                    Отменить бронирование
                  </Button>
                </Group>
              </Stack>
            </Card>
          ))}
        </Stack>
      </Card>

      {totalPages > 1 && (
        <Group justify="center" mt="md">
          <Pagination value={page} onChange={setPage} total={totalPages} />
        </Group>
      )}

      <Modal
        opened={cancelModal.opened}
        onClose={() => setCancelModal({ opened: false, bookingId: null })}
        title="Отмена бронирования"
      >
        <Text mb="md">Вы уверены, что хотите отменить это бронирование?</Text>
        <Group justify="flex-end">
          <Button
            variant="default"
            onClick={() => setCancelModal({ opened: false, bookingId: null })}
          >
            Нет
          </Button>
          <Button color="red" onClick={handleCancelBooking}>
            Да, отменить
          </Button>
        </Group>
      </Modal>
    </div>
  );
}
