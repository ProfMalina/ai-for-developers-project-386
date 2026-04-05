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
  Box,
} from '@mantine/core';
import { DatePickerInput } from '@mantine/dates';
import { notifications } from '@mantine/notifications';
import dayjs from 'dayjs';
import { guestApi } from '../../api/client';
import { ApiValidationError } from '../../api/client';
import type { EventType, TimeSlot, CreateBookingRequest } from '../../types/api';

export function BookingPage() {
  const { eventTypeId } = useParams<{ eventTypeId: string }>();
  const navigate = useNavigate();

  const [eventType, setEventType] = useState<EventType | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [availableSlots, setAvailableSlots] = useState<TimeSlot[]>([]);
  const [selectedSlot, setSelectedSlot] = useState<TimeSlot | null>(null);
  const [fetchingSlots, setFetchingSlots] = useState(false);

  // Form state
  const [guestName, setGuestName] = useState('');
  const [guestEmail, setGuestEmail] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const fetchEventType = useCallback(async () => {
    try {
      setLoading(true);
      const data = await guestApi.getPublicEventType(eventTypeId!);
      setEventType(data);
      document.title = `Book ${data.name} - Calendar Booking`;
    } catch {
      notifications.show({
        title: 'Error',
        message: 'Event type not found',
        color: 'red',
      });
      navigate('/');
    } finally {
      setLoading(false);
    }
  }, [eventTypeId, navigate]);

  useEffect(() => {
    if (!eventTypeId) return;
    fetchEventType();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [fetchEventType]);

  const fetchAvailableSlots = async (date: Date) => {
    if (!eventTypeId) return;

    try {
      setFetchingSlots(true);
      const dateFrom = dayjs(date).startOf('day').toISOString();
      const dateTo = dayjs(date).endOf('day').toISOString();

      const response = await guestApi.getAvailableSlots(eventTypeId, {
        dateFrom,
        dateTo,
        pageSize: 100,
      });

      setAvailableSlots(response.items);
      if (response.items.length === 0) {
        notifications.show({
          title: 'No Slots',
          message: 'No available slots for this date',
          color: 'yellow',
        });
      }
    } catch (error) {
      console.error('Failed to fetch slots:', error);
      notifications.show({
        title: 'Error',
        message: 'Failed to load available slots',
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
    }
  };

  const handleBooking = async () => {
    if (!selectedSlot || !guestName.trim() || !guestEmail.trim()) {
      notifications.show({
        title: 'Validation Error',
        message: 'Please fill in all required fields',
        color: 'red',
      });
      return;
    }

    // Email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(guestEmail)) {
      notifications.show({
        title: 'Validation Error',
        message: 'Please enter a valid email address',
        color: 'red',
      });
      return;
    }

    try {
      setSubmitting(true);
      const bookingData: CreateBookingRequest = {
        eventTypeId: eventTypeId!,
        startTime: selectedSlot.startTime,
        guestName: guestName.trim(),
        guestEmail: guestEmail.trim(),
      };

      await guestApi.createBooking(bookingData);

      notifications.show({
        title: 'Success!',
        message: 'Your booking has been created successfully',
        color: 'green',
      });

      navigate('/');
    } catch (error) {
      console.error('Failed to create booking:', error);

      if (error instanceof ApiValidationError && error.data?.fieldErrors) {
        const messages = error.data.fieldErrors.map((e) => e.message).join(', ');
        notifications.show({
          title: 'Validation Error',
          message: messages,
          color: 'red',
        });
      } else if (error instanceof Error) {
        notifications.show({
          title: 'Error',
          message: error.message,
          color: 'red',
        });
      } else {
        notifications.show({
          title: 'Error',
          message: 'Failed to create booking',
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
      <Title order={1} mb="md">
        Book: {eventType.name}
      </Title>

      <Card withBorder mb="xl" padding="lg">
        <Stack gap="md">
          <Group>
            <Title order={3}>{eventType.name}</Title>
            <Badge color="blue" size="lg">
              {eventType.durationMinutes} min
            </Badge>
          </Group>
          <Text>{eventType.description}</Text>
        </Stack>
      </Card>

      <Card withBorder mb="xl" padding="lg">
        <Title order={4} mb="md">
          1. Select Date
        </Title>
        <DatePickerInput
          value={selectedDate}
          onChange={handleDateChange}
          placeholder="Pick a date"
          label="Date"
          minDate={dayjs().add(1, 'day').format('YYYY-MM-DD')}
          size="md"
        />
      </Card>

      {fetchingSlots ? (
        <Card withBorder mb="xl" padding="lg">
          <LoadingOverlay visible={fetchingSlots} />
          <Title order={4} mb="md">
            2. Select Time Slot
          </Title>
        </Card>
      ) : availableSlots.length > 0 ? (
        <Card withBorder mb="xl" padding="lg">
          <Title order={4} mb="md">
            2. Select Time Slot
          </Title>
          <Group gap="xs" wrap="wrap">
            {availableSlots.map((slot) => (
              <Button
                key={slot.id}
                variant={selectedSlot?.id === slot.id ? 'filled' : 'outline'}
                onClick={() => setSelectedSlot(slot)}
                size="md"
              >
                {dayjs(slot.startTime).format('HH:mm')}
              </Button>
            ))}
          </Group>
        </Card>
      ) : selectedDate ? (
        <Card withBorder mb="xl" padding="lg">
          <Title order={4} mb="md">
            2. Select Time Slot
          </Title>
          <Text c="dimmed" ta="center" py="xl">
            No available slots for this date
          </Text>
        </Card>
      ) : null}

      {selectedSlot && (
        <Card withBorder mb="xl" padding="lg">
          <Title order={4} mb="md">
            3. Your Details
          </Title>
          <Stack gap="md">
            <TextInput
              label="Your Name"
              placeholder="John Doe"
              value={guestName}
              onChange={(e) => setGuestName(e.target.value)}
              required
              size="md"
            />
            <TextInput
              label="Email"
              placeholder="john@example.com"
              value={guestEmail}
              onChange={(e) => setGuestEmail(e.target.value)}
              required
              size="md"
            />
            <Divider />
            <Box>
              <Text fw={500} mb="xs">
                Booking Summary:
              </Text>
              <Stack gap="xs">
                <Text>
                  <Text component="span" fw={500}>Event:</Text> {eventType.name}
                </Text>
                <Text>
                  <Text component="span" fw={500}>Date:</Text>{' '}
                  {dayjs(selectedSlot.startTime).format('MMMM DD, YYYY')}
                </Text>
                <Text>
                  <Text component="span" fw={500}>Time:</Text>{' '}
                  {dayjs(selectedSlot.startTime).format('HH:mm')} -{' '}
                  {dayjs(selectedSlot.endTime).format('HH:mm')}
                </Text>
                <Text>
                  <Text component="span" fw={500}>Duration:</Text> {eventType.durationMinutes} minutes
                </Text>
              </Stack>
            </Box>
            <Button
              onClick={handleBooking}
              loading={submitting}
              size="lg"
              fullWidth
            >
              Confirm Booking
            </Button>
          </Stack>
        </Card>
      )}
    </Container>
  );
}
