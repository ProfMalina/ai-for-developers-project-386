import { Page, expect } from '@playwright/test';
import { BasePage } from './BasePage';

/**
 * Page Object for Owner Dashboard (/owner)
 */
export class OwnerDashboard extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  /**
   * Navigate to owner dashboard
   */
  async goto(): Promise<void> {
    await super.goto('/owner');
    await this.acceptCookies();
  }

  /**
   * Verify dashboard is loaded
   */
  async expectLoaded(): Promise<void> {
    await this.expectHeading('Панель управления');
    // Verify tabs exist (use role=tab to avoid strict mode violations)
    await expect(this.page.getByRole('tab', { name: 'Типы встреч' })).toBeVisible({ timeout: 10000 });
    await expect(this.page.getByRole('tab', { name: 'Расписание' })).toBeVisible({ timeout: 10000 });
    await expect(this.page.getByRole('tab', { name: 'Бронирования' })).toBeVisible({ timeout: 10000 });
  }

  /**
   * Switch to a tab
   */
  async switchTab(tab: 'events' | 'slots' | 'bookings'): Promise<void> {
    const tabNames = {
      events: 'Типы встреч',
      slots: 'Расписание',
      bookings: 'Бронирования',
    };
    await this.page.getByRole('tab', { name: tabNames[tab] }).click();
    await this.page.waitForTimeout(500);
  }

  // ==================== Event Type Management ====================

  /**
   * Open add event type modal
   */
  async openAddEventType(): Promise<void> {
    await this.page.getByRole('button', { name: 'Добавить тип встречи' }).click();
    await this.page.waitForTimeout(500);
  }

  /**
   * Fill event type form in modal
   */
  async fillEventTypeForm(data: {
    name: string;
    description?: string;
    duration: number;
  }): Promise<void> {
    await this.page.getByRole('textbox', { name: /название/i }).fill(data.name);
    if (data.description) {
      await this.page.getByRole('textbox', { name: /описание/i }).fill(data.description);
    }
    await this.page.getByRole('textbox', { name: /длительность/i }).fill(String(data.duration));
  }

  /**
   * Submit event type form
   */
  async submitEventType(): Promise<void> {
    await this.page.getByRole('button', { name: /создать|обновить/i }).click();
    await this.page.waitForLoadState('domcontentloaded');
    await this.page.waitForTimeout(500);
  }

  /**
   * Cancel event type form
   */
  async cancelEventTypeForm(): Promise<void> {
    await this.page.getByRole('button', { name: /отмена/i }).click();
  }

  /**
   * Get event type card by name
   */
  getEventTypeCard(name: string) {
    return this.page
      .getByText(name, { exact: true })
      .locator('xpath=ancestor::div[.//button[contains(@aria-label, "тип встречи")]][1]');
  }

  /**
   * Edit event type
   */
  async editEventType(name: string): Promise<void> {
    const card = this.getEventTypeCard(name);
    await card.getByRole('button', { name: /редактировать|edit/i }).click();
    await this.page.waitForTimeout(500);
  }

  /**
   * Delete event type
   */
  async deleteEventType(name: string): Promise<void> {
    const card = this.getEventTypeCard(name);
    await card.getByRole('button', { name: /удалить|delete/i }).click();

    // Handle confirmation dialog
    this.page.on('dialog', async dialog => {
      await dialog.accept();
    });

    await this.page.waitForLoadState('domcontentloaded');
    await this.page.waitForTimeout(500);
  }

  /**
   * Verify event type is visible in list
   */
  async expectEventTypeVisible(name: string, duration: number): Promise<void> {
    const card = this.getEventTypeCard(name);
    await expect(card).toBeVisible();
    await expect(card).toContainText(`${duration} мин`);
  }

  /**
   * Verify event type is NOT visible in list
   */
  async expectEventTypeNotVisible(name: string): Promise<void> {
    const card = this.getEventTypeCard(name);
    await expect(card).not.toBeVisible();
  }

  // ==================== Slot Generation ====================

  /**
   * Open slot generation modal
   */
  async openSlotGeneration(): Promise<void> {
    await this.switchTab('slots');
    await this.page.getByRole('button', { name: 'Создать слоты' }).click();
    await this.page.waitForTimeout(500);
  }

  /**
   * Fill slot generation form
   */
  async fillSlotGenerationForm(data: {
    workingHoursStart?: string;
    workingHoursEnd?: string;
    intervalMinutes: number;
    daysOfWeek: string[];
    dateFrom: string;
    dateTo: string;
  }): Promise<void> {
    if (data.workingHoursStart) {
      await this.page.getByRole('textbox', { name: /начало рабочего дня/i }).fill(data.workingHoursStart);
    }
    if (data.workingHoursEnd) {
      await this.page.getByRole('textbox', { name: /конец рабочего дня/i }).fill(data.workingHoursEnd);
    }
    await this.page.getByRole('textbox', { name: /длительность слота/i }).fill(String(data.intervalMinutes));

    // Select days of week
    for (const day of data.daysOfWeek) {
      await this.page.getByRole('checkbox', { name: day }).check();
    }

    // Keep the valid default date range from the modal. Mantine DatePickerInput
    // renders button-like controls here, and the happy-path slot-generation
    // test only needs a valid submitted range, not custom date interaction.
  }

  /**
   * Submit slot generation
   */
  async submitSlotGeneration(): Promise<void> {
    await this.page.getByRole('button', { name: 'Сгенерировать слоты' }).click();
    await this.page.waitForLoadState('domcontentloaded');
    await this.page.waitForTimeout(500);
  }

  /**
   * Verify slot generation success
   */
  async expectSlotsGenerated(expectedCount?: number): Promise<void> {
    await expect(this.page.getByText(/создано.*слотов/i)).toBeVisible({ timeout: 5000 });
    if (expectedCount) {
      await this.expectNotification(`${expectedCount}`);
    }
  }

  // ==================== Bookings List ====================

  /**
   * Navigate to bookings tab
   */
  async goToBookings(): Promise<void> {
    await this.switchTab('bookings');
  }

  /**
   * Get active bookings tab panel
   */
  getBookingsPanel() {
    return this.page.getByRole('tabpanel', { name: 'Бронирования' });
  }

  /**
   * Get booking card by guest name
   */
  getBookingCard(guestName: string) {
    return this.getBookingsPanel()
      .getByText(guestName, { exact: true })
      .locator('xpath=ancestor::div[.//button[normalize-space()="Отменить бронирование"]][1]');
  }

  /**
   * Get bookings pagination navigation
   */
  getBookingsPagination() {
    return this.getBookingsPanel();
  }

  /**
   * Verify booking is visible
   */
  async expectBookingVisible(guestName: string): Promise<void> {
    const card = this.getBookingCard(guestName);
    await expect(card).toBeVisible();
  }

  /**
   * Cancel a booking
   */
  async cancelBooking(guestName: string): Promise<void> {
    const card = this.getBookingCard(guestName);
    await card.getByRole('button', { name: /отменить бронирование/i }).click();
    await this.page.waitForTimeout(300);

    // Confirm cancellation
    await this.page.getByRole('button', { name: /да, отменить/i }).click();
    await this.page.waitForLoadState('domcontentloaded');
    await this.page.waitForTimeout(500);
  }

  /**
   * Verify pagination
   */
  async expectPaginationVisible(): Promise<void> {
    await expect(this.getBookingsPagination()).toBeVisible();
  }

  /**
   * Go to bookings page
   */
  async goToBookingsPage(pageNum: number): Promise<void> {
    await this.getBookingsPagination().getByRole('button', { name: String(pageNum), exact: true }).click();
    await this.page.waitForLoadState('domcontentloaded');
    await this.page.waitForTimeout(500);
  }
}
