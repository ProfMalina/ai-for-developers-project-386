import { Page, expect } from '@playwright/test';

/**
 * Base Page Object for all pages
 * Provides common functionality
 */
export class BasePage {
  readonly page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  /**
   * Navigate to a URL
   */
  async goto(url: string): Promise<void> {
    await this.page.goto(url);
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Get current URL
   */
  getUrl(): string {
    return this.page.url();
  }

  /**
   * Wait for page to be ready
   */
  async waitForReady(): Promise<void> {
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Accept cookie banner if present
   */
  async acceptCookies(): Promise<void> {
    const cookieBanner = this.page.getByRole('button', { name: /принять|accept|согласен|ok/i }).first();
    if (await cookieBanner.isVisible({ timeout: 2000 }).catch(() => false)) {
      await cookieBanner.click();
    }
  }

  /**
   * Take screenshot (for debugging)
   */
  async screenshot(name: string): Promise<void> {
    await this.page.screenshot({ path: `e2e/screens/${name}.png`, fullPage: true });
  }

  /**
   * Verify page title or heading
   */
  async expectHeading(text: string): Promise<void> {
    await expect(this.page.getByRole('heading', { name: text })).toBeVisible();
  }

  /**
   * Verify element is visible
   */
  async expectVisible(selector: string): Promise<void> {
    await expect(this.page.locator(selector)).toBeVisible();
  }

  /**
   * Verify element contains text
   */
  async expectContainsText(selector: string, text: string): Promise<void> {
    await expect(this.page.locator(selector)).toContainText(text);
  }

  /**
   * Click element by role
   */
  async clickByRole(role: string, name: string | RegExp): Promise<void> {
    await this.page.getByRole(role as any, { name }).click();
  }

  /**
   * Wait for notification (Mantine)
   */
  async expectNotification(text: string): Promise<void> {
    await expect(this.page.getByText(text)).toBeVisible({ timeout: 5000 });
  }

  /**
   * Wait for error message
   */
  async expectError(text: string): Promise<void> {
    await expect(this.page.getByText(text)).toBeVisible({ timeout: 5000 });
  }
}
