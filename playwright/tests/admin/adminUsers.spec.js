/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../utils/admin/adminTest';

test('Admin Users List Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await adminPage.waitForPage.adminPage();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  const columnLabels = ['Id', 'Email', 'First name', 'Last name', 'User Id', 'Active'];

  await adminPage.expectRoleLabelsByText('columnheader', columnLabels);
});

test('Admin User Create Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await adminPage.waitForPage.adminPage();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  await page.getByRole('link', { name: 'Create' }).click();
  await adminPage.waitForPage.adminPage();

  await expect(page.getByRole('heading', { name: 'Create Admin Users' })).toBeVisible();

  // we need to add the date to the email so that it is unique every time (only one record per email allowed in db)
  const testEmail = `cy.admin_user.${Date.now()}@example.com`;

  // create an admin user
  await page.getByLabel('Email').fill(testEmail);

  const firstName = page.getByLabel('First name');
  await firstName.focus();
  await firstName.fill('Cypress');

  const lastName = page.getByLabel('Last name');
  await lastName.focus();
  await lastName.fill('Test');

  // The autocomplete form results in multiple matching elements, so
  // pick the input element
  await page.getByLabel('Organization').click();
  await page.getByRole('option').first().click();

  await page.getByRole('button').filter({ hasText: 'Save' }).click();
  await adminPage.waitForPage.adminPage();

  // redirected to edit details page
  const adminUserID = await page.locator('#id').inputValue();

  await expect(page.getByRole('heading', { name: `Admin Users #${adminUserID}` })).toBeVisible();

  await expect(page.locator('#email')).toHaveValue(testEmail);
  await expect(page.locator('#firstName')).toHaveValue('Cypress');
  await expect(page.locator('#lastName')).toHaveValue('Test');
  await expect(page.locator('#active')).toHaveText('Yes');
});

test('Admin Users Show Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await adminPage.waitForPage.adminPage();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  // Click first office user row
  await page.locator('tbody >> tr').first().click();
  await adminPage.waitForPage.adminPage();

  // Get first id field and check that it's in the URL
  const id = await page.locator('.ra-field-id > span').first().textContent();
  expect(page.url()).toContain(id);

  const firstName = await page.locator('.ra-field-firstName > span').textContent();
  const lastName = await page.locator('.ra-field-lastName > span').textContent();

  await expect(page.getByRole('heading', { name: `${firstName} ${lastName}` })).toBeVisible();

  const labels = [
    'Id',
    'Email',
    'First name',
    'Last name',
    'User Id',
    'Organization Id',
    'Active',
    'Created at',
    'Updated at',
  ];

  await adminPage.expectLabels(labels);
});

test('Admin Users Edit Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  // create a new admin user to edit
  // using an existing one may stop on a concurrent playwright session
  const adminUser = await adminPage.testHarness.buildDefaultAdminUser();
  const adminUserId = adminUser.id;

  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await adminPage.waitForPage.adminPage();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  // go directly to the show page for the admin user because if we
  // have created more than 25 admin users, the user may not be on the
  // first page and we don't have a search functionality for admin users
  await page.goto(`/system/admin-users/${adminUserId}/show`);
  await adminPage.waitForPage.adminPage();

  await page.getByRole('link', { name: 'Edit' }).click();
  await adminPage.waitForPage.adminPage();

  // Potentially flaky if there are multiple pages of admin users
  expect(page.url()).toContain(adminUserId);

  const disabledFields = ['id', 'userId', 'createdAt', 'updatedAt'];
  for (const field of disabledFields) {
    await expect(page.locator(`#${field}`)).toBeDisabled();
  }

  const firstName = page.getByLabel('First name');
  await firstName.focus();
  await firstName.clear();
  await firstName.fill('NewFirst');

  const lastName = page.getByLabel('Last name');
  await lastName.focus();
  await lastName.clear();
  await lastName.fill('NewLast');

  // set the user to the active status they did NOT have before
  const activeStatus = await page.locator('div:has(label :text-is("Active")) >> input[name="active"]').inputValue();

  const newStatus = (activeStatus !== 'true').toString();
  await page.locator('div:has(label :text-is("Active")) >> #active').click();
  await page.locator(`ul[aria-labelledby="active-label"] >> li[data-value="${newStatus}"]`).click();

  // Save updates
  await page.getByRole('button', { name: 'Save' }).click();

  await adminPage.waitForPage.adminPage();

  // back to list of all users
  expect(page.url()).not.toContain(adminUserId);

  // look at admin user's page to ensure changes were saved
  // Potentially flaky if there are multiple pages of admin users
  await page.getByText(adminUserId).click();
  await adminPage.waitForPage.adminPage();

  await expect(page.locator('.ra-field-firstName > span')).toHaveText('NewFirst');
  await expect(page.locator('.ra-field-lastName > span')).toHaveText('NewLast');
  await expect(page.getByTestId(newStatus)).toBeVisible();

  // go back to list of all users and ensure updated
  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await adminPage.waitForPage.adminPage();
  await expect(page.locator(`tr:has(:text("${adminUserId}")) >> td.column-firstName`)).toHaveText('NewFirst');
  await expect(page.locator(`tr:has(:text("${adminUserId}")) >> td.column-lastName`)).toHaveText('NewLast');
  await expect(page.locator(`tr:has(:text("${adminUserId}")) >> td.column-active >> svg`)).toHaveAttribute(
    'data-testid',
    newStatus,
  );
});
