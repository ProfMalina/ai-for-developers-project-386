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
  ThemeIcon,
  SimpleGrid,
} from '@mantine/core';
import { IconClock, IconCalendarEvent } from '@tabler/icons-react';
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
    } catch {
      // Ошибки обрабатываются на уровне API клиента
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    document.title = 'Забронировать встречу';
    fetchEventTypes();
  }, [fetchEventTypes]);

  return (
    <Container size="lg" py="xl">
      <Stack gap="xl" align="center">
        <Stack gap="xs" align="center">
          <Title order={1} size={42} ta="center">
            Забронировать встречу
          </Title>
          <Text c="dimmed" size="lg" ta="center">
            Выберите тип встречи и удобное для вас время
          </Text>
        </Stack>

        <Card withBorder radius="md" p="xl" pos="relative" w="100%">
          <LoadingOverlay visible={loading} />
          <Stack gap="md">
            {eventTypes.length === 0 && !loading ? (
              <Text c="dimmed" ta="center" py="xl">
                В данный момент нет доступных типов встреч
              </Text>
            ) : (
              <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="md">
                {eventTypes.map((eventType) => (
                  <Card
                    key={eventType.id}
                    withBorder
                    padding="lg"
                    radius="md"
                    style={{ transition: 'box-shadow 0.2s' }}
                  >
                    <Stack gap="sm">
                      <Group justify="space-between" wrap="nowrap">
                        <ThemeIcon size={40} radius="md" variant="light" color="blue">
                          <IconCalendarEvent size={20} />
                        </ThemeIcon>
                        <Badge
                          size="lg"
                          variant="light"
                          leftSection={<IconClock size={12} />}
                        >
                          {eventType.durationMinutes} мин
                        </Badge>
                      </Group>

                      <Title order={3} size={20}>
                        {eventType.name}
                      </Title>

                      <Text size="sm" c="dimmed" lineClamp={3} style={{ minHeight: 60 }}>
                        {eventType.description}
                      </Text>

                      <Button
                        component={Link}
                        to={`/book/${eventType.id}`}
                        fullWidth
                        mt="md"
                        size="md"
                      >
                        Забронировать
                      </Button>
                    </Stack>
                  </Card>
                ))}
              </SimpleGrid>
            )}
          </Stack>
        </Card>

        {totalPages > 1 && (
          <Group justify="center" mt="xl">
            <Pagination
              value={page}
              onChange={setPage}
              total={totalPages}
              size="md"
            />
          </Group>
        )}
      </Stack>
    </Container>
  );
}
