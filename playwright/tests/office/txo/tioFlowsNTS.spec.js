// @ts-check
import { test, expect } from '../../utils/office/officeTest';

import { TioFlowPage } from './tioTestFixture';

test.describe('TIO user', () => {
  let tioFlowPage;
  test('A TIO can review and understand calculations from INPK pricing on an iNTS shipment', async ({
    page,
    officePage,
  }) => {
    const move =
      await officePage.testHarness.buildInternationalHHGIntoInternationalNTSMoveWithServiceItemsandPaymentRequestsForTIO();
    await officePage.signInAsNewTIOUser();
    tioFlowPage = new TioFlowPage(officePage, move);
    await tioFlowPage.waitForLoading();
    await officePage.tioNavigateToMove(tioFlowPage.moveLocator);
    await page.waitForSelector('#payment-requests');
    // Review the second payment request on this move (Testharness makes the 2nd one INPK)
    await page
      .locator('#payment-requests div')
      .filter({ has: page.locator('h6', { hasText: /-2$/ }) })
      .locator('button[data-testid="reviewBtn"]')
      .click();

    // Show the calculation breakdown
    await page.locator('button[data-testid="toggleCalculations"]').click();
    await expect(page.locator('[data-testid="ServiceItemCalculations"]')).toBeVisible();

    //
    // Verify the breakdown
    //

    // CWT
    await expect(
      page.locator('[data-testid="column"]', { hasText: 'Billable weight (cwt)' }).locator('[data-testid="value"]'),
    ).toHaveText('9.8 cwt');

    // IHPK base price
    await expect(
      page.locator('[data-testid="column"]', { hasText: 'International Pack price' }).locator('[data-testid="value"]'),
    ).toHaveText('69.97');

    // Reference date (Requested pickup)
    await expect(
      page
        .locator('[data-testid="column"]', { hasText: 'International Pack price' })
        .locator('[data-testid="details"] >> text=Requested pickup'),
    ).toContainText('15 Mar 2024');

    // Contract escalation factor
    await expect(
      page.locator('[data-testid="column"]', { hasText: 'Price escalation factor' }).locator('[data-testid="value"]'),
    ).toHaveText('1.11000');

    // NTS packing factor
    await expect(
      page.locator('[data-testid="column"]', { hasText: 'NTS packing factor' }).locator('[data-testid="value"]'),
    ).toHaveText('1.45');

    // Total = (Base price * escalation) * cwt * nts factor
    await expect(
      page.locator('[data-testid="column"]', { hasText: 'Total:' }).locator('[data-testid="value"]'),
    ).toHaveText('$1,103.70');
  });
});
