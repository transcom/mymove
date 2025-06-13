import { test, expect } from '../utils/admin/adminTest';

const requestAccountPrivileges = process.env.FEATURE_FLAG_REQUEST_ACCOUNT_PRIVILEGES;
// Helper function to search for a user by email
async function searchForUserByEmail(page, email) {
  await page.getByPlaceholder('Search').click();
  await page.getByPlaceholder('Search').fill(email);
  await page.getByPlaceholder('Search').press('Enter');
}

// Helper function to verify no privileges assigned for an office user
async function verifyNoPrivilegesAssigned(page, adminPage, email) {
  await page.goto('/system/office-users');
  await searchForUserByEmail(page, email);
  await page.getByText(email).click();
  await adminPage.waitForPage.adminPage();
  await expect(page.getByText(/No privileges assigned to this office user/i)).toBeVisible();
}

// Helper function to verify privilege assigned for an office user
async function verifyPrivilegeAssigned(page, adminPage, email, privilege = 'supervisor') {
  await page.goto('/system/office-users');
  await searchForUserByEmail(page, email);
  await page.getByText(email).click();
  await adminPage.waitForPage.adminPage();
  await expect(page.getByText(/privilege name/i)).toBeVisible();
  // Use a case-insensitive string match for privilege
  await expect(page.getByText(privilege, { exact: false })).toBeVisible();
}

test.describe('RequestedOfficeUserShow', () => {
  test('approve requested office user with out privilege ', async ({ page, adminPage }) => {
    // Build a requested office user with a privilege request
    const officeUser = await adminPage.testHarness.buildRequestedOfficeUser();
    const testEmail = officeUser.okta_email;

    // Sign in as a new admin user
    await adminPage.signInAsNewAdminUser();

    // Navigate to the requested office users admin page
    await page.goto('/system/requested-office-users');

    // Search for the test office user by email
    await searchForUserByEmail(page, testEmail);

    // Click on the user's email to go to their detail page
    await page.getByText(testEmail).click();
    await adminPage.waitForPage.adminPage();

    // Click the Approve button to open the privilege confirm dialog
    await page
      .getByRole('button', { name: /approve/i })
      .first()
      .click();

    // Assert the privilege confirm dialog is closed
    await expect(page.locator('[data-testid="RequestedOfficeUserPrivilegeConfirm"]')).toBeHidden();

    // Use helper to verify no privileges assigned
    await verifyNoPrivilegesAssigned(page, adminPage, testEmail);
  });

  test('approve and verify rejected requested privilege not assigned to office user', async ({ page, adminPage }) => {
    test.skip(requestAccountPrivileges === 'false', 'Skip if request account privileges feature flag is disabled');
    // Build a requested office user with a privilege request
    const officeUser = await adminPage.testHarness.buildRequestedOfficeUserWithPrivilege();
    const testEmail = officeUser.okta_email;

    // Sign in as a new admin user
    await adminPage.signInAsNewAdminUser();

    // Navigate to the requested office users admin page
    await page.goto('/system/requested-office-users');

    // Search for the test office user by email
    await searchForUserByEmail(page, testEmail);

    // Click on the user's email to go to their detail page
    await page.getByText(testEmail).click();
    await adminPage.waitForPage.adminPage();

    // Click the Approve button to open the privilege confirm dialog
    await page
      .getByRole('button', { name: /approve/i })
      .first()
      .click();

    // Assert the privilege confirm dialog is visible
    await expect(page.locator('[data-testid="RequestedOfficeUserPrivilegeConfirm"]')).toBeVisible();

    // Find the checkbox labeled 'Supervisor' and uncheck it if checked
    const supervisorCheckbox = page.getByRole('checkbox', { name: /supervisor/i });
    if (await supervisorCheckbox.isChecked()) {
      await supervisorCheckbox.uncheck();
    }

    // Click the Confirm button to approve the privilege (with checkbox unchecked)
    await page.getByRole('button', { name: /confirm/i }).click();

    await expect(page.locator('[data-testid="RequestedOfficeUserPrivilegeConfirm"]')).toBeHidden();
    await verifyNoPrivilegesAssigned(page, adminPage, testEmail);
  });

  test('approve and verify privilege assigned to office user', async ({ page, adminPage }) => {
    test.skip(requestAccountPrivileges === 'false', 'Skip if request account privileges feature flag is disabled');
    // Build a requested office user with a privilege request
    const officeUser = await adminPage.testHarness.buildRequestedOfficeUserWithPrivilege();
    const testEmail = officeUser.okta_email;

    // Sign in as a new admin user
    await adminPage.signInAsNewAdminUser();

    // Navigate to the requested office users admin page
    await page.goto('/system/requested-office-users');

    // Search for the test office user by email
    await searchForUserByEmail(page, testEmail);

    // Click on the user's email to go to their detail page
    await page.getByText(testEmail).click();
    await adminPage.waitForPage.adminPage();

    // Click the Approve button to open the privilege confirm dialog
    await page
      .getByRole('button', { name: /approve/i })
      .first()
      .click();

    // Assert the privilege confirm dialog is visible
    await expect(page.locator('[data-testid="RequestedOfficeUserPrivilegeConfirm"]')).toBeVisible();

    // Find the checkbox labeled 'Supervisor' and check it if not already checked
    const supervisorCheckbox = page.getByRole('checkbox', { name: /supervisor/i });
    await expect(supervisorCheckbox).toBeVisible();
    await expect(await supervisorCheckbox.isChecked()).toBeTruthy();

    // Click the Confirm button to approve the privilege
    await page.getByRole('button', { name: /confirm/i }).click();

    // Assert the privilege confirm dialog is closed
    await expect(page.locator('[data-testid="RequestedOfficeUserPrivilegeConfirm"]')).toBeHidden();

    // Use helper to verify privilege assigned
    await verifyPrivilegeAssigned(page, adminPage, testEmail);
  });
});

