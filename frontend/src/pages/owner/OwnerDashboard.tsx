import { useState, useEffect } from 'react';
import { Container, Tabs, Title } from '@mantine/core';
import { EventTypeManagement } from '../../components/owner/EventTypeManagement';
import { BookingsList } from '../../components/owner/BookingsList';

export function OwnerDashboard() {
  const [activeTab, setActiveTab] = useState<string | null>('events');

  useEffect(() => {
    document.title = 'Owner Dashboard - Calendar Booking';
  }, []);

  return (
    <Container size="xl" py="xl">
      <Title order={1} mb="xl">Owner Dashboard</Title>

      <Tabs value={activeTab} onChange={setActiveTab}>
        <Tabs.List>
          <Tabs.Tab value="events">Event Types</Tabs.Tab>
          <Tabs.Tab value="bookings">Upcoming Bookings</Tabs.Tab>
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
