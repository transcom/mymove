/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../utils/admin/adminTest';

/**
 * @param {import('@playwright/test').Page} page
 * @param {string} email
 */
async function searchForOfficeUser(page, email) {
  await page.getByLabel('Search').click();
  await page.getByLabel('Search').fill(email);
  await page.getByLabel('Search').press('Enter');
}

test.describe('Office Users List Page', () => {
  test('successfully navigates to Office Users page', async ({ page, adminPage }) => {
    await adminPage.testHarness.buildOfficeUserWithTOOAndTIO();
    await adminPage.signInAsNewAdminUser();

    // we should be at the office users list page by default,
    // but let's go somewhere else and then come back to make sure the side nav link works
    await page.getByRole('menuitem', { name: 'Moves' }).click();
    expect(page.url()).toContain('/system/moves');

    // now we'll come back to the office users page:
    await page.getByRole('menuitem', { name: 'Office Users', exact: true }).click();
    expect(page.url()).toContain('/system/office-users');
    await expect(page.locator('header')).toContainText('Office Users');

    const columnLabels = ['Id', 'Email', 'First name', 'Last name', 'Transportation Office', 'User Id', 'Active'];

    await adminPage.expectRoleLabelsByText('columnheader', columnLabels);
  });
});

test.describe('Office User Create Page', () => {
  test('pulls up create page for an office user', async ({ page, adminPage }) => {
    await adminPage.signInAsNewAdminUser();
    // we tested the side nav in the previous test,
    // so let's work with the assumption that we were already redirected to this page:
    expect(page.url()).toContain('/system/requested-office-users');
    await page.getByRole('menuitem', { name: 'Office Users', exact: true }).click();
    expect(page.url()).toContain('/system/office-users');

    await page.getByRole('link', { name: 'Create' }).click();
    await expect(page.getByRole('heading', { name: 'Create Office Users', exact: true })).toBeVisible();

    expect(page.url()).toContain('/system/office-users/create');

    // we need to add the date to the email so that it is unique every time (only one record per email allowed in db)
    const testEmail = `cy.admin_user.${Date.now()}@example.com`;

    // create an office user
    const firstName = page.getByLabel('First name');
    await firstName.focus();
    await firstName.fill('Cypress');

    // middle initials are not working for me in chrome outside of
    // playwright - ahobson 2022-12-29
    // await page.getByLabel('Middle initials').fill('X');
    const lastName = page.getByLabel('Last name');
    await lastName.focus();
    await lastName.fill('Test');

    const email = page.getByLabel('Email');
    await email.focus();
    await email.fill(testEmail);

    const phone = page.getByLabel('Telephone');
    await phone.focus();
    await phone.fill('222-555-1234');

    await page.getByLabel('Services Counselor').click();
    await page.getByLabel('Supervisor').click();
    // The autocomplete form results in multiple matching elements, so
    // pick the input element
    await page.getByLabel('Transportation Office').fill('JPPSO Testy McTest');
    // the autocomplete might return multiples because of concurrent
    // tests running that are adding offices
    await page.getByRole('option', { name: 'JPPSO Testy McTest' }).first().click();

    await page.getByRole('button', { name: 'Save' }).click();
    await adminPage.waitForPage.adminPage();

    // redirected to edit details page
    const officeUserID = await page.locator('#id').inputValue();

    await expect(page.getByRole('heading', { name: `Office Users #${officeUserID}` })).toBeVisible();

    await expect(page.locator('#email')).toHaveValue(testEmail);
    await expect(page.locator('#firstName')).toHaveValue('Cypress');
    await expect(page.locator('#lastName')).toHaveValue('Test');
    // middle initials are not working for me in chrome outside of
    // playwright - ahobson 2022-12-29
    // await expect(page.locator('#middleInitials')).toHaveValue('X');
    await expect(page.locator('#telephone')).toHaveValue('222-555-1234');
    await expect(page.locator('#active')).toHaveText('Yes');
  });
});

