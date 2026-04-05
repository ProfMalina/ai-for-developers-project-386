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
import { ownerApi } from '../../api/client';
import type { Booking } from '../../types/api';

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
    } catch (error) {
      console.error('Failed to fetch bookings:', error);
      notifications.show({
        title: 'Error',
        message: 'Failed to load bookings',
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
        title: 'Success',
        message: 'Booking cancelled successfully',
        color: 'green',
      });
      setCancelModal({ opened: false, bookingId: null });
      fetchBookings();
    } catch (error) {
      console.error('Failed to cancel booking:', error);
      notifications.show({
        title: 'Error',
        message: 'Failed to cancel booking',
        color: 'red',
      });
    }
  };

  if (bookings.length === 0 && !loading) {
    return (
      <Card withBorder>
        <Text c="dimmed" ta="center" py="xl">
          No bookings yet
        </Text>
      </Card>
    );
  }

  return (
    <div>
      <Text size="lg" fw={500} mb="md">
        Upcoming Bookings
      </Text>

      <Card withBorder pos="relative">
        <LoadingOverlay visible={loading} />
        <Stack gap="md">
          {bookings.map((booking) => (
            <Card key={booking.id} withBorder padding="md">
              <Group justify="space-between" mb="xs">
                <Text fw={500}>{booking.guestName}</Text>
                <Badge color={dayjs(booking.startTime).isBefore(dayjs()) ? 'gray' : 'green'}>
                  {dayjs(booking.startTime).isBefore(dayjs()) ? 'Past' : 'Upcoming'}
                </Badge>
              </Group>
              <Stack gap="xs">
                <Text size="sm">
                  <Text component="span" fw={500}>Email:</Text> {booking.guestEmail}
                </Text>
                <Text size="sm">
                  <Text component="span" fw={500}>Time:</Text>{' '}
                  {dayjs(booking.startTime).format('MMMM DD, YYYY HH:mm')} -{' '}
                  {dayjs(booking.endTime).format('HH:mm')}
                </Text>
                <Text size="sm">
                  <Text component="span" fw={500}>Event Type ID:</Text> {booking.eventTypeId}
                </Text>
                <Text size="xs" c="dimmed">
                  Created: {dayjs(booking.createdAt).format('MMMM DD, YYYY HH:mm')}
                </Text>
              </Stack>
              <Group justify="flex-end" mt="md">
                <Button
                  color="red"
                  variant="light"
                  onClick={() => setCancelModal({ opened: true, bookingId: booking.id })}
                  disabled={dayjs(booking.startTime).isBefore(dayjs())}
                >
                  Cancel Booking
                </Button>
              </Group>
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
        title="Cancel Booking"
      >
        <Text mb="md">Are you sure you want to cancel this booking?</Text>
        <Group justify="flex-end">
          <Button
            variant="default"
            onClick={() => setCancelModal({ opened: false, bookingId: null })}
          >
            No
          </Button>
          <Button color="red" onClick={handleCancelBooking}>
            Yes, Cancel
          </Button>
        </Group>
      </Modal>
    </div>
  );
}
