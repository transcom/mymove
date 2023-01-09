// @ts-check
const { test, expect } = require('../utils/officeTest');

test.describe('The document viewer', () => {
  test.describe('When not logged in', () => {
    test('shows page not found', async ({ page }) => {
      await page.goto('/moves/foo/documents');
      await expect(page.getByText('Welcome')).toBeVisible();
      // sign in button not in header
      await expect(page.locator('#main').getByRole('button', { name: 'Sign in' })).toBeVisible();
    });
  });
});