test.describe('Office Users Show Page', () => {
  test('pulls up details page for an office user', async ({ page, adminPage }) => {
    await adminPage.testHarness.buildOfficeUserWithTOOAndTIO();
    await adminPage.signInAsNewAdminUser();

    expect(page.url()).toContain('/system/requested-office-users');
    await adminPage.waitForPage.adminPage();
    await page.getByRole('menuitem', { name: 'Office Users', exact: true }).click();
    expect(page.url()).toContain('/system/office-users');

    // Click first office user row
    await page.locator('tbody >> tr').first().click();

    // Get first id field and check that it's in the URL
    const id = await page.locator('.ra-field-id > span').first().textContent();
    expect(page.url()).toContain(id);

    // check that the office user's name is shown in the page title
    const firstName = await page.locator('.ra-field-firstName > span').textContent();
    const lastName = await page.locator('.ra-field-lastName > span').textContent();

    await expect(page.getByRole('heading', { name: `${firstName} ${lastName}` })).toBeVisible();

    const labels = [
      'Id',
      'User Id',
      'Email',
      'First name',
      'Middle initials',
      'Last name',
      'Telephone',
      'Active',
      'Transportation Office',
      'Created at',
      'Updated at',
    ];
    await adminPage.expectLabels(labels);
  });
});

