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

  test('TOO can create a boat shipment and view shipment card info', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO();
    await officePage.signInAsNewTOOUser();
    tooFlowPage = new TooFlowPage(officePage, move);
    await tooFlowPage.waitForLoading();
    await officePage.tooNavigateToMove(tooFlowPage.moveLocator);

    const deliveryDate = new Date().toLocaleDateString('en-US');
    await page.getByTestId('dropdown').selectOption({ label: 'Boat' });

    await expect(page.getByRole('heading', { level: 1 })).toHaveText('Add shipment details');
    await expect(page.getByTestId('tag')).toHaveText('Boat');

    await page.getByLabel('Year').fill('2022');
    await page.getByLabel('Make').fill('make');
    await page.getByLabel('Model').fill('model');
    await page.getByTestId('lengthFeet').fill('22');
    await page.getByTestId('widthFeet').fill('22');
    await page.getByTestId('heightFeet').fill('22');
    await page.locator('label[for="hasTrailerYes"]').click();
    await page.locator('label[for="isRoadworthyYes"]').click();
    await page.locator('label[for="typeTowAway"]').click();

    await page.locator('#requestedPickupDate').fill(deliveryDate);
    await page.locator('#requestedPickupDate').blur();
    await page.getByText('Use current address').click();
    await page.locator('#requestedDeliveryDate').fill('16 Mar 2022');
    await page.locator('#requestedDeliveryDate').blur();

    // Save the shipment
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByTestId('ShipmentContainer')).toHaveCount(2);

    await expect(page.getByText('Shipment method')).toBeVisible();
    await expect(page.getByTestId('shipmentType')).toHaveText('BTA');
    await expect(page.getByText('Boat year')).toBeVisible();
    await expect(page.getByTestId('year')).toHaveText('2022');
    await expect(page.getByText('Boat make')).toBeVisible();
    await expect(page.getByTestId('make')).toHaveText('make');
    await expect(page.getByText('Boat model')).toBeVisible();
    await expect(page.getByTestId('model')).toHaveText('model');
    await expect(page.getByText('Dimensions')).toBeVisible();
    await expect(page.getByTestId('dimensions')).toHaveText("22' L x 22' W x 22' H");
    await expect(page.getByText('Trailer', { exact: true })).toBeVisible();
    await expect(page.getByTestId('trailer')).toHaveText('Yes');
    await expect(page.getByText('Is trailer roadworthy')).toBeVisible();
    await expect(page.getByTestId('isRoadworthy')).toHaveText('Yes');
  });

  test('TOO can edit a boat shipment', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildBoatHaulAwayMoveNeedsTOOApproval();
    await officePage.signInAsNewTOOUser();
    tooFlowPage = new TooFlowPage(officePage, move);
    await tooFlowPage.waitForLoading();
    await officePage.tooNavigateToMove(tooFlowPage.moveLocator);

    const deliveryDate = new Date().toLocaleDateString('en-US');

    await expect(page.getByText('Shipment method')).toBeVisible();
    await expect(page.getByTestId('shipmentType')).not.toHaveText('BTA');
    await expect(page.getByText('Boat year')).toBeVisible();
    await expect(page.getByTestId('year')).not.toHaveText('2022');
    await expect(page.getByText('Boat make', { exact: true })).toBeVisible();
    await expect(page.getByTestId('make')).not.toHaveText('make');
    await expect(page.getByText('Boat model', { exact: true })).toBeVisible();
    await expect(page.getByTestId('model')).not.toHaveText('model');
    await expect(page.getByText('Dimensions')).toBeVisible();
    await expect(page.getByTestId('dimensions')).not.toHaveText("22' L x 22' W x 22' H");
    await expect(page.getByText('Trailer', { exact: true })).toBeVisible();
    await expect(page.getByTestId('trailer')).toHaveText('Yes');
    await expect(page.getByText('Is trailer roadworthy')).toBeVisible();
    await expect(page.getByTestId('isRoadworthy')).not.toHaveText('Yes');

    // Edit the shipment
    await page.getByRole('button', { name: 'Edit shipment' }).click();
    await expect(page.getByRole('heading', { level: 1 })).toHaveText('Edit shipment details');
    await expect(page.getByTestId('tag')).toHaveText('Boat');

    await page.getByLabel('Year').fill('2020');
    await page.getByLabel('Make').fill('changedmake');
    await page.getByLabel('Model').fill('changedmodel');
    await page.getByTestId('lengthFeet').fill('123');
    await page.getByTestId('widthFeet').fill('222');
    await page.getByTestId('heightFeet').fill('333');
    await page.locator('label[for="hasTrailerNo"]').click();
    await page.locator('label[for="typeHaulAway"]').click();

    await page.locator('#requestedPickupDate').fill(deliveryDate);
    await page.locator('#requestedPickupDate').blur();
    await page.getByText('Use current address').click();
    await page.locator('#requestedDeliveryDate').fill('19 Mar 2022');
    await page.locator('#requestedDeliveryDate').blur();
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByText('Shipment method')).toBeVisible();
    await expect(page.getByTestId('shipmentType')).toHaveText('BHA');
    await expect(page.getByText('Boat year')).toBeVisible();
    await expect(page.getByTestId('year')).toHaveText('2020');
    await expect(page.getByText('Boat make')).toBeVisible();
    await expect(page.getByTestId('make')).toHaveText('changedmake');
    await expect(page.getByText('Boat model')).toBeVisible();
    await expect(page.getByTestId('model')).toHaveText('changedmodel');
    await expect(page.getByText('Dimensions')).toBeVisible();
    await expect(page.getByTestId('dimensions')).toHaveText("123' L x 222' W x 333' H");
    await expect(page.getByText('Trailer', { exact: true })).toBeVisible();
    await expect(page.getByTestId('trailer')).toHaveText('No');
    await expect(page.getByText('Is trailer roadworthy')).not.toBeVisible();
  });

  test('TOO can delete a boat shipment', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildBoatHaulAwayMoveNeedsTOOApproval();
    await officePage.signInAsNewTOOUser();
    tooFlowPage = new TooFlowPage(officePage, move);
    await tooFlowPage.waitForLoading();
    await officePage.tooNavigateToMove(tooFlowPage.moveLocator);

    await page.getByRole('button', { name: 'Edit shipment' }).click();

    await expect(page.getByRole('heading', { level: 1 })).toHaveText('Edit shipment details');
    await expect(page.getByTestId('tag')).toHaveText('Boat');
    await page.getByRole('button', { name: 'Delete shipment' }).click();

    await expect(page.getByRole('heading', { level: 3 })).toHaveText('Are you sure?');
    await page.getByTestId('modal').getByTestId('button').click();
    await officePage.waitForPage.moveDetails();

    await expect(page.getByTestId('ShipmentContainer')).toHaveCount(0);
  });

  test('TOO can approve a boat shipment', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildBoatHaulAwayMoveNeedsTOOApproval();
    await officePage.signInAsNewTOOUser();
    tooFlowPage = new TooFlowPage(officePage, move);
    await tooFlowPage.waitForLoading();
    await officePage.tooNavigateToMove(tooFlowPage.moveLocator);

    await expect(page.getByTestId('ShipmentContainer')).toBeVisible();
    await page.getByTestId('ShipmentContainer').getByTestId('checkbox').locator('label').click();
    await page.getByTestId('shipmentApproveButton').click();
    await page.getByRole('button', { name: 'Approve and send' }).click();
    await page.getByTestId('MoveDetails-Tab').click();
    await expect(page.getByRole('heading', { name: 'Approved Shipments' })).toBeVisible();
  });
});
