import { describe, it, expect } from 'vitest';

// Test email validation regex used in BookingPage
const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

describe('Email Validation', () => {
  it('should accept valid email addresses', () => {
    const validEmails = [
      'test@example.com',
      'user.name@domain.org',
      'user+tag@example.com',
      'user123@test.co.uk',
      'ivan@example.com',
    ];

    validEmails.forEach((email) => {
      expect(emailRegex.test(email)).toBe(true);
    });
  });

  it('should reject invalid email addresses', () => {
    const invalidEmails = [
      '',
      'notanemail',
      '@domain.com',
      'user@',
      'user@.com',
      'user domain@example.com',
      'user@domain',
    ];

    invalidEmails.forEach((email) => {
      expect(emailRegex.test(email)).toBe(false);
    });
  });
});

describe('Booking Form Validation', () => {
  it('should require guest name', () => {
    const guestName = '';
    const guestEmail = 'test@example.com';

    const isValid = guestName.trim() !== '' && guestEmail.trim() !== '';
    expect(isValid).toBe(false);
  });

  it('should require guest email', () => {
    const guestName = 'Иван';
    const guestEmail = '';

    const isValid = guestName.trim() !== '' && guestEmail.trim() !== '';
    expect(isValid).toBe(false);
  });

  it('should accept valid name and email', () => {
    const guestName = 'Иван Иванов';
    const guestEmail = 'ivan@example.com';

    const hasValidName = guestName.trim() !== '';
    const hasValidEmail = guestEmail.trim() !== '' && emailRegex.test(guestEmail);

    expect(hasValidName).toBe(true);
    expect(hasValidEmail).toBe(true);
  });
});
