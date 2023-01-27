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
    const move = await customerPage.testHarness.buildApprovedMoveWithPPM();
    customerPpmPage = new CustomerPpmPage(customerPage, move);
  });

  //
  // https://playwright.dev/docs/test-parameterize
  //
  // use forEach to avoid
  // https://eslint.org/docs/latest/rules/no-loop-func
  [true, false].forEach((isMobile) => {
    const viewportName = isMobile ? 'mobile' : 'desktop';
    [true, false].forEach((selectAdvance) => {
      const advanceText = selectAdvance ? 'with' : 'without';
      test.describe(`with ${viewportName} viewport`, async () => {
        if (isMobile) {
          useMobileViewport();
        }
        test(`can submit actual PPM shipment info ${advanceText} an advance`, async () => {
          await customerPpmPage.signInAndNavigateToAboutPage({ selectAdvance });
        });
      });
    });
  });
});
