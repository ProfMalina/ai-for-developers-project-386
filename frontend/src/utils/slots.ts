import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

dayjs.extend(utc);
dayjs.extend(timezone);

/**
 * Check if a time slot is available (hasn't started yet)
 */
export const isSlotAvailable = (slotStartTime: string): boolean => {
  return dayjs.utc(slotStartTime).local().isAfter(dayjs());
};

/**
 * Check if two time slots overlap
 */
export const doSlotsOverlap = (
  start1: string,
  end1: string,
  start2: string,
  end2: string
): boolean => {
  const s1 = dayjs(start1);
  const e1 = dayjs(end1);
  const s2 = dayjs(start2);
  const e2 = dayjs(end2);

  return s1 < e2 && s2 < e1;
};

/**
 * Format time slot for display
 */
export const formatSlotTime = (startTime: string, format = 'HH:mm'): string => {
  return dayjs.utc(startTime).local().format(format);
};

/**
 * Format date for display
 */
export const formatDate = (date: string, locale = 'ru'): string => {
  return dayjs(date).locale(locale).format('D MMMM YYYY');
};

/**
 * Calculate slot end time based on duration
 */
export const calculateSlotEndTime = (startTime: string, durationMinutes: number): string => {
  return dayjs(startTime).add(durationMinutes, 'minute').toISOString();
};
