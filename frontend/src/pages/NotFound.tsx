import { Container, Title, Text, Button, Group } from '@mantine/core';
import { Link } from 'react-router-dom';

export function NotFound() {
  return (
    <Container size="sm" py="xl" ta="center">
      <Title order={1} size={96}>404</Title>
      <Title order={2} mb="md">Page Not Found</Title>
      <Text c="dimmed" mb="xl">
        The page you're looking for doesn't exist.
      </Text>
      <Group justify="center">
        <Button component={Link} to="/" size="lg">
          Go Home
        </Button>
      </Group>
    </Container>
  );
}
