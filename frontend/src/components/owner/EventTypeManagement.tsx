import { useState, useEffect, useCallback } from 'react';
import {
  Card,
  Button,
  Group,
  Text,
  Modal,
  TextInput,
  Textarea,
  NumberInput,
  Stack,
  Badge,
  ActionIcon,
  LoadingOverlay,
  Pagination,
  SimpleGrid,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { IconPlus, IconEdit, IconTrash, IconClock } from '@tabler/icons-react';
import { ownerApi } from '../../api/client';
import type { EventType, CreateEventTypeRequest, UpdateEventTypeRequest } from '../../types/api';

export function EventTypeManagement() {
  const [eventTypes, setEventTypes] = useState<EventType[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalOpened, setModalOpened] = useState(false);
  const [editingEvent, setEditingEvent] = useState<EventType | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  // Form state
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [durationMinutes, setDurationMinutes] = useState<number | string>(30);

  const fetchEventTypes = useCallback(async () => {
    try {
      setLoading(true);
      const response = await ownerApi.getEventTypes({ page, pageSize: 10 });
      setEventTypes(response.items);
      setTotalPages(response.pagination.totalPages);
    } catch {
      notifications.show({
        title: 'Ошибка',
        message: 'Не удалось загрузить типы встреч',
        color: 'red',
      });
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    fetchEventTypes();
  }, [fetchEventTypes]);

  const openCreateModal = () => {
    setEditingEvent(null);
    setName('');
    setDescription('');
    setDurationMinutes(30);
    setModalOpened(true);
  };

  const openEditModal = (eventType: EventType) => {
    setEditingEvent(eventType);
    setName(eventType.name);
    setDescription(eventType.description);
    setDurationMinutes(eventType.durationMinutes);
    setModalOpened(true);
  };

  const handleSubmit = async () => {
    if (!name.trim() || !description.trim() || !durationMinutes) {
      notifications.show({
        title: 'Ошибка валидации',
        message: 'Все поля обязательны для заполнения',
        color: 'red',
      });
      return;
    }

    if (typeof durationMinutes === 'number' && durationMinutes < 5) {
      notifications.show({
        title: 'Ошибка валидации',
        message: 'Длительность должна быть не менее 5 минут',
        color: 'red',
      });
      return;
    }

    try {
      if (editingEvent) {
        const updateData: UpdateEventTypeRequest = {
          name,
          description,
          durationMinutes: durationMinutes as number,
        };
        await ownerApi.updateEventType(editingEvent.id, updateData);
        notifications.show({
          title: 'Успешно',
          message: 'Тип встречи обновлен',
          color: 'green',
        });
      } else {
        const createData: CreateEventTypeRequest = {
          name,
          description,
          durationMinutes: durationMinutes as number,
        };
        await ownerApi.createEventType(createData);
        notifications.show({
          title: 'Успешно',
          message: 'Тип встречи создан',
          color: 'green',
        });
      }
      setModalOpened(false);
      fetchEventTypes();
    } catch {
      notifications.show({
        title: 'Ошибка',
        message: 'Не удалось сохранить тип встречи',
        color: 'red',
      });
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот тип встречи?')) {
      return;
    }

    try {
      await ownerApi.deleteEventType(id);
      notifications.show({
        title: 'Успешно',
        message: 'Тип встречи удален',
        color: 'green',
      });
      fetchEventTypes();
    } catch {
      notifications.show({
        title: 'Ошибка',
        message: 'Не удалось удалить тип встречи',
        color: 'red',
      });
    }
  };

  return (
    <div>
      <Group justify="space-between" mb="md">
        <Text size="lg" fw={500}>Управление типами встреч</Text>
        <Button leftSection={<IconPlus size={16} />} onClick={openCreateModal}>
          Добавить тип встречи
        </Button>
      </Group>

      <Card withBorder pos="relative" radius="md">
        <LoadingOverlay visible={loading} />
        <Stack gap="md">
          {eventTypes.length === 0 && !loading ? (
            <Text c="dimmed" ta="center" py="xl">
              Типы встреч еще не созданы. Нажмите «Добавить тип встречи» для создания.
            </Text>
          ) : (
            <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="md">
              {eventTypes.map((eventType) => (
                <Card key={eventType.id} withBorder padding="md" radius="md">
                  <Stack gap="xs">
                    <Group justify="space-between" wrap="nowrap">
                      <Text fw={600} size="lg" lineClamp={1}>
                        {eventType.name}
                      </Text>
                      <Badge
                        variant="light"
                        leftSection={<IconClock size={12} />}
                      >
                        {eventType.durationMinutes} мин
                      </Badge>
                    </Group>
                    <Text size="sm" c="dimmed" lineClamp={2} style={{ minHeight: 40 }}>
                      {eventType.description}
                    </Text>
                    <Group justify="flex-end" mt="xs">
                      <ActionIcon
                        aria-label={`Редактировать тип встречи ${eventType.name}`}
                        color="blue"
                        variant="subtle"
                        onClick={() => openEditModal(eventType)}
                      >
                        <IconEdit size={18} />
                      </ActionIcon>
                      <ActionIcon
                        aria-label={`Удалить тип встречи ${eventType.name}`}
                        color="red"
                        variant="subtle"
                        onClick={() => handleDelete(eventType.id)}
                      >
                        <IconTrash size={18} />
                      </ActionIcon>
                    </Group>
                  </Stack>
                </Card>
              ))}
            </SimpleGrid>
          )}
        </Stack>
      </Card>

      {totalPages > 1 && (
        <Group justify="center" mt="md">
          <Pagination value={page} onChange={setPage} total={totalPages} />
        </Group>
      )}

      <Modal
        opened={modalOpened}
        onClose={() => setModalOpened(false)}
        title={editingEvent ? 'Редактировать тип встречи' : 'Создать тип встречи'}
        size="md"
      >
        <Stack>
          <TextInput
            label="Название"
            placeholder="Например: Консультация, Мастер-класс"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
          <Textarea
            label="Описание"
            placeholder="Опишите этот тип встречи"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            required
            minRows={3}
          />
          <NumberInput
            label="Длительность (минуты)"
            description="Минимум 5 минут, максимум 1440 минут (24 часа)"
            value={durationMinutes}
            onChange={(value) => setDurationMinutes(value)}
            min={5}
            max={1440}
            required
          />
          <Group justify="flex-end" mt="md">
            <Button variant="default" onClick={() => setModalOpened(false)}>
              Отмена
            </Button>
            <Button onClick={handleSubmit}>
              {editingEvent ? 'Обновить' : 'Создать'}
            </Button>
          </Group>
        </Stack>
      </Modal>
    </div>
  );
}
