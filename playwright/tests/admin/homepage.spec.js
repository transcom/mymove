/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../utils/adminTest';

test('admin home page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  // redirects to office users page after login
  await expect(page.getByRole('heading', { name: 'Office Users' })).toBeVisible();
});
