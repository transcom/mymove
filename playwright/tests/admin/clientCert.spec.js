// @ts-check
import { test, expect } from '../utils/admin/adminTest';

test('Client Cert List Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  await page.getByRole('menuitem', { name: 'Client Certs' }).click();
  await adminPage.waitForPage.adminPage();
  await expect(page.getByRole('heading', { name: 'Client Certs' })).toBeVisible();

  const columnLabels = [
    'Subject',
    'Id',
    'Sha256 digest',
    'User Id',
    'Prime API',
    'PPTAS API',
    'PPTAS Affiliation',
    'Orders API',
    'USAF Orders Read',
    'USAF Orders Write',
    'Army Orders Read',
    'Army Orders Write',
    'USCG Orders Read',
    'USCG Orders Write',
    'USMC Orders Read',
    'USMC Orders Write',
    'Navy Orders Read',
    'Navy Orders Write',
  ];

  await adminPage.expectRoleLabelsByText('columnheader', columnLabels);
});

test('Client Cert Create Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  await page.getByRole('menuitem', { name: 'Client Certs' }).click();
  await adminPage.waitForPage.adminPage();
  await expect(page.getByRole('heading', { name: 'Client Certs' })).toBeVisible();

  await page.getByRole('link', { name: 'Create' }).click();
  await adminPage.waitForPage.adminPage();

  await expect(page.getByRole('heading', { name: 'Create Client Certs' })).toBeVisible();

  // we need to add the date to the email so that it is unique every time (only one record per email allowed in db)
  const testEmail = `cy.admin_user.${Date.now()}@example.com`;
  await page.getByLabel('Email').fill(testEmail);

  const firstName = page.getByLabel('Subject');
  await firstName.focus();
  await firstName.fill('Test');

  const lastName = page.getByLabel('Sha256 digest');
  await lastName.focus();
  await lastName.fill('test digest');

  await page.getByLabel('Allow prime').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow p p t a s').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('PPTAS Affiliation').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow air force orders read').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow air force orders write').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow army orders read').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow army orders write').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow coast guard orders read').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow coast guard orders write').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow marine corps orders read').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow marine corps orders write').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow navy orders read').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow navy orders write').click();
  await page.getByRole('option').first().click();
});

test('Client Cert Show Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  await page.getByRole('menuitem', { name: 'Client Certs' }).click();
  await adminPage.waitForPage.adminPage();
  await expect(page.getByRole('heading', { name: 'Client Certs' })).toBeVisible();

  // Click first client cert row
  await page.locator('tbody >> tr').first().click();
  await adminPage.waitForPage.adminPage();

  // Get first id field and check that it's in the URL
  const id = await page.locator('.ra-field-id > span').first().textContent();
  expect(page.url()).toContain(id);

  const labels = [
    'Id',
    'Subject',
    'Sha256 digest',
    'User Id',
    'Allow Orders API',
    'Allow air force orders read',
    'Allow air force orders write',
    'Allow army orders read',
    'Allow army orders write',
    'Allow coast guard orders read',
    'Allow coast guard orders write',
    'Allow marine corps orders read',
    'Allow marine corps orders write',
    'Allow navy orders read',
    'Allow navy orders write',
    'Allow prime',
    'Allow PPTAS',
    'PPTAS Affiliation',
    'Created at',
    'Updated at',
  ];

  await adminPage.expectLabels(labels);
});

test('Admin Users Edit Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  await page.getByRole('menuitem', { name: 'Client Certs' }).click();
  await adminPage.waitForPage.adminPage();
  await expect(page.getByRole('heading', { name: 'Client Certs' })).toBeVisible();

  // Click first client cert row
  await page.locator('tbody >> tr').first().click();
  await adminPage.waitForPage.adminPage();

  // Get first id field and check that it's in the URL
  const id = await page.locator('.ra-field-id > span').first().textContent();
  expect(page.url()).toContain(id);

  await page.getByRole('link', { name: 'Edit' }).click();
  await adminPage.waitForPage.adminPage();

  // Get first id field and check that it's in the URL
  expect(page.url()).toContain(id);

  const disabledFields = ['id', 'userId', 'createdAt', 'updatedAt'];
  for (const field of disabledFields) {
    await expect(page.locator(`#${field}`)).toBeDisabled();
  }

  const firstName = page.getByLabel('Subject');
  await firstName.focus();
  await firstName.clear();
  await firstName.fill('NewSubject');

  const lastName = page.getByLabel('Sha256 digest');
  await lastName.focus();
  await lastName.clear();
  await lastName.fill('TestDigest');

  await page.getByLabel('Allow orders a p i').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow air force orders read').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow air force orders write').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow army orders read').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow army orders write').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow coast guard orders read').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow coast guard orders write').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow marine corps orders read').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow marine corps orders write').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow navy orders read').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow navy orders write').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow prime').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('Allow PPTAS').click();
  await page.getByRole('option').first().click();

  await page.getByLabel('PPTAS Affiliation').click();
  await page.getByRole('option').first().click();
});
