/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../utils/admin/adminTest';

test.describe('Users Page', () => {
  test('successfully navigates to users page', async ({ page, adminPage }) => {
    await adminPage.signInAsNewAdminUser();

    await page.getByRole('menuitem', { name: 'Users', exact: true }).click();
    await adminPage.waitForPage.adminPage();
    expect(page.url()).toContain('/system/users');
    await expect(page.locator('header')).toContainText('Users');

    const columnLabels = ['Id', 'Email', 'Active', 'Created at'];
    for (const label of columnLabels) {
      await expect(page.locator('table')).toContainText(label);
    }
  });
});

test.describe('Users Details Show Page', () => {
  test('pulls up details page for a user', async ({ page, adminPage }) => {
    const officeUser = await adminPage.testHarness.buildOfficeUserWithTOOAndTIO();
    const email = officeUser.login_gov_email;
    await adminPage.signInAsNewAdminUser();

    await page.getByRole('menuitem', { name: 'Users', exact: true }).click();
    await adminPage.waitForPage.adminPage();
    expect(page.url()).toContain('/system/users');

    await page.getByLabel('Search by User Id or Email').click();
    await page.getByLabel('Search by User Id or Email').fill(email);
    await page.getByLabel('Search by User Id or Email').press('Enter');
    await expect(page.getByRole('cell', { name: email })).toHaveCount(1);
    await page.getByText(email).click();
    await adminPage.waitForPage.adminPage();

    // check that the user's login.gov email is shown in the page title
    await expect(page.getByRole('heading', { name: email })).toBeVisible();

    const labels = ['User ID', 'User email', 'Active', 'Created at', 'Updated at'];
    for (const label of labels) {
      await expect(page.getByRole('paragraph').filter({ hasText: label })).toBeVisible();
    }
  });
});

test.describe('Users Details Edit Page', () => {
  test('pulls up edit page for a user', async ({ page, adminPage }) => {
    const officeUser = await adminPage.testHarness.buildOfficeUserWithTOOAndTIO();
    const email = officeUser.login_gov_email;
    await adminPage.signInAsNewAdminUser();

    await page.getByRole('menuitem', { name: 'Users', exact: true }).click();
    await adminPage.waitForPage.adminPage();
    expect(page.url()).toContain('/system/users');

    await page.getByLabel('Search by User Id or Email').click();
    await page.getByLabel('Search by User Id or Email').fill(email);
    await page.getByLabel('Search by User Id or Email').press('Enter');
    await expect(page.getByRole('cell', { name: email })).toHaveCount(1);
    await page.getByText(email).click();
    await adminPage.waitForPage.adminPage();

    await page.getByRole('link', { name: 'Edit' }).click();
    await adminPage.waitForPage.adminPage();

    // check page content
    const pageContent = [
      'Id',
      'Login gov email',
      'Active',
      'Revoke admin session',
      'Revoke office session',
      'Revoke mil session',
    ];
    for (const label of pageContent) {
      await expect(page.locator('label').getByText(label, { exact: true })).toBeVisible();
    }

    // set the user to the active status they did NOT have before
    const activeStatus = await page.locator('div:has(label :text-is("Active")) >> input[name="active"]').inputValue();

    const newStatus = (activeStatus !== 'true').toString();
    await page.locator('div:has(label :text-is("Active")) >> #active').click();
    await page.locator(`ul[aria-labelledby="active-label"] >> li[data-value="${newStatus}"]`).click();

    await page.getByRole('button', { name: 'Save' }).click();
    await adminPage.waitForPage.adminPage();

    // check that user was deactivated
    await expect(page.locator(`tr:has(:text("${email}")) >> td.column-active >> svg`)).toHaveAttribute(
      'data-testid',
      newStatus,
    );
  });
});
