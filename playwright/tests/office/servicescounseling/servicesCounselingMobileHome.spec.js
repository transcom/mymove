// @ts-check
import { test, expect } from './servicesCounselingTestFixture';

test.describe('Services counselor user', () => {
  test.beforeEach(async ({ scPage }) => {
    const move = await scPage.testHarness.buildHHGMoveWithNTSAndNeedsSC();
    await scPage.navigateToMove(move.locator);
  });

  test('Services Counselor can create a mobile home shipment and view shipment card info', async ({ page, scPage }) => {
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

    await page.getByLabel('Counselor remarks').fill('Sample counselor remarks');

    // Save the shipment
    await page.getByRole('button', { name: 'Save' }).click();
    await scPage.waitForPage.moveDetails();

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
