/**
 * Shared test utilities and helpers
 */

import { Page, expect } from '@playwright/test';

/**
 * Wait for page to be fully loaded (network idle)
 */
export async function waitForPageLoad(page: Page): Promise<void> {
  await page.waitForLoadState('networkidle');
}

/**
 * Select a date in Mantine DatePickerInput
 * Uses Mantine's date picker structure
 */
export async function selectDate(page: Page, date: Date): Promise<void> {
  const day = date.getDate();
  const month = date.getMonth();
  const year = date.getFullYear();

  // Click on the date input to open calendar
  await page.getByRole('textbox', { name: /дата/i }).click();

  // Navigate to correct month/year if needed
  // Click on the day
  await page.getByRole('button', { name: String(day) }).click();
}

/**
 * Select a time slot button
 */
export async function selectTimeSlot(page: Page, time: string): Promise<void> {
  await page.getByRole('button', { name: time }).click();
}

/**
 * Fill booking form fields
 */
export async function fillBookingForm(page: Page, data: {
  name: string;
  email: string;
}): Promise<void> {
  await page.getByRole('textbox', { name: /имя/i }).fill(data.name);
  await page.getByRole('textbox', { name: /email/i }).fill(data.email);
}

/**
 * Fill event type form (for owner)
 */
export async function fillEventTypeForm(page: Page, data: {
  name: string;
  description?: string;
  duration: number;
}): Promise<void> {
  await page.getByRole('textbox', { name: /название/i }).fill(data.name);
  if (data.description) {
    await page.getByRole('textbox', { name: /описание/i }).fill(data.description);
  }
  await page.getByRole('spinbutton', { name: /длительность/i }).fill(String(data.duration));
}

/**
 * Accept cookie consent banner if it appears
 */
export async function acceptCookieBanner(page: Page): Promise<void> {
  const cookieBanner = page.getByRole('button', { name: /принять|accept|ok/i }).first();
  if (await cookieBanner.isVisible()) {
    await cookieBanner.click();
  }
}

/**
 * Switch language (if i18n is implemented)
 */
export async function switchLanguage(page: Page, language: 'ru' | 'en'): Promise<void> {
  const langButton = page.getByRole('button', { name: new RegExp(language === 'ru' ? 'english|en' : 'русский|ru', 'i') });
  if (await langButton.isVisible()) {
    await langButton.click();
  }
}

/**
 * Switch theme (if theme switching is implemented)
 */
export async function switchTheme(page: Page, theme: 'light' | 'dark'): Promise<void> {
  const themeButton = page.getByRole('button', { name: /тема|theme/i });
  if (await themeButton.isVisible()) {
    await themeButton.click();
  }
}

/**
 * Verify API error response handling
 */
export function expectApiError(page: Page, statusCode: number, message?: string): void {
  page.on('response', async (response) => {
    if (response.status() === statusCode) {
      if (message) {
        await expect(page.getByText(message)).toBeVisible();
      }
    }
  });
}

/**
 * Wait for notification message (Mantine notifications)
 */
export async function waitForNotification(page: Page, text: string): Promise<void> {
  await expect(page.getByText(text)).toBeVisible({ timeout: 5000 });
}

/**
 * Check responsive layout at specific viewport
 */
export async function setViewport(page: Page, width: number, height: number = 800): Promise<void> {
  await page.setViewportSize({ width, height });
  await waitForPageLoad(page);
}
