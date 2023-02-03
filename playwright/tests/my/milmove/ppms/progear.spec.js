/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, forEachViewport } = require('./customerPpmTestFixture');

test.describe('Progear', () => {
  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPMProgearWeightTicket();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.navigateToProgearPage();
    });

    test(`progear page loads`, async ({ customerPpmPage, page }) => {
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
