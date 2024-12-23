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

    const columnLabels = [
      'Id',
      'Email',
      'First name',
      'Last name',
      'Primary Transportation Office',
      'User Id',
      'Active',
    ];

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
    await page.getByLabel('Transportation Office').nth(0).fill('PPPO Scott AFB - USAF');
    // the autocomplete might return multiples because of concurrent
    // tests running that are adding offices
    await page.getByRole('option', { name: 'PPPO Scott AFB - USAF' }).first().click();

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

  test('has correct supervisor role permissions', async ({ page, adminPage }) => {
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

    const lastName = page.getByLabel('Last name');
    await lastName.focus();
    await lastName.fill('Test');

    const email = page.getByLabel('Email');
    await email.focus();
    await email.fill(testEmail);

    const phone = page.getByLabel('Telephone');
    await phone.focus();
    await phone.fill('222-555-1234');

    // Define constants for all roles checkboxes to be tested
    const customerCheckbox = page.getByLabel('Customer', { exact: true });
    const contractingOfficerCheckbox = page.getByLabel('Contracting Officer', { exact: true });
    const servicesCounselorCheckbox = page.getByLabel('Services Counselor', { exact: true });
    const primeSimulatorCheckbox = page.getByLabel('Prime Simulator', { exact: true });
    const qualityAssuranceEvaluatorCheckbox = page.getByLabel('Quality Assurance Evaluator', { exact: true });
    const customerServiceRepersentativeCheckbox = page.getByLabel('Customer Service Representative', { exact: true });
    const governmentSurveillanceRepresentativeCheckbox = page.getByLabel('Government Surveillance Representative', {
      exact: true,
    });
    const headquartersCheckbox = page.getByLabel('Headquarters', { exact: true });
    const taskOrderingOfficerCheckbox = page.getByLabel('Task Ordering Officer', { exact: true });
    const taskInvoicingOfficerCheckbox = page.getByLabel('Task Invoicing Officer', { exact: true });

    // Define constants for privileges
    const supervisorCheckbox = page.getByLabel('Supervisor', { exact: true });

    // Check roles that cannot have supervisor priveleges
    await customerCheckbox.click();
    await supervisorCheckbox.click();
    await expect(customerCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await customerCheckbox.click();
    await contractingOfficerCheckbox.click();
    await supervisorCheckbox.click();
    await expect(contractingOfficerCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await contractingOfficerCheckbox.click();
    await primeSimulatorCheckbox.click();
    await supervisorCheckbox.click();
    await expect(primeSimulatorCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await primeSimulatorCheckbox.click();
    await qualityAssuranceEvaluatorCheckbox.click();
    await supervisorCheckbox.click();
    await expect(qualityAssuranceEvaluatorCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await qualityAssuranceEvaluatorCheckbox.click();
    await customerServiceRepersentativeCheckbox.click();
    await supervisorCheckbox.click();
    await expect(customerServiceRepersentativeCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await customerServiceRepersentativeCheckbox.click();
    await governmentSurveillanceRepresentativeCheckbox.click();
    await supervisorCheckbox.click();
    await expect(governmentSurveillanceRepresentativeCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await governmentSurveillanceRepresentativeCheckbox.click();
    await headquartersCheckbox.click();
    await supervisorCheckbox.click();
    await expect(headquartersCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await headquartersCheckbox.click();

    // Check roles that can have supervisor priveleges
    await taskOrderingOfficerCheckbox.click();
    await supervisorCheckbox.click();
    await expect(taskOrderingOfficerCheckbox).toBeChecked();
    await expect(supervisorCheckbox).toBeChecked();
    await taskOrderingOfficerCheckbox.click();
    await supervisorCheckbox.click();
    await taskInvoicingOfficerCheckbox.click();
    await supervisorCheckbox.click();
    await expect(taskInvoicingOfficerCheckbox).toBeChecked();
    await expect(supervisorCheckbox).toBeChecked();
    await taskInvoicingOfficerCheckbox.click();
    await supervisorCheckbox.click();
    await servicesCounselorCheckbox.click();
    await supervisorCheckbox.click();
    await expect(servicesCounselorCheckbox).toBeChecked();
    await expect(supervisorCheckbox).toBeChecked();
    // Check selecting roles after having supervisor selected for unallowed roles
    await customerCheckbox.click();
    await expect(customerCheckbox).not.toBeChecked();
    await contractingOfficerCheckbox.click();
    await expect(contractingOfficerCheckbox).not.toBeChecked();
    await primeSimulatorCheckbox.click();
    await expect(primeSimulatorCheckbox).not.toBeChecked();
    await qualityAssuranceEvaluatorCheckbox.click();
    await expect(qualityAssuranceEvaluatorCheckbox).not.toBeChecked();
    await customerServiceRepersentativeCheckbox.click();
    await expect(customerServiceRepersentativeCheckbox).not.toBeChecked();
    await governmentSurveillanceRepresentativeCheckbox.click();
    await expect(governmentSurveillanceRepresentativeCheckbox).not.toBeChecked();
    await headquartersCheckbox.click();
    await expect(headquartersCheckbox).not.toBeChecked();

    // Check selecting roles after having supervisor selected for allowed roles
    await taskOrderingOfficerCheckbox.click();
    await expect(taskOrderingOfficerCheckbox).toBeChecked();
    await taskInvoicingOfficerCheckbox.click();
    await expect(taskInvoicingOfficerCheckbox).toBeChecked();
    await servicesCounselorCheckbox.click();
    await expect(servicesCounselorCheckbox).toBeChecked();

    // Continue test to ensure form still submits with valid information
    await servicesCounselorCheckbox.click();
    await expect(servicesCounselorCheckbox).toBeChecked();

    // The autocomplete form results in multiple matching elements, so
    // pick the input element
    await page.getByLabel('Transportation Office').nth(0).fill('PPPO Scott AFB - USAF');
    // the autocomplete might return multiples because of concurrent
    // tests running that are adding offices
    await page.getByRole('option', { name: 'PPPO Scott AFB - USAF' }).first().click();

    await page.getByRole('button', { name: 'Save' }).click();
    await adminPage.waitForPage.adminPage();

    // redirected to edit details page
    const officeUserID = await page.locator('#id').inputValue();

    await expect(page.getByRole('heading', { name: `Office Users #${officeUserID}` })).toBeVisible();

    await expect(page.locator('#email')).toHaveValue(testEmail);
    await expect(page.locator('#firstName')).toHaveValue('Cypress');
    await expect(page.locator('#lastName')).toHaveValue('Test');
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
      'Transportation Offices',
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
    await expect(page.getByLabel('Transportation Office').nth(0)).toBeEditable();

    // Add a Transportation Office Assignment
    await page.getByTestId('addTransportationOfficeButton').click();
    // n = 2 because of the disabled GBLOC input
    await expect(page.getByLabel('Transportation Office').nth(2)).toBeEditable();
    await page.getByLabel('Transportation Office').nth(2).fill('AGFM');
    // the autocomplete might return multiples because of concurrent
    // tests running that are adding offices
    await page.getByRole('option', { name: 'JPPSO - North East (AGFM) - USAF' }).first().click();
    // Set as primary transportation office
    await page.getByLabel('Primary Office').nth(1).click();
    await page.getByText('You cannot designate more than one primary transportation office.');
    await page.getByLabel('Primary Office').nth(1).click();

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

  test('has correct supervisor role permissions', async ({ page, adminPage }) => {
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
    await expect(page.getByLabel('Transportation Office').nth(0)).toBeEditable();

    // Add a Transportation Office Assignment
    await page.getByTestId('addTransportationOfficeButton').click();
    // n = 2 because of the disabled GBLOC input
    await expect(page.getByLabel('Transportation Office').nth(2)).toBeEditable();
    await page.getByLabel('Transportation Office').nth(2).fill('AGFM');
    // the autocomplete might return multiples because of concurrent
    // tests running that are adding offices
    await page.getByRole('option', { name: 'JPPSO - North East (AGFM) - USAF' }).first().click();
    // Set as primary transportation office
    await page.getByLabel('Primary Office').nth(1).click();
    await page.getByText('You cannot designate more than one primary transportation office.');
    await page.getByLabel('Primary Office').nth(1).click();

    // set the user to the active status they did NOT have before
    const activeStatus = await page.locator('div:has(label :text-is("Active")) >> input[name="active"]').inputValue();

    const newStatus = (activeStatus !== 'true').toString();
    await page.locator('div:has(label :text-is("Active")) >> #active').click();
    await page.locator(`ul[aria-labelledby="active-label"] >> li[data-value="${newStatus}"]`).click();

    // Define constants for all roles checkboxes to be tested
    const customerCheckbox = page.getByLabel('Customer', { exact: true });
    const contractingOfficerCheckbox = page.getByLabel('Contracting Officer', { exact: true });
    const servicesCounselorCheckbox = page.getByLabel('Services Counselor', { exact: true });
    const primeSimulatorCheckbox = page.getByLabel('Prime Simulator', { exact: true });
    const qualityAssuranceEvaluatorCheckbox = page.getByLabel('Quality Assurance Evaluator', { exact: true });
    const customerServiceRepersentativeCheckbox = page.getByLabel('Customer Service Representative', { exact: true });
    const governmentSurveillanceRepresentativeCheckbox = page.getByLabel('Government Surveillance Representative', {
      exact: true,
    });
    const headquartersCheckbox = page.getByLabel('Headquarters', { exact: true });
    const taskOrderingOfficerCheckbox = page.getByLabel('Task Ordering Officer', { exact: true });
    const taskInvoicingOfficerCheckbox = page.getByLabel('Task Invoicing Officer', { exact: true });

    // Define constants for privileges
    const supervisorCheckbox = page.getByLabel('Supervisor', { exact: true });

    // Check roles that cannot have supervisor priveleges
    await customerCheckbox.click();
    await supervisorCheckbox.click();
    await expect(customerCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await customerCheckbox.click();
    await contractingOfficerCheckbox.click();
    await supervisorCheckbox.click();
    await expect(contractingOfficerCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await contractingOfficerCheckbox.click();
    await primeSimulatorCheckbox.click();
    await supervisorCheckbox.click();
    await expect(primeSimulatorCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await primeSimulatorCheckbox.click();
    await qualityAssuranceEvaluatorCheckbox.click();
    await supervisorCheckbox.click();
    await expect(qualityAssuranceEvaluatorCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await qualityAssuranceEvaluatorCheckbox.click();
    await customerServiceRepersentativeCheckbox.click();
    await supervisorCheckbox.click();
    await expect(customerServiceRepersentativeCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await customerServiceRepersentativeCheckbox.click();
    await governmentSurveillanceRepresentativeCheckbox.click();
    await supervisorCheckbox.click();
    await expect(governmentSurveillanceRepresentativeCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await governmentSurveillanceRepresentativeCheckbox.click();
    await headquartersCheckbox.click();
    await supervisorCheckbox.click();
    await expect(headquartersCheckbox).toBeChecked();
    await expect(supervisorCheckbox).not.toBeChecked();
    await headquartersCheckbox.click();

    // Check roles that can have supervisor priveleges
    await taskOrderingOfficerCheckbox.click();
    await supervisorCheckbox.click();
    await expect(taskOrderingOfficerCheckbox).toBeChecked();
    await expect(supervisorCheckbox).toBeChecked();
    await taskOrderingOfficerCheckbox.click();
    await supervisorCheckbox.click();
    await taskInvoicingOfficerCheckbox.click();
    await supervisorCheckbox.click();
    await expect(taskInvoicingOfficerCheckbox).toBeChecked();
    await expect(supervisorCheckbox).toBeChecked();
    await taskInvoicingOfficerCheckbox.click();
    await supervisorCheckbox.click();
    await servicesCounselorCheckbox.click();
    await supervisorCheckbox.click();
    await expect(servicesCounselorCheckbox).toBeChecked();
    await expect(supervisorCheckbox).toBeChecked();
    // Check selecting roles after having supervisor selected for unallowed roles
    await customerCheckbox.click();
    await expect(customerCheckbox).not.toBeChecked();
    await contractingOfficerCheckbox.click();
    await expect(contractingOfficerCheckbox).not.toBeChecked();
    await primeSimulatorCheckbox.click();
    await expect(primeSimulatorCheckbox).not.toBeChecked();
    await qualityAssuranceEvaluatorCheckbox.click();
    await expect(qualityAssuranceEvaluatorCheckbox).not.toBeChecked();
    await customerServiceRepersentativeCheckbox.click();
    await expect(customerServiceRepersentativeCheckbox).not.toBeChecked();
    await governmentSurveillanceRepresentativeCheckbox.click();
    await expect(governmentSurveillanceRepresentativeCheckbox).not.toBeChecked();
    await headquartersCheckbox.click();
    await expect(headquartersCheckbox).not.toBeChecked();

    // Check selecting roles after having supervisor selected for allowed roles
    await taskOrderingOfficerCheckbox.click();
    await expect(taskOrderingOfficerCheckbox).toBeChecked();
    await taskInvoicingOfficerCheckbox.click();
    await expect(taskInvoicingOfficerCheckbox).toBeChecked();
    await servicesCounselorCheckbox.click();
    await expect(servicesCounselorCheckbox).toBeChecked();

    // Continue test to ensure form still submits with valid information
    await servicesCounselorCheckbox.click();
    await expect(servicesCounselorCheckbox).toBeChecked();

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
