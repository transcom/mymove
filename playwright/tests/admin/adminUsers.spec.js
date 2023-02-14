/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../utils/adminTest';

test('Admin Users List Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await adminPage.waitForAdminPageToLoad();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  const columnLabels = ['Id', 'Email', 'First name', 'Last name', 'User Id', 'Active'];

  await adminPage.expectRoleLabelsByText('columnheader', columnLabels);
});

test('Admin User Create Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await adminPage.waitForAdminPageToLoad();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  await page.getByRole('button', { name: 'Create' }).click();
  await adminPage.waitForAdminPageToLoad();
  await expect(page.getByRole('heading', { name: 'Create Admin Users' })).toBeVisible();

  // we need to add the date to the email so that it is unique every time (only one record per email allowed in db)
  const testEmail = `cy.admin_user.${Date.now()}@example.com`;

  // create an admin user
  await page.getByLabel('Email').fill(testEmail);
  await page.getByLabel('First name').fill('Cypress');
  await page.getByLabel('Last name').fill('Test');

  await page.getByLabel('Organization').click();
  await page.getByRole('option').first().click();
  await page.getByRole('button').filter({ hasText: 'Save' }).click();
  await adminPage.waitForAdminPageToLoad();

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
  await adminPage.waitForAdminPageToLoad();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  // click on first row
  await page.locator('tr[resource="admin_users"]').first().click();
  await adminPage.waitForAdminPageToLoad();

  const id = await page.locator('div:has(label :text-is("Id")) > div > span').textContent();
  expect(page.url()).toContain(id);

  const firstName = await page.locator('label:has-text("First name") + div').textContent();
  const lastName = await page.locator('label:has-text("Last name") + div').textContent();

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

  await adminPage.expectLocatorLabelsByText('label', labels, { exact: true });
});

test('Admin Users Edit Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  // create a new admin user to edit
  // using an existing one may stop on a concurrent playwright session
  const adminUser = await adminPage.testHarness.buildDefaultAdminUser();
  const adminUserId = adminUser.id;

  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await adminPage.waitForAdminPageToLoad();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  // go directly to the show page for the admin user because if we
  // have created more than 25 admin users, the user may not be on the
  // first page and we don't have a search functionality for admin users
  await page.goto(`/system/admin_users/${adminUserId}/show`);
  await adminPage.waitForAdminPageToLoad();

  await page.getByRole('button', { name: 'Edit' }).click();
  await adminPage.waitForAdminPageToLoad();
  expect(page.url()).toContain(adminUserId);

  const disabledFields = ['id', 'email', 'userId', 'createdAt', 'updatedAt'];
  for (const field of disabledFields) {
    await expect(page.locator(`#${field}`)).toBeDisabled();
  }

  await page.getByLabel('First name').clear();
  await page.getByLabel('First name').fill('NewFirst');

  await page.getByLabel('Last name').clear();
  await page.getByLabel('Last name').fill('NewLast');
  // set the user to the active status they did NOT have before
  const activeStatus = await page.locator('div:has(label :text-is("Active")) >> input[name="active"]').inputValue();

  const newStatus = (activeStatus !== 'true').toString();
  await page.locator('div:has(label :text-is("Active")) >> #active').click();
  await page.locator(`ul[aria-labelledby="active-label"] >> li[data-value="${newStatus}"]`).click();

  await page.getByRole('button', { name: 'Save' }).click();
  await adminPage.waitForAdminPageToLoad();

  // back to list of all users
  expect(page.url()).not.toContain(adminUserId);

  await expect(page.locator(`tr:has(:text("${adminUserId}")) >> td.column-active >> svg`)).toHaveAttribute(
    'data-testid',
    newStatus,
  );

  await expect(page.locator(`tr:has(:text("${adminUserId}")) >> td.column-firstName`)).toHaveText('NewFirst');
  await expect(page.locator(`tr:has(:text("${adminUserId}")) >> td.column-lastName`)).toHaveText('NewLast');
});
