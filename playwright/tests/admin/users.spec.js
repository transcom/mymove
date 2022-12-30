// @ts-check
const { test, expect } = require('../utils/adminTest');

test.describe('Users Page', () => {
  test('successfully navigates to users page', async ({ page, adminPage }) => {
    await adminPage.signInAsNewAdminUser();

    await page.getByRole('menuitem', { name: 'Users', exact: true }).click();
    await adminPage.waitForAdminPageToLoad();
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
    const officeUser = await adminPage.buildOfficeUserWithTOOAndTIO();
    const email = officeUser.login_gov_email;
    await adminPage.signInAsNewAdminUser();

    await page.getByRole('menuitem', { name: 'Users', exact: true }).click();
    await adminPage.waitForAdminPageToLoad();
    expect(page.url()).toContain('/system/users');

    await page.getByText(email).click();
    await adminPage.waitForAdminPageToLoad();

    // check that the user's login.gov email is shown in the page title
    await expect(page.getByRole('heading', { name: email })).toBeVisible();

    const labels = ['User ID', 'User email', 'Active', 'Created at', 'Updated at'];
    for (const label of labels) {
      await expect(page.locator('label').getByText(label, { exact: true })).toBeVisible();
    }
  });
});

test.describe('Users Details Edit Page', () => {
  test('pulls up edit page for a user', async ({ page, adminPage }) => {
    const officeUser = await adminPage.buildOfficeUserWithTOOAndTIO();
    const email = officeUser.login_gov_email;
    await adminPage.signInAsNewAdminUser();

    await page.getByRole('menuitem', { name: 'Users', exact: true }).click();
    await adminPage.waitForAdminPageToLoad();
    expect(page.url()).toContain('/system/users');

    await page.getByText(email).click();
    await adminPage.waitForAdminPageToLoad();

    await page.getByRole('button', { name: 'Edit' }).click();
    await adminPage.waitForAdminPageToLoad();

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
    await adminPage.waitForAdminPageToLoad();

    // check that user was deactivated
    await expect(page.locator(`tr:has(:text("${email}")) >> td.column-active >> svg`)).toHaveAttribute(
      'data-testid',
      newStatus,
    );
  });
});
