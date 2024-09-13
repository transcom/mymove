/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../../utils/office/officeTest';

import { TooFlowPage } from './tooTestFixture';

test.describe('TOO user', () => {
  /** @type {TooFlowPage} */
  let tooFlowPage;

  test.beforeEach(async ({ officePage }) => {
    const move = await officePage.testHarness.buildHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO();
    await officePage.signInAsNewTOOUser();
    tooFlowPage = new TooFlowPage(officePage, move);
    await tooFlowPage.waitForLoading();
    await officePage.tooNavigateToMove(tooFlowPage.moveLocator);
  });

  test('TOO can create a mobile home shipment and view shipment card info', async ({ page }) => {
    const deliveryDate = new Date().toLocaleDateString('en-US');
    await page.getByTestId('dropdown').selectOption({ label: 'Mobile Home' });

    await expect(page.getByRole('heading', { level: 1 })).toHaveText('Add shipment details');
    await expect(page.getByTestId('tag')).toHaveText('Mobile Home');

    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');

    await page.locator('#requestedPickupDate').fill(deliveryDate);
    await page.locator('#requestedPickupDate').blur();
    await page.getByText('Use current address').click();
    await page.locator('#requestedDeliveryDate').fill('16 Mar 2022');
    await page.locator('#requestedDeliveryDate').blur();

    // Save the shipment
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByTestId('ShipmentContainer')).toHaveCount(2);

    await expect(page.getByText('Mobile home year')).toBeVisible();
    await expect(page.getByTestId('year')).toHaveText('2022');
    await expect(page.getByText('Mobile home make')).toBeVisible();
    await expect(page.getByTestId('make')).toHaveText('make');
    await expect(page.getByText('Mobile home model')).toBeVisible();
    await expect(page.getByTestId('model')).toHaveText('model');
    await expect(page.getByText('Dimensions')).toBeVisible();
    await expect(page.getByTestId('dimensions')).toHaveText("22' L x 22' W x 22' H");
  });
});
