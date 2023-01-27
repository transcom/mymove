/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, useMobileViewport, CustomerPpmPage } = require('./customerPpmTestFixture');

test.describe('About Your PPM', () => {
  /** @type {CustomerPpmPage} */
  let customerPpmPage;

  test.beforeEach(async ({ customerPage }) => {
    const move = await customerPage.testHarness.buildApprovedMoveWithPPMWithAboutFormComplete();
    customerPpmPage = new CustomerPpmPage(customerPage, move);
    await customerPpmPage.signInAndNavigateToWeightTicketPage();
  });

  //
  // https://playwright.dev/docs/test-parameterize
  //
  // use forEach to avoid
  // https://eslint.org/docs/latest/rules/no-loop-func
  [true, false].forEach((isMobile) => {
    const viewportName = isMobile ? 'mobile' : 'desktop';
    test.describe(`with ${viewportName} viewport`, async () => {
      if (isMobile) {
        useMobileViewport();
      }

      test('proceed with weight ticket documents', async () => {
        await customerPpmPage.submitWeightTicketPage();
      });
      test('proceed with claiming trailer', async () => {
        await customerPpmPage.submitWeightTicketPage({ hasTrailer: true, ownTrailer: true });
      });
      test('proceed without claiming trailer', async () => {
        await customerPpmPage.submitWeightTicketPage({ hasTrailer: true, ownTrailer: false });
      });
      test('proceed with constructed weight ticket documents', async () => {
        await customerPpmPage.submitWeightTicketPage({ useConstructedWeight: true });
      });
    });
  });
});
