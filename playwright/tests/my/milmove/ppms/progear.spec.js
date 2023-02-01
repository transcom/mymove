/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, forEachViewport, CustomerPpmPage } = require('./customerPpmTestFixture');

test.describe('Progear', () => {
  /** @type {CustomerPpmPage} */
  let customerPpmPage;

  test.beforeEach(async ({ customerPage }) => {
    const move = await customerPage.testHarness.buildApprovedMoveWithPPMProgearWeightTicket();
    customerPpmPage = new CustomerPpmPage(customerPage, move);
    await customerPpmPage.signInAndNavigateToProgearPage();
  });

  forEachViewport(async () => {
    test(`progear page loads`, async ({ page }) => {
      await customerPpmPage.submitProgearPage({ belongsToSelf: true });

      const set2Heading = page.getByRole('heading', { name: 'Set 2' });
      await expect(set2Heading).toBeVisible();
      const progearSection = set2Heading.locator('../..').last();
      await expect(progearSection).toContainText('Pro-gear');
      await expect(progearSection).toContainText('Radio equipment');
      // contains won't work when text content is divided between multiple semantic html elements
      await expect(progearSection).toContainText('Weight:');
      await expect(progearSection).toContainText('2,000 lbs');
    });
  });
});
