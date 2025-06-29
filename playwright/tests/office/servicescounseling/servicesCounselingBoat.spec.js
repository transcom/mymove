// @ts-check
import { test, expect } from './servicesCounselingTestFixture';

test.describe('Services counselor user', () => {
  test('Services Counselor can create a boat shipment and view shipment card info', async ({ page, scPage }) => {
    const move = await scPage.testHarness.buildHHGMoveWithNTSAndNeedsSC();
    await scPage.navigateToMove(move.locator);

    const deliveryDate = new Date(Date.now() + 24 * 60 * 60 * 1000).toLocaleDateString('en-US');
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
    await page.getByText('Use pickup address').click();
    await page.locator('#requestedDeliveryDate').fill('16 Mar 2022');
    await page.locator('#requestedDeliveryDate').blur();

    await page.getByLabel('Counselor remarks').fill('Sample counselor remarks');

    // Save the shipment
    await page.getByRole('button', { name: 'Save' }).click();
    await scPage.waitForPage.moveDetails();

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

  test('Services Counselor can edit a boat shipment', async ({ page, scPage }) => {
    const move = await scPage.testHarness.buildBoatHaulAwayMoveNeedsSC();
    await scPage.navigateToMove(move.locator);

    const deliveryDate = new Date(Date.now() + 24 * 60 * 60 * 1000).toLocaleDateString('en-US');

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

    await page.getByRole('button', { name: 'Edit shipment' }).click();

    await expect(page.getByRole('heading', { level: 1 })).toHaveText('Edit shipment details');
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
    await page.getByText('Use pickup address').click();
    await page.locator('#requestedDeliveryDate').fill('16 Mar 2022');
    await page.locator('#requestedDeliveryDate').blur();

    await page.getByLabel('Counselor remarks').fill('Sample counselor remarks');

    // Save the shipment
    await page.getByRole('button', { name: 'Save' }).click();
    await scPage.waitForPage.moveDetails();

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

  test('Services Counselor can delete a boat shipment', async ({ page, scPage }) => {
    const move = await scPage.testHarness.buildBoatHaulAwayMoveNeedsSC();
    await scPage.navigateToMove(move.locator);

    await page.getByRole('button', { name: 'Edit shipment' }).click();

    await expect(page.getByRole('heading', { level: 1 })).toHaveText('Edit shipment details');
    await expect(page.getByTestId('tag')).toHaveText('Boat');
    await page.getByRole('button', { name: 'Delete shipment' }).click();

    await expect(page.getByText('Are you sure?')).toBeVisible();
    await page.getByTestId('modal').getByTestId('button').click();
    await scPage.waitForPage.moveDetails();

    await expect(page.getByTestId('ShipmentContainer')).toHaveCount(0);
  });

  test('Services Counselor can submit a boat shipment', async ({ page, scPage }) => {
    const move = await scPage.testHarness.buildBoatHaulAwayMoveNeedsSC();
    await scPage.navigateToMove(move.locator);

    await page.getByRole('button', { name: 'Submit move details' }).click();
    await expect(page.getByText('Are you sure?')).toBeVisible();
    await page.getByRole('button', { name: 'Yes, submit' }).click();
    await expect(page.getByText('Move submitted.')).toBeVisible();
  });
});