test.describe('RequestedOfficeUserEdit', () => {
  test('assign privilege and approve a requested office user who initially has no privileges', async ({
    page,
    adminPage,
  }) => {
    // Build a requested office user with no privileges
    const officeUser = await adminPage.testHarness.buildRequestedOfficeUser();
    const testEmail = officeUser.okta_email;

    // Sign in as a new admin user
    await adminPage.signInAsNewAdminUser();

    // Go to the requested office users admin page
    await page.goto('/system/requested-office-users');

    // Search for the office user by email
    await searchForUserByEmail(page, testEmail);

    // Open the user's detail page
    await page.getByText(testEmail).click();
    await adminPage.waitForPage.adminPage();

    // Open the Edit page for the user
    await page.getByRole('link', { name: 'Edit' }).first().click();
    await adminPage.waitForPage.adminPage();

    // Ensure Supervisor privilege is not checked, then check it
    const supervisorCheckbox = page.getByRole('checkbox', { name: /supervisor/i });
    await expect(supervisorCheckbox).toBeVisible();
    await expect(await supervisorCheckbox.isChecked()).toBeFalsy();
    await supervisorCheckbox.check();

    await page
      .getByRole('button', { name: /approve/i })
      .first()
      .click();

    // Verify Supervisor privilege is assigned
    await verifyPrivilegeAssigned(page, adminPage, testEmail);
  });

  test('approve and verify requested office user with privilege', async ({ page, adminPage }) => {
    test.skip(requestAccountPrivileges === 'false', 'Skip if request account privileges feature flag is disabled');
    // Build a requested office user with Supervisor privilege requested
    const officeUser = await adminPage.testHarness.buildRequestedOfficeUserWithPrivilege();
    const testEmail = officeUser.okta_email;

    // Sign in as a new admin user
    await adminPage.signInAsNewAdminUser();

    // Go to the requested office users admin page
    await page.goto('/system/requested-office-users');

    // Search for the office user by email
    await searchForUserByEmail(page, testEmail);

    // Open the user's detail page
    await page.getByText(testEmail).click();
    await adminPage.waitForPage.adminPage();

    // Open the Edit page for the user
    await page.getByRole('link', { name: 'Edit' }).first().click();
    await adminPage.waitForPage.adminPage();

    // Ensure Supervisor privilege is checked
    const supervisorEditCheckbox = page.getByRole('checkbox', { name: /supervisor/i });
    await expect(supervisorEditCheckbox).toBeVisible();
    await expect(await supervisorEditCheckbox.isChecked()).toBeTruthy();

    // Open the privilege confirm dialog
    await page
      .getByRole('button', { name: /approve/i })
      .first()
      .click();

    // Assert the privilege confirm dialog is visible
    await expect(page.locator('[data-testid="RequestedOfficeUserPrivilegeConfirm"]')).toBeVisible();

    // Ensure Supervisor privilege is checked in the dialog
    const supervisorDialogCheckbox = page.getByRole('checkbox', { name: /supervisor/i });
    if (await !supervisorDialogCheckbox.isChecked()) {
      await supervisorDialogCheckbox.check();
    }

    // Confirm the privilege assignment
    await page.getByRole('button', { name: /confirm/i }).click();

    // Assert the dialog is closed and privilege is assigned
    await expect(page.locator('[data-testid="RequestedOfficeUserPrivilegeConfirm"]')).toBeHidden();
    await verifyPrivilegeAssigned(page, adminPage, testEmail);
  });

  // Reject privilege and approve requested office user who requested a privilege
  test('reject privilege and approve requested office user who requested a privilege', async ({ page, adminPage }) => {
    test.skip(requestAccountPrivileges === 'false', 'Skip if request account privileges feature flag is disabled');
    // Build a requested office user with Supervisor privilege requested
    const officeUser = await adminPage.testHarness.buildRequestedOfficeUserWithPrivilege();
    const testEmail = officeUser.okta_email;

    // Sign in as a new admin user
    await adminPage.signInAsNewAdminUser();

    // Go to the requested office users admin page
    await page.goto('/system/requested-office-users');

    // Search for the office user by email
    await searchForUserByEmail(page, testEmail);

    // Open the user's detail page
    await page.getByText(testEmail).click();
    await adminPage.waitForPage.adminPage();

    // Open the Edit page for the user
    await page.getByRole('link', { name: 'Edit' }).first().click();
    await adminPage.waitForPage.adminPage();

    // Open the privilege confirm dialog
    await page
      .getByRole('button', { name: /approve/i })
      .first()
      .click();

    // Assert the privilege confirm dialog is visible
    await expect(page.locator('[data-testid="RequestedOfficeUserPrivilegeConfirm"]')).toBeVisible();

    // Uncheck Supervisor privilege in the dialog if checked
    const supervisorCheckbox = page.getByRole('checkbox', { name: /supervisor/i });
    if (await supervisorCheckbox.isChecked()) {
      await supervisorCheckbox.uncheck();
    }

    // Confirm the privilege rejection
    await page.getByRole('button', { name: /confirm/i }).click();

    // Assert the dialog is closed and no privileges are assigned
    await expect(page.locator('[data-testid="RequestedOfficeUserPrivilegeConfirm"]')).toBeHidden();
    await verifyNoPrivilegesAssigned(page, adminPage, testEmail);
  });
});
