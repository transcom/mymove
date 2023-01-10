// @ts-check
const { test, expect } = require('../utils/adminTest');

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
    await page.getByRole('menuitem', { name: 'Office users' }).click();
    expect(page.url()).toContain('/system/office_users');
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
    expect(page.url()).toContain('/system/office_users');

    await page.getByRole('button', { name: 'Create' }).click();
    await expect(page.getByRole('heading', { name: 'Create Office Users' })).toBeVisible();

    expect(page.url()).toContain('/system/office_users/create');

    // we need to add the date to the email so that it is unique every time (only one record per email allowed in db)
    const testEmail = `cy.admin_user.${Date.now()}@example.com`;

    // create an office user
    await page.getByLabel('First name').fill('Cypress');

    // middle initials are not working for me in chrome outside of
    // playwright - ahobson 2022-12-29
    // await page.getByLabel('Middle initials').fill('X');
    await page.getByLabel('Last name').fill('Test');

    await page.getByLabel('Email').fill(testEmail);
    await page.getByLabel('Telephone').fill('222-555-1234');
    await page.getByText('Services Counselor').click();
    await page.getByLabel('Transportation Office').fill('JPPSO Testy McTest');
    // the autocomplete might return multiples because of concurrent
    // tests running that are adding offices
    await page.getByRole('option', { name: 'JPPSOTestyMcTest' }).first().click();

    await page.getByRole('button', { name: 'Save' }).click();
    await adminPage.waitForAdminPageToLoad();

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

    expect(page.url()).toContain('/system/office_users');

    await page.locator('tr[resource="office_users"]').first().click();

    // check that the office user's name is shown in the page title
    const id = await page.locator('div:has(label :text-is("Id")) > div > span').textContent();
    expect(page.url()).toContain(id);

    const firstName = await page.locator('label:has-text("First name") + div').textContent();
    const lastName = await page.locator('label:has-text("Last name") + div').textContent();

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
      'Roles',
      'Transportation Office',
      'Created at',
      'Updated at',
    ];
    await adminPage.expectLocatorLabelsByText('label', labels, { exact: true });
  });
});

test.describe('Office Users Edit Page', () => {
  test('pulls up edit page for an office user', async ({ page, adminPage }) => {
    const officeUser = await adminPage.testHarness.buildOfficeUserWithTOOAndTIO();
    const email = officeUser.login_gov_email;

    await adminPage.signInAsNewAdminUser();

    expect(page.url()).toContain('/system/office_users');
    await searchForOfficeUser(page, email);
    await page.getByText(email).click();
    await adminPage.waitForAdminPageToLoad();

    await page.getByRole('button', { name: 'Edit' }).click();
    await adminPage.waitForAdminPageToLoad();

    const disabledFields = ['id', 'email', 'userId', 'createdAt', 'updatedAt'];
    for (const field of disabledFields) {
      await expect(page.locator(`#${field}`)).toBeDisabled();
    }

    await page.getByLabel('First name').clear();
    await page.getByLabel('First name').fill('NewFirst');

    await page.getByLabel('Last name').clear();
    await page.getByLabel('Last name').fill('NewLast');

    await expect(page.getByLabel('Transportation Office')).toBeEditable();

    // set the user to the active status they did NOT have before
    const activeStatus = await page.locator('div:has(label :text-is("Active")) >> input[name="active"]').inputValue();

    const newStatus = (activeStatus !== 'true').toString();
    await page.locator('div:has(label :text-is("Active")) >> #active').click();
    await page.locator(`ul[aria-labelledby="active-label"] >> li[data-value="${newStatus}"]`).click();

    await page.getByRole('button', { name: 'Save' }).click();
    await adminPage.waitForAdminPageToLoad();

    await searchForOfficeUser(page, email);
    await expect(page.locator(`tr:has(:text("${email}")) >> td.column-active >> svg`)).toHaveAttribute(
      'data-testid',
      newStatus,
    );

    await expect(page.locator(`tr:has(:text("${email}")) >> td.column-firstName`)).toHaveText('NewFirst');
    await expect(page.locator(`tr:has(:text("${email}")) >> td.column-lastName`)).toHaveText('NewLast');
  });
});
