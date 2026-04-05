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
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { IconPlus, IconEdit, IconTrash } from '@tabler/icons-react';
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
    } catch (error) {
      console.error('Failed to fetch event types:', error);
      notifications.show({
        title: 'Error',
        message: 'Failed to load event types',
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
        title: 'Validation Error',
        message: 'All fields are required',
        color: 'red',
      });
      return;
    }

    if (typeof durationMinutes === 'number' && durationMinutes < 5) {
      notifications.show({
        title: 'Validation Error',
        message: 'Duration must be at least 5 minutes',
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
          title: 'Success',
          message: 'Event type updated successfully',
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
          title: 'Success',
          message: 'Event type created successfully',
          color: 'green',
        });
      }
      setModalOpened(false);
      fetchEventTypes();
    } catch (error) {
      console.error('Failed to save event type:', error);
      notifications.show({
        title: 'Error',
        message: error instanceof Error ? error.message : 'Failed to save event type',
        color: 'red',
      });
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this event type?')) {
      return;
    }

    try {
      await ownerApi.deleteEventType(id);
      notifications.show({
        title: 'Success',
        message: 'Event type deleted successfully',
        color: 'green',
      });
      fetchEventTypes();
    } catch (error) {
      console.error('Failed to delete event type:', error);
      notifications.show({
        title: 'Error',
        message: 'Failed to delete event type',
        color: 'red',
      });
    }
  };

  return (
    <div>
      <Group justify="space-between" mb="md">
        <Text size="lg" fw={500}>Event Types</Text>
        <Button leftSection={<IconPlus size={16} />} onClick={openCreateModal}>
          Add Event Type
        </Button>
      </Group>

      <Card withBorder pos="relative">
        <LoadingOverlay visible={loading} />
        <Stack gap="md">
          {eventTypes.length === 0 && !loading ? (
            <Text c="dimmed" ta="center" py="xl">
              No event types yet. Click "Add Event Type" to create one.
            </Text>
          ) : (
            eventTypes.map((eventType) => (
              <Card key={eventType.id} withBorder padding="md">
                <Group justify="space-between" mb="xs">
                  <Text fw={500}>{eventType.name}</Text>
                  <Group>
                    <Badge color="blue">{eventType.durationMinutes} min</Badge>
                    <ActionIcon color="blue" variant="subtle" onClick={() => openEditModal(eventType)}>
                      <IconEdit size={16} />
                    </ActionIcon>
                    <ActionIcon color="red" variant="subtle" onClick={() => handleDelete(eventType.id)}>
                      <IconTrash size={16} />
                    </ActionIcon>
                  </Group>
                </Group>
                <Text size="sm" c="dimmed">{eventType.description}</Text>
              </Card>
            ))
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
        title={editingEvent ? 'Edit Event Type' : 'Create Event Type'}
        size="md"
      >
        <Stack>
          <TextInput
            label="Name"
            placeholder="e.g., Consultation, Workshop"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
          <Textarea
            label="Description"
            placeholder="Describe this event type"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            required
            minRows={3}
          />
          <NumberInput
            label="Duration (minutes)"
            description="Minimum 5 minutes, maximum 1440 minutes (24 hours)"
            value={durationMinutes}
            onChange={(value) => setDurationMinutes(value)}
            min={5}
            max={1440}
            required
          />
          <Group justify="flex-end" mt="md">
            <Button variant="default" onClick={() => setModalOpened(false)}>
              Cancel
            </Button>
            <Button onClick={handleSubmit}>
              {editingEvent ? 'Update' : 'Create'}
            </Button>
          </Group>
        </Stack>
      </Modal>
    </div>
  );
}
