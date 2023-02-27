/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../utils/adminTest';

test('Moves Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  // make sure at least one move exists
  await adminPage.testHarness.buildDefaultMove();

  await page.getByRole('menuitem', { name: 'Moves' }).click();
  expect(page.url()).toContain('/system/moves');
  await expect(page.getByRole('heading', { name: 'Moves' })).toBeVisible();

  const columnLabels = ['Id', 'Order Id', 'Service Member Id', 'Locator', 'Status', 'Show', 'Created at', 'Updated at'];
  await adminPage.expectRoleLabelsByText('columnheader', columnLabels);
});

test('Moves Details Show Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  // make sure at least one move exists
  await adminPage.testHarness.buildDefaultMove();

  await page.getByRole('menuitem', { name: 'Moves' }).click();
  await adminPage.waitForAdminPageToLoad();
  await page.locator('span[reference="moves"]').first().click();
  await adminPage.waitForAdminPageToLoad();

  const id = await page.locator('div:has(label :text-is("Id")) > div > span').textContent();
  expect(page.url()).toContain(id);
  await expect(page.getByRole('heading', { name: `Move ID: ${id}` })).toBeVisible();

  const labels = [
    'Id',
    'Locator',
    'Status',
    'Show',
    'Order Id',
    'Created at',
    'Updated at',
    'User Id',
    'Service member Id',
    'Service member first name',
    'Service member middle name',
    'Service member last name',
  ];

  await adminPage.expectLabels(labels);
});

test('Moves Details Edit Page', async ({ page, adminPage }) => {
  await adminPage.signInAsNewAdminUser();

  const move = await adminPage.testHarness.buildDefaultMove();
  const moveId = move.id;
  const moveLocator = move.locator;

  await page.getByRole('menuitem', { name: 'Moves' }).click();
  await adminPage.waitForAdminPageToLoad();

  // use locator search to find move in case move is not on first page
  // entering the move locator should auto search without a click
  await page.getByLabel('locator').fill(moveLocator);

  // click on row for newly created move
  // if this test has been run many times locally, this might fail
  // because the new move is not on the first page of results
  await page.locator(`tr:has(:text("${moveId}"))`).click();
  await adminPage.waitForAdminPageToLoad();
  expect(page.url()).toContain(moveId);

  await page.getByRole('link', { name: 'Edit' }).click();
  await adminPage.waitForAdminPageToLoad();
  expect(page.url()).toContain(moveId);

  const disabledFields = [
    'id',
    'locator',
    'status',
    'ordersId',
    'createdAt',
    'updatedAt',
    'serviceMember.userId',
    'serviceMember.id',
    'serviceMember.firstName',
    'serviceMember.middleName',
    'serviceMember.lastName',
  ];
  for (const field of disabledFields) {
    await expect(page.locator(`[id=${field.replace('.', '\\.')}]`)).toBeDisabled();
  }

  // set the move to the show status it did NOT have before
  const showStatus = await page.locator('div:has(label :text-is("Show")) >> input[name="show"]').inputValue();

  const newStatus = (showStatus !== 'true').toString();
  await page.locator('div:has(label :text-is("Show")) >> #show').click();
  await page.locator(`ul[aria-labelledby="show-label"] >> li[data-value="${newStatus}"]`).click();

  await page.getByRole('button', { name: 'Save' }).click();

  // back to list of all moves
  expect(page.url()).not.toContain(moveId);

  await expect(page.locator(`tr:has(:text("${moveId}")) >> td.column-show >> span`)).toHaveText(newStatus);
});
