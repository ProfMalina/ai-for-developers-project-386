import { Page, expect } from '@playwright/test';
import { BasePage } from './BasePage';

/**
 * Page Object for Booking Page (/book/:eventTypeId)
 */
export class BookingPage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  /**
   * Navigate to booking page for specific event type
   */
  async goto(eventTypeId: string): Promise<void> {
    await super.goto(`/book/${eventTypeId}`);
  }

  /**
   * Verify booking page is loaded
   */
  async expectLoaded(eventTypeName: string): Promise<void> {
    await this.expectHeading(`Бронирование: ${eventTypeName}`);
    // Verify stepper is present (Mantine Stepper renders steps with specific text)
    await expect(this.page.getByText('Шаг 1: Выберите дату').first()).toBeVisible();
    await expect(this.page.getByText('Шаг 2: Выберите время').first()).toBeVisible();
    await expect(this.page.getByText('Шаг 3: Ваши данные').first()).toBeVisible();
  }

  /**
   * Step 1: Select a date
   */
  async selectDate(date: Date): Promise<void> {
    // Click date input
    await this.page.getByRole('textbox', { name: /дата встречи/i }).click();
    await this.page.waitForTimeout(500);

    // Click the day in calendar
    const dayButton = this.page.getByRole('button', { name: String(date.getDate()) }).first();
    await dayButton.click();
    await this.page.waitForTimeout(500);
  }

  /**
   * Step 2: Select a time slot
   */
  async selectTimeSlot(time: string): Promise<void> {
    await this.page.getByRole('button', { name: time }).click();
    await this.page.waitForTimeout(500);
  }

  /**
   * Verify time slots are displayed
   */
  async expectTimeSlotsVisible(): Promise<void> {
    const slotButtons = this.page.getByRole('button', { name: /\d{2}:\d{2}/ });
    await expect(slotButtons.first()).toBeVisible({ timeout: 5000 });
  }

  /**
   * Verify time slot is disabled (booked)
   */
  async expectTimeSlotDisabled(time: string): Promise<void> {
    const slotButton = this.page.getByRole('button', { name: time });
    await expect(slotButton).toBeDisabled();
  }

  /**
   * Step 3: Fill guest details
   */
  async fillGuestDetails(name: string, email: string): Promise<void> {
    await this.page.getByRole('textbox', { name: /ваше имя/i }).fill(name);
    await this.page.getByRole('textbox', { name: /email/i }).fill(email);
  }

  /**
   * Verify summary card shows correct info
   */
  async expectSummaryInfo(info: {
    eventType?: string;
    date?: string;
    time?: string;
    duration?: string;
  }): Promise<void> {
    const summary = this.page.locator('article').last();

    if (info.eventType) {
      await expect(summary).toContainText(info.eventType);
    }
    if (info.date) {
      await expect(summary).toContainText(info.date);
    }
    if (info.time) {
      await expect(summary).toContainText(info.time);
    }
    if (info.duration) {
      await expect(summary).toContainText(info.duration);
    }
  }

  /**
   * Submit booking
   */
  async submitBooking(): Promise<void> {
    await this.page.getByRole('button', { name: 'Подтвердить бронирование' }).click();
    await this.page.waitForLoadState('domcontentloaded');
    await this.page.waitForTimeout(1000);
  }

  /**
   * Verify booking success
   */
  async expectBookingSuccess(): Promise<void> {
    // Should redirect to home page or show success message
    await this.expectNotification('Бронирование успешно создано')
      .catch(async () => {
        // Alternative: check if redirected to home
        await expect(this.page).toHaveURL(/.*\/$/);
      });
  }

  /**
   * Verify booking error (conflict)
   */
  async expectBookingError(): Promise<void> {
    await this.expectError('Ошибка')
      .catch(() => {
        // Alternative error message
        this.expectError('уже забронирован');
      });
  }

  /**
   * Verify form validation error
   */
  async expectValidationError(field: 'name' | 'email', message: string): Promise<void> {
    const errorElement = field === 'name'
      ? this.page.getByText(/имя обязательно/i)
      : this.page.getByText(/некорректный email/i);
    await expect(errorElement).toBeVisible();
  }

  /**
   * Try to submit with invalid data (should fail validation)
   */
  async submitInvalidBooking(): Promise<void> {
    // Try to click next/submit when form is invalid
    const nextButton = this.page.getByRole('button', { name: /далее|подтвердить/i });
    if (await nextButton.isEnabled()) {
      await nextButton.click();
    }
  }

  /**
   * Go back to step 1
   */
  async goToStep(step: number): Promise<void> {
    await this.page.getByText(`Шаг ${step}`).click();
    await this.page.waitForTimeout(300);
  }
}
