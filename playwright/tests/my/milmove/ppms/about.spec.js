/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, forEachViewport, expect } from './customerPpmTestFixture';

test.describe('About Your PPM', () => {
  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPM();
      await customerPpmPage.signInForPPMWithMove(move);
    });

    [true, false].forEach((selectAdvance) => {
      const advanceText = selectAdvance ? 'with' : 'without';
      test(`can submit actual PPM shipment info ${advanceText} an advance`, async ({ customerPpmPage }) => {
        await customerPpmPage.clickOnGoToMoveButton();
        await customerPpmPage.navigateToAboutPage({ selectAdvance });
      });

      test(`W-2 Address ${advanceText} an advance`, async ({ customerPpmPage }) => {
        await customerPpmPage.clickOnGoToMoveButton();
        await customerPpmPage.navigateToAboutPage({ selectAdvance });
        await customerPpmPage.submitWeightTicketPage();

        await expect(customerPpmPage.page.getByRole('heading', { name: 'Review' })).toBeVisible();
        await expect(customerPpmPage.page.getByRole('heading', { name: 'About Your PPM' })).toBeVisible();
        await customerPpmPage.page.getByTestId('aboutYourPPMEditLink').click();

        await expect(customerPpmPage.page).toHaveURL(/\/moves\/[^/]+\/shipments\/[^/]+\/about/);
        await expect(customerPpmPage.page.getByRole('heading', { name: 'About your PPM' })).toBeVisible();

        // Verify the presence of the "W-2 address" section
        await expect(customerPpmPage.page.locator('h2', { hasText: 'W-2 address' })).toBeVisible();
        await expect(customerPpmPage.page.getByText('What is the address on your W-2?')).toBeVisible();

        // Verify the W-2 address fields are present
        await expect(customerPpmPage.page.locator('input[name="w2Address.streetAddress1"]')).toBeVisible();
        await expect(customerPpmPage.page.locator('input[id="w2Address-input"]')).toBeVisible();

        // Verify the W-2 address1 field value matches the expected value
        await expect(customerPpmPage.page.locator('input[name="w2Address.streetAddress1"]')).toHaveValue(
          '1819 S Cedar Street',
        );
        // Verify the W-2 country is present and contains the Country
        await expect(customerPpmPage.page.getByText('UNITED STATES (US)').nth(2)).toBeVisible();
        // Verify the W-2 location is present and contains City, State, ZIP, and County
        await expect(customerPpmPage.page.getByText('YUMA, AZ 85367 (YUMA)')).toBeVisible();
      });
    });
  });
});
