/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../utils/adminTest';

test.describe('Offices Page', () => {
  test('successfully navigates to offices page', async ({ page, adminPage }) => {
    await adminPage.signInAsNewAdminUser();
    await page.getByRole('menuitem', { name: 'Offices' }).click();

    expect(page.url()).toContain('/system/offices');
    await expect(page.locator('header')).toContainText('Offices');
    await expect(page.getByLabel('Search by Office Name')).toBeEditable();

    const columnLabels = ['Id', 'Name', 'Latitude', 'Longitude', 'Gbloc'];
    await adminPage.expectRoleLabelsByText('columnheader', columnLabels);
  });
});
