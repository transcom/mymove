// @ts-check
const { test, expect } = require('@playwright/test');

const { signInAsNewAdminUser } = require('../utils/signIn');

test('Admin Users List Page', async ({ page }) => {
  await page.goto('/');
  await signInAsNewAdminUser(page);

  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  const columnLabels = ['Id', 'Email', 'First name', 'Last name', 'User Id', 'Active'];

  for (const label of columnLabels) {
    await expect(page.getByRole('columnheader').getByText(label, { exact: true })).toBeVisible();
  }
});

test('Admin User Create Page', async ({ page }) => {
  await page.goto('/');
  await signInAsNewAdminUser(page);

  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  await page.getByRole('button', { name: 'Create' }).click();
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

  // redirected to edit details page
  const adminUserID = await page.locator('#id').inputValue();

  await expect(page.getByRole('heading', { name: `Admin Users #${adminUserID}` })).toBeVisible();

  await expect(page.locator('#email')).toHaveValue(testEmail);
  await expect(page.locator('#firstName')).toHaveValue('Cypress');
  await expect(page.locator('#lastName')).toHaveValue('Test');
  await expect(page.locator('#active')).toHaveText('Yes');
});

test('Admin Users Show Page', async ({ page }) => {
  await page.goto('/');
  await signInAsNewAdminUser(page);

  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  // click on first row
  await page.locator('tr[resource="admin_users"]').first().click();

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

  for (const label of labels) {
    await expect(page.locator('label').getByText(label, { exact: true })).toBeVisible();
  }
});

test('Admin Users Edit Page', async ({ page }) => {
  await page.goto('/');
  await signInAsNewAdminUser(page);

  await page.getByRole('menuitem', { name: 'Admin Users' }).click();
  await expect(page.getByRole('heading', { name: 'Admin Users' })).toBeVisible();

  // click on first row
  await page.locator('tr[resource="admin_users"]').first().click();

  const id = await page.locator('div:has(label :text-is("Id")) > div > span').textContent();

  await page.getByRole('button', { name: 'Edit' }).click();
  expect(page.url()).toContain(id);

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

  // back to list of all users
  expect(page.url()).not.toContain(id);

  await expect(page.locator(`tr:has(:text("${id}")) >> td.column-active >> svg`)).toHaveAttribute(
    'data-testid',
    newStatus,
  );

  await expect(page.locator(`tr:has(:text("${id}")) >> td.column-firstName`)).toHaveText('NewFirst');
  await expect(page.locator(`tr:has(:text("${id}")) >> td.column-lastName`)).toHaveText('NewLast');
});
