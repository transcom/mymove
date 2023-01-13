/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, expect } = require('../utils/adminTest');

test.describe('Webhook Subscriptions', () => {
  test('successfully navigates to the webhook subscriptions list page', async ({ page, adminPage }) => {
    await adminPage.testHarness.buildWebhookSubscription();
    await adminPage.signInAsNewAdminUser();

    await page.getByRole('menuitem', { name: 'Webhook Subscriptions' }).click();
    await adminPage.waitForAdminPageToLoad();
    await expect(page.getByRole('heading', { name: 'Webhook Subscriptions' })).toBeVisible();

    const columnLabels = ['Id', 'Event key', 'Severity', 'Callback url', 'Status', 'Updated at'];
    await adminPage.expectRoleLabelsByText('columnheader', columnLabels);
  });
});

test.describe('WebhookSubscriptions Details Show Page', () => {
  test('pulls up details page for a webhook subscription', async ({ page, adminPage }) => {
    const webhook = await adminPage.testHarness.buildWebhookSubscription();
    const id = webhook.ID;
    await adminPage.signInAsNewAdminUser();

    await page.getByRole('menuitem', { name: 'Webhook Subscriptions' }).click();
    await adminPage.waitForAdminPageToLoad();
    await expect(page.getByRole('heading', { name: 'Webhook Subscriptions' })).toBeVisible();
    await page.getByText(id).click();
    await adminPage.waitForAdminPageToLoad();

    // check that the webhookSubscription's ID is shown in the page title
    expect(page.url()).toContain(id);
    await expect(page.getByRole('heading', { name: `Webhook Subscription ID: ${id}` })).toBeVisible();

    const labels = [
      'Id',
      'Subscriber Id',
      'Status',
      'Event key',
      'Callback url',
      'Created at',
      'Updated at',
      'Severity',
    ];
    await adminPage.expectLocatorLabelsByText('label', labels);
  });
});

test.describe('WebhookSubscriptions Details Edit Page', () => {
  test('pulls up edit page for a webhook subscription', async ({ page, adminPage }) => {
    const webhook = await adminPage.testHarness.buildWebhookSubscription();
    const id = webhook.ID;
    await adminPage.signInAsNewAdminUser();

    await page.getByRole('menuitem', { name: 'Webhook Subscriptions' }).click();
    await adminPage.waitForAdminPageToLoad();
    await expect(page.getByRole('heading', { name: 'Webhook Subscriptions' })).toBeVisible();
    await page.getByText(id).click();
    await adminPage.waitForAdminPageToLoad();

    // check labels on edit page
    const labels = [
      'Id',
      'Subscriber Id',
      'Status',
      'Event key',
      'Callback url',
      'Created at',
      'Updated at',
      'Severity',
    ];
    await adminPage.expectLocatorLabelsByText('label', labels);

    await page.getByRole('button', { name: 'Edit' }).click();
    await adminPage.waitForAdminPageToLoad();

    // Change webhook subscription status
    // await page.locator('label', { hasText: 'Status' }).click();

    await page.getByRole('button', { name: 'Status Active' }).click();
    await page.getByRole('option', { name: 'Disabled' }).click();

    await page.getByRole('button', { name: 'Save' }).click();
    await adminPage.waitForAdminPageToLoad();

    // Check that the webhook subscription status was changed
    await expect(page.locator(`tr:has(:text("${id}")) >> td.column-status`)).toHaveText('DISABLED');
  });
});

test.describe('Webhook Subscription Create Page', () => {
  test('pulls up create page for a webhook subscription', async ({ page, adminPage }) => {
    const webhook = await adminPage.testHarness.buildWebhookSubscription();
    const subId = webhook.SubscriberID;
    await adminPage.signInAsNewAdminUser();

    await page.getByRole('menuitem', { name: 'Webhook Subscriptions' }).click();
    await adminPage.waitForAdminPageToLoad();
    await expect(page.getByRole('heading', { name: 'Webhook Subscriptions' })).toBeVisible();
    await page.getByRole('button', { name: 'Create' }).click();

    await page.getByLabel('Subscriber Id').fill(subId);
    await page.getByLabel('Event key').fill('PaymentRequest.Update');
    await page.getByLabel('Callback url').fill('https://test1.example.com');
    await page.locator('div[id="status"]').click();
    await page.locator('li[data-value="ACTIVE"]').click();
    await page.getByRole('button').filter({ hasText: 'Save' }).click();
    await adminPage.waitForAdminPageToLoad();

    const id = await page.locator('#id').inputValue();
    expect(page.url()).toContain(id);
    await expect(page.getByRole('heading', { name: `Webhook Subscriptions #${id}` })).toBeVisible();

    await expect(page.locator('#subscriberId')).toHaveValue(subId);
    await expect(page.locator('#eventKey')).toHaveValue('PaymentRequest.Update');
    await expect(page.locator('#callbackUrl')).toHaveValue('https://test1.example.com');
    await expect(page.locator('#status')).toHaveText('Active');
    await expect(page.locator('#severity')).toHaveText('0');
  });
});
