import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  isSlotAvailable,
  doSlotsOverlap,
  formatSlotTime,
  formatDate,
  calculateSlotEndTime,
} from './slots';

describe('Slot Utilities', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  describe('isSlotAvailable', () => {
    it('should return true for future slots', () => {
      const futureTime = '2026-04-08T10:00:00Z';
      vi.setSystemTime(new Date('2026-04-07T12:00:00Z'));

      expect(isSlotAvailable(futureTime)).toBe(true);
    });

    it('should return false for past slots', () => {
      const pastTime = '2026-04-06T10:00:00Z';
      vi.setSystemTime(new Date('2026-04-07T12:00:00Z'));

      expect(isSlotAvailable(pastTime)).toBe(false);
    });

    it('should return false for current time slot', () => {
      const currentTime = '2026-04-07T12:00:00Z';
      vi.setSystemTime(new Date('2026-04-07T12:00:00Z'));

      expect(isSlotAvailable(currentTime)).toBe(false);
    });
  });

  describe('doSlotsOverlap', () => {
    it('should return true for overlapping slots', () => {
      const start1 = '2026-04-08T10:00:00Z';
      const end1 = '2026-04-08T11:00:00Z';
      const start2 = '2026-04-08T10:30:00Z';
      const end2 = '2026-04-08T11:30:00Z';

      expect(doSlotsOverlap(start1, end1, start2, end2)).toBe(true);
    });

    it('should return false for non-overlapping slots', () => {
      const start1 = '2026-04-08T10:00:00Z';
      const end1 = '2026-04-08T11:00:00Z';
      const start2 = '2026-04-08T11:00:00Z';
      const end2 = '2026-04-08T12:00:00Z';

      expect(doSlotsOverlap(start1, end1, start2, end2)).toBe(false);
    });

    it('should return false when slot2 is before slot1', () => {
      const start1 = '2026-04-08T11:00:00Z';
      const end1 = '2026-04-08T12:00:00Z';
      const start2 = '2026-04-08T09:00:00Z';
      const end2 = '2026-04-08T10:00:00Z';

      expect(doSlotsOverlap(start1, end1, start2, end2)).toBe(false);
    });

    it('should return true when one slot contains another', () => {
      const start1 = '2026-04-08T10:00:00Z';
      const end1 = '2026-04-08T12:00:00Z';
      const start2 = '2026-04-08T10:30:00Z';
      const end2 = '2026-04-08T11:30:00Z';

      expect(doSlotsOverlap(start1, end1, start2, end2)).toBe(true);
    });
  });

  describe('formatSlotTime', () => {
    it('should format time correctly', () => {
      const time = '2026-04-08T10:00:00Z';
      vi.setSystemTime(new Date('2026-04-08T12:00:00Z')); // UTC+2 timezone simulation

      const formatted = formatSlotTime(time);
      expect(typeof formatted).toBe('string');
      expect(formatted).toMatch(/\d{2}:\d{2}/);
    });
  });

  describe('formatDate', () => {
    it('should format date in Russian locale', () => {
      const date = '2026-04-08T10:00:00Z';

      const formatted = formatDate(date, 'ru');
      expect(typeof formatted).toBe('string');
      expect(formatted).toContain('2026');
    });
  });

  describe('calculateSlotEndTime', () => {
    it('should calculate end time correctly for 30 minutes', () => {
      const startTime = '2026-04-08T10:00:00Z';
      const durationMinutes = 30;

      const endTime = calculateSlotEndTime(startTime, durationMinutes);
      expect(endTime).toBe('2026-04-08T10:30:00.000Z');
    });

    it('should calculate end time correctly for 60 minutes', () => {
      const startTime = '2026-04-08T10:00:00Z';
      const durationMinutes = 60;

      const endTime = calculateSlotEndTime(startTime, durationMinutes);
      expect(endTime).toBe('2026-04-08T11:00:00.000Z');
    });

    it('should handle day boundary crossing', () => {
      const startTime = '2026-04-08T23:30:00Z';
      const durationMinutes = 60;

      const endTime = calculateSlotEndTime(startTime, durationMinutes);
      expect(endTime).toBe('2026-04-09T00:30:00.000Z');
    });
  });
});