test.describe('Office Users Edit Page', () => {
  test('pulls up edit page for an office user', async ({ page, adminPage }) => {
    const officeUser = await adminPage.testHarness.buildOfficeUserWithTOOAndTIO();
    const email = officeUser.okta_email;

    await adminPage.signInAsNewAdminUser();

    expect(page.url()).toContain('/system/requested-office-users');
    await page.getByRole('menuitem', { name: 'Office Users', exact: true }).click();
    expect(page.url()).toContain('/system/office-users');
    await searchForOfficeUser(page, email);
    await page.getByText(email).click();
    await adminPage.waitForPage.adminPage();

    await page.getByRole('link', { name: 'Edit' }).click();
    await adminPage.waitForPage.adminPage();

    const disabledFields = ['id', 'email', 'userId', 'createdAt', 'updatedAt'];
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

    // The autocomplete form results in multiple matching elements, so
    // pick the input element
    await expect(page.getByLabel('Transportation Office')).toBeEditable();

    // set the user to the active status they did NOT have before
    const activeStatus = await page.locator('div:has(label :text-is("Active")) >> input[name="active"]').inputValue();

    const newStatus = (activeStatus !== 'true').toString();
    await page.locator('div:has(label :text-is("Active")) >> #active').click();
    await page.locator(`ul[aria-labelledby="active-label"] >> li[data-value="${newStatus}"]`).click();

    const tooCheckbox = page.getByLabel('Task Ordering Officer');
    const tioCheckbox = page.getByLabel('Task Invoicing Officer');

    if (tioCheckbox.isChecked() && tooCheckbox.isChecked()) {
      await tooCheckbox.click();
    }

    await page.getByRole('button', { name: 'Save' }).click();
    await adminPage.waitForPage.adminPage();

    await searchForOfficeUser(page, email);
    await expect(page.locator(`tr:has(:text("${email}")) >> td.column-active >> svg`)).toHaveAttribute(
      'data-testid',
      newStatus,
    );

    await expect(page.locator(`tr:has(:text("${email}")) >> td.column-firstName`)).toHaveText('NewFirst');
    await expect(page.locator(`tr:has(:text("${email}")) >> td.column-lastName`)).toHaveText('NewLast');
  });

  test('prevents safety move priv selection with Customer role', async ({ page, adminPage }) => {
    const officeUser = await adminPage.testHarness.buildOfficeUserWithCustomer();
    const email = officeUser.okta_email;

    // create a new admin user to edit
    // using an existing one may stop on a concurrent playwright session
    const adminUser = await adminPage.testHarness.buildDefaultSuperAdminUser();
    await adminPage.signInAsExistingAdminUser(adminUser.user_id);

    expect(page.url()).toContain('/system/requested-office-users');
    await page.getByRole('menuitem', { name: 'Office Users', exact: true }).click();
    expect(page.url()).toContain('/system/office-users');
    await searchForOfficeUser(page, email);
    await page.getByText(email).click();
    await adminPage.waitForPage.adminPage();

    await page.getByRole('link', { name: 'Edit' }).click();
    await adminPage.waitForPage.adminPage();

    const safetyMoveCheckbox = page.getByLabel('Safety Moves');
    const customerCheckbox = page.getByLabel('Customer', { exact: true });

    await expect(customerCheckbox).toBeChecked();
    await safetyMoveCheckbox.click();
    await expect(safetyMoveCheckbox).not.toBeChecked();
  });

  test('prevents safety move priv selection with Contracting Officer role', async ({ page, adminPage }) => {
    const officeUser = await adminPage.testHarness.buildOfficeUserWithContractingOfficer();
    const email = officeUser.okta_email;

    // create a new admin user to edit
    // using an existing one may stop on a concurrent playwright session
    const adminUser = await adminPage.testHarness.buildDefaultSuperAdminUser();
    await adminPage.signInAsExistingAdminUser(adminUser.user_id);

    expect(page.url()).toContain('/system/requested-office-users');
    await page.getByRole('menuitem', { name: 'Office Users', exact: true }).click();
    expect(page.url()).toContain('/system/office-users');
    await searchForOfficeUser(page, email);
    await page.getByText(email).click();
    await adminPage.waitForPage.adminPage();

    await page.getByRole('link', { name: 'Edit' }).click();
    await adminPage.waitForPage.adminPage();

    const safetyMoveCheckbox = page.getByLabel('Safety Moves');
    const customerCheckbox = page.getByLabel('Contracting Officer', { exact: true });

    await expect(customerCheckbox).toBeChecked();
    await safetyMoveCheckbox.click();
    await expect(safetyMoveCheckbox).not.toBeChecked();
  });

  test('prevents safety move priv selection with Prime Simulator role', async ({ page, adminPage }) => {
    const officeUser = await adminPage.testHarness.buildOfficeUserWithPrimeSimulator();
    const email = officeUser.okta_email;

    // create a new admin user to edit
    // using an existing one may stop on a concurrent playwright session
    const adminUser = await adminPage.testHarness.buildDefaultSuperAdminUser();
    await adminPage.signInAsExistingAdminUser(adminUser.user_id);

    expect(page.url()).toContain('/system/requested-office-users');
    await page.getByRole('menuitem', { name: 'Office Users', exact: true }).click();
    expect(page.url()).toContain('/system/office-users');
    await searchForOfficeUser(page, email);
    await page.getByText(email).click();
    await adminPage.waitForPage.adminPage();

    await page.getByRole('link', { name: 'Edit' }).click();
    await adminPage.waitForPage.adminPage();

    const safetyMoveCheckbox = page.getByLabel('Safety Moves');
    const customerCheckbox = page.getByLabel('Prime Simulator', { exact: true });

    await expect(customerCheckbox).toBeChecked();
    await safetyMoveCheckbox.click();
    await expect(safetyMoveCheckbox).not.toBeChecked();
  });

  test('prevents safety move priv selection with Government Surveillance Representative role', async ({
    page,
    adminPage,
  }) => {
    const officeUser = await adminPage.testHarness.buildOfficeUserWithGSR();
    const email = officeUser.okta_email;

    // create a new admin user to edit
    // using an existing one may stop on a concurrent playwright session
    const adminUser = await adminPage.testHarness.buildDefaultSuperAdminUser();
    await adminPage.signInAsExistingAdminUser(adminUser.user_id);

    expect(page.url()).toContain('/system/requested-office-users');
    await page.getByRole('menuitem', { name: 'Office Users', exact: true }).click();
    expect(page.url()).toContain('/system/office-users');
    await searchForOfficeUser(page, email);
    await page.getByText(email).click();
    await adminPage.waitForPage.adminPage();

    await page.getByRole('link', { name: 'Edit' }).click();
    await adminPage.waitForPage.adminPage();

    const safetyMoveCheckbox = page.getByLabel('Safety Moves');
    const customerCheckbox = page.getByLabel('Government Surveillance Representative', { exact: true });

    await expect(customerCheckbox).toBeChecked();
    await safetyMoveCheckbox.click();
    await expect(safetyMoveCheckbox).not.toBeChecked();
  });
});
