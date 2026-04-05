import { Container, Title, Text, Button, Group } from '@mantine/core';
import { Link } from 'react-router-dom';

export function NotFound() {
  return (
    <Container size="sm" py="xl" ta="center">
      <Title order={1} size={96}>404</Title>
      <Title order={2} mb="md">Страница не найдена</Title>
      <Text c="dimmed" mb="xl">
        Страница, которую вы ищете, не существует.
      </Text>
      <Group justify="center">
        <Button component={Link} to="/" size="lg">
          На главную
        </Button>
      </Group>
    </Container>
  );
}
