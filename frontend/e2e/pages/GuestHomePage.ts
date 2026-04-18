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
   * Note: Mantine Card renders as a generic container, not a semantic article
   */
  getEventTypeCard(name: string) {
    return this.page
      .getByRole('heading', { name, exact: true })
      .locator('xpath=ancestor::div[.//a[normalize-space()="Забронировать"]][1]');
  }

  /**
   * Verify event type is displayed
   */
  async expectEventTypeVisible(name: string, duration: number): Promise<void> {
    // First verify the text is visible
    await expect(this.page.getByText(name).first()).toBeVisible();
    // Then verify the duration badge
    await expect(this.page.getByText(`${duration} мин`).first()).toBeVisible();
  }

  /**
   * Click "Забронировать" button on event type card
   * Note: Button uses component={Link} so it renders as <a>, not <button>
   */
  async bookEventType(name: string): Promise<void> {
    const card = this.getEventTypeCard(name);
    await card.getByRole('link', { name: 'Забронировать' }).click();
    await this.page.waitForLoadState('domcontentloaded');
    await this.page.waitForTimeout(500);
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
    await this.page.waitForLoadState('domcontentloaded');
    await this.page.waitForTimeout(500);
  }
}
