import { Group, Button, Container, Title, rem } from '@mantine/core';
import { Link, useLocation } from 'react-router-dom';

export function Header() {
  const location = useLocation();
  const isOwner = location.pathname.startsWith('/owner');

  return (
    <Container fluid px="md">
      <Group justify="space-between" h="100%">
        <Group gap="xs">
          <Title order={3} c="blue" style={{ fontSize: rem(22) }}>
            <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
              📅 Бронирование
            </Link>
          </Title>
        </Group>
        <Group gap="xs">
          <Button
            variant={isOwner ? 'outline' : 'filled'}
            component={Link}
            to="/"
            size="compact-sm"
          >
            Гость
          </Button>
          <Button
            variant={isOwner ? 'filled' : 'outline'}
            component={Link}
            to="/owner"
            size="compact-sm"
          >
            Владелец
          </Button>
        </Group>
      </Group>
    </Container>
  );
}
