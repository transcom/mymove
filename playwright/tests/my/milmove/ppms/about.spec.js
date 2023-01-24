/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { test, CustomerPpmPage } = require('./customerPpmTestFixture');

test.describe('About Your PPM', () => {
  /** @type {CustomerPpmPage} */
  let customerPpmPage;

  test.beforeEach(async ({ customerPage }) => {
    const move = await customerPage.testHarness.buildApprovedMoveWithPPM();
    customerPpmPage = new CustomerPpmPage(customerPage, move);
  });

  // https://playwright.dev/docs/test-parameterize
  // use forEach to avoid https://eslint.org/docs/latest/rules/no-loop-func
  ['mobile', 'desktop'].forEach((viewport) => {
    test.describe(`with ${viewport} viewport`, async () => {
      if (viewport === 'mobile') {
        // https://playwright.dev/docs/emulation#viewport
        test.use({ viewport: { width: 479, height: 875 } });
      }
      [true, false].forEach((selectAdvance) => {
        const withAdvance = selectAdvance ? 'with' : 'without';
        test(`can submit actual PPM shipment info ${withAdvance} an advance`, async () => {
          await customerPpmPage.signInAndNavigateToAboutPage(selectAdvance);
        });
      });
    });
  });
});
