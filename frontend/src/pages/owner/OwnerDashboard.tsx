import { useState } from 'react';
import { Container, Tabs, Title } from '@mantine/core';
import { EventTypeManagement } from '../../components/owner/EventTypeManagement';
import { BookingsList } from '../../components/owner/BookingsList';

export function OwnerDashboard() {
  const [activeTab, setActiveTab] = useState<string | null>('events');

  return (
    <Container size="xl" py="xl">
      <Title order={1} size={36} mb="xl">
        Панель управления
      </Title>

      <Tabs value={activeTab} onChange={setActiveTab}>
        <Tabs.List>
          <Tabs.Tab value="events">Типы встреч</Tabs.Tab>
          <Tabs.Tab value="bookings">Бронирования</Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="events" pt="xl">
          <EventTypeManagement />
        </Tabs.Panel>

        <Tabs.Panel value="bookings" pt="xl">
          <BookingsList />
        </Tabs.Panel>
      </Tabs>
    </Container>
  );
}
