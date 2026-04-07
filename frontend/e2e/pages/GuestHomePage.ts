import { Page, expect } from '@playwright/test';
import { BasePage } from './BasePage';

/**
 * Page Object for Guest Home Page (/)
 */
export class GuestHomePage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  /**
   * Navigate to guest home page
   */
  async goto(): Promise<void> {
    await super.goto('/');
    await this.acceptCookies();
  }

  /**
   * Verify guest home page is loaded
   */
  async expectLoaded(): Promise<void> {
    await this.expectHeading('Забронировать встречу');
    await expect(this.page.getByText('Выберите тип встречи и удобное для вас время')).toBeVisible();
  }

  /**
   * Get event type card by name
   */
  getEventTypeCard(name: string) {
    return this.page.locator('article').filter({ hasText: name });
  }

  /**
   * Verify event type is displayed
   */
  async expectEventTypeVisible(name: string, duration: number): Promise<void> {
    const card = this.getEventTypeCard(name);
    await expect(card).toBeVisible();
    await expect(card).toContainText(`${duration} мин`);
  }

  /**
   * Click "Забронировать" button on event type card
   */
  async bookEventType(name: string): Promise<void> {
    const card = this.getEventTypeCard(name);
    await card.getByRole('button', { name: 'Забронировать' }).click();
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Verify pagination is present
   */
  async expectPaginationVisible(): Promise<void> {
    await expect(this.page.getByRole('navigation')).toBeVisible();
  }

  /**
   * Go to next page (pagination)
   */
  async goToPage(pageNum: number): Promise<void> {
    await this.page.getByRole('button', { name: String(pageNum) }).click();
    await this.page.waitForLoadState('networkidle');
  }
}
