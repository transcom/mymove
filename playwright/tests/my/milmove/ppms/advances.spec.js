/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, useMobileViewport, CustomerPpmPage } = require('./customerPpmTestFixture');

test.describe('About Your PPM', () => {
  /** @type {CustomerPpmPage} */
  let customerPpmPage;

  test.beforeEach(async ({ customerPage }) => {
    const move = await customerPage.testHarness.buildUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights();
    customerPpmPage = new CustomerPpmPage(customerPage, move);
    await customerPpmPage.signInAndNavigateFromHomePageToExistingPPMDateAndLocationPage();
    await customerPpmPage.navigateFromDateAndLocationPageToEstimatedWeightsPage();
    await customerPpmPage.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
    await customerPpmPage.navigateFromEstimatedIncentivePageToAdvancesPage();
  });

  test('does not allow SM to progress if form is in an invalid state', async ({ page }) => {
    await page.locator('label[for="hasRequestedAdvanceYes"]').click();

    // missing advance
    const advanceInput = page.locator('input[name="advanceAmountRequested"]');
    // need to click the first time before clear in playwright
    await advanceInput.click();
    await advanceInput.clear();
    await advanceInput.blur();
    const errorMessage = page.locator('[class="usa-error-message"]');
    await expect(errorMessage).toContainText('Required');
    await expect(page.locator('[class="usa-error-message"] + div').locator('input')).toHaveAttribute(
      'id',
      'advanceAmountRequested',
    );
    const saveButton = page.getByRole('button', { name: 'Save & Continue' });

    await expect(saveButton).toBeDisabled();
    await advanceInput.type('1');
    await advanceInput.blur();
    await expect(page.locator('[class="usa-error-message"]')).not.toBeVisible();

    // advance violates min
    await advanceInput.clear();
    await advanceInput.type('0');
    await advanceInput.blur();
    await expect(errorMessage).toContainText(
      "The minimum advance request is $1. If you don't want an advance, select No.",
    );
    await expect(page.locator('[class="usa-error-message"] + div').locator('input')).toHaveAttribute(
      'id',
      'advanceAmountRequested',
    );
    await expect(saveButton).toBeDisabled();

    await advanceInput.clear();
    await advanceInput.type('1');
    await advanceInput.blur();
    await expect(errorMessage).not.toBeVisible();

    // advance violates max (over 60% of incentive)
    await advanceInput.clear();
    await advanceInput.type('6001');
    await advanceInput.blur();
    await expect(errorMessage).toContainText('Enter an amount $6,000 or less');
    await expect(page.locator('[class="usa-error-message"] + div').locator('input')).toHaveAttribute(
      'id',
      'advanceAmountRequested',
    );
    await expect(saveButton).toBeDisabled();
    await advanceInput.clear();
    await advanceInput.type('1');
    await advanceInput.blur();
    await expect(errorMessage).not.toBeVisible();
  });

  //
  // https://playwright.dev/docs/test-parameterize
  //
  // use forEach to avoid
  // https://eslint.org/docs/latest/rules/no-loop-func
  [true, false].forEach((isMobile) => {
    const viewportName = isMobile ? 'mobile' : 'desktop';
    [true, false].forEach((addAdvance) => {
      const advanceText = addAdvance ? 'request' : 'opt to not receive';
      test.describe(`with ${viewportName} viewport`, async () => {
        if (isMobile) {
          useMobileViewport();
        }
        test(`can ${advanceText} an advance`, async () => {
          await customerPpmPage.submitsAdvancePage({ addAdvance, isMobile });
        });
      });
    });
  });
});
