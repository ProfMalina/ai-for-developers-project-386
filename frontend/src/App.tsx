import { Routes, Route } from 'react-router-dom';
import { AppShell } from '@mantine/core';
import { Header } from './components/layout/Header';
import { OwnerDashboard } from './pages/owner/OwnerDashboard';
import { GuestHome } from './pages/guest/GuestHome';
import { BookingPage } from './pages/guest/BookingPage';
import { NotFound } from './pages/NotFound';

export default function App() {
  return (
    <AppShell
      header={{ height: 60 }}
      padding="md"
    >
      <AppShell.Header>
        <Header />
      </AppShell.Header>
      <AppShell.Main>
        <Routes>
          {/* Owner routes */}
          <Route path="/owner" element={<OwnerDashboard />} />

          {/* Guest/Public routes */}
          <Route path="/" element={<GuestHome />} />
          <Route path="/book/:eventTypeId" element={<BookingPage />} />

          {/* 404 */}
          <Route path="*" element={<NotFound />} />
        </Routes>
      </AppShell.Main>
    </AppShell>
  );
}
