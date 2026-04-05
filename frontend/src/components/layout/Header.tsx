import { Group, Button, Container, Title } from '@mantine/core';
import { Link, useLocation } from 'react-router-dom';

export function Header() {
  const location = useLocation();
  const isOwner = location.pathname.startsWith('/owner');

  return (
    <Container fluid>
      <Group justify="space-between" h="100%">
        <Group>
          <Title order={3} c="blue">
            <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
              📅 Calendar Booking
            </Link>
          </Title>
        </Group>
        <Group>
          <Button
            variant={isOwner ? 'light' : 'outline'}
            component={Link}
            to="/"
          >
            Guest View
          </Button>
          <Button
            variant={isOwner ? 'filled' : 'outline'}
            component={Link}
            to="/owner"
          >
            Owner View
          </Button>
        </Group>
      </Group>
    </Container>
  );
}
