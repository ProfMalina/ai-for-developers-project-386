import { useState, useEffect, useCallback } from 'react';
import {
  Container,
  Title,
  Text,
  Card,
  Group,
  Button,
  Stack,
  LoadingOverlay,
  Pagination,
  Badge,
} from '@mantine/core';
import { Link } from 'react-router-dom';
import { guestApi } from '../../api/client';
import type { EventType } from '../../types/api';

export function GuestHome() {
  const [eventTypes, setEventTypes] = useState<EventType[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  const fetchEventTypes = useCallback(async () => {
    try {
      setLoading(true);
      const response = await guestApi.getPublicEventTypes({ page, pageSize: 10 });
      setEventTypes(response.items);
      setTotalPages(response.pagination.totalPages);
    } catch (error) {
      console.error('Failed to fetch event types:', error);
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    document.title = 'Book a Meeting - Calendar Booking';
    fetchEventTypes();
  }, [fetchEventTypes]);

  return (
    <Container size="lg" py="xl">
      <Title order={1} ta="center" mb="xl">
        Book a Meeting
      </Title>
      <Text c="dimmed" ta="center" mb="xl" size="lg">
        Select an event type and choose a convenient time slot
      </Text>

      <Card withBorder pos="relative">
        <LoadingOverlay visible={loading} />
        <Stack gap="md">
          {eventTypes.length === 0 && !loading ? (
            <Text c="dimmed" ta="center" py="xl">
              No event types available at the moment
            </Text>
          ) : (
            eventTypes.map((eventType) => (
              <Card key={eventType.id} withBorder padding="lg">
                <Group justify="space-between" mb="xs">
                  <Title order={3}>{eventType.name}</Title>
                  <Badge color="blue" size="lg">
                    {eventType.durationMinutes} min
                  </Badge>
                </Group>
                <Text c="dimmed" mb="md" lineClamp={3}>
                  {eventType.description}
                </Text>
                <Group justify="flex-end">
                  <Button
                    component={Link}
                    to={`/book/${eventType.id}`}
                    size="md"
                  >
                    Book Now
                  </Button>
                </Group>
              </Card>
            ))
          )}
        </Stack>
      </Card>

      {totalPages > 1 && (
        <Group justify="center" mt="xl">
          <Pagination value={page} onChange={setPage} total={totalPages} />
        </Group>
      )}
    </Container>
  );
}
