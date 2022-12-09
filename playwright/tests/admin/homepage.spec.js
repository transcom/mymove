// @ts-check
const { test, expect } = require('@playwright/test');

const { signInAsNewAdminUser } = require('../utils/signIn');

test('admin home page', async ({ page }) => {
  await page.goto('/');

  await signInAsNewAdminUser(page);

  // redirects to office users page after login
  await expect(page.getByRole('heading', { name: 'Office Users' })).toBeVisible();
});
