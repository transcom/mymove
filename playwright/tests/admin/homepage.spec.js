// @ts-check
const { test, expect } = require('../utils/adminTest');

test('admin home page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  // redirects to office users page after login
  await expect(page.getByRole('heading', { name: 'Office Users' })).toBeVisible();
});
