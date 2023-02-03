/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
const { expect, test, forEachViewport, CustomerPpmPage } = require('./customerPpmTestFixture');

/**
 * CustomerPpmOnboardingPage test fixture. Our linting rules (like
 * no-use-before-define) pushes us towards grouping all these helpers
 * into a class. It also follows the examples at
 * https://playwright.dev/docs/test-fixtures
 *
 * @extends CustomerPpmPage
 */
class CustomerPpmOnboardingPage extends CustomerPpmPage {
  /**
   * verify page and submit to go to next page
   */
  async verifyEstimatedWeightsAndProGear() {
    await this.page.getByRole('button', { name: 'Back' }).click();

    await expect(this.page.locator('input[name="estimatedWeight"]')).toHaveValue('4,000');
    await expect(this.page.locator('label[for="hasProGearYes"]')).toBeChecked();
    await expect(this.page.locator('input[name="proGearWeight"]')).toBeVisible();
    await expect(this.page.locator('input[name="proGearWeight"]')).toHaveValue('500');

    await expect(this.page.locator('input[name="spouseProGearWeight"]')).toBeVisible();
    await expect(this.page.locator('input[name="spouseProGearWeight"]')).toHaveValue('400');

    await this.navigateFromEstimatedWeightsPageToEstimatedIncentivePage();
  }

  /**
   */
  async verifyShipmentSpecificInfoOnEstimatedIncentivePage() {
    const shipmentInfo = this.page.locator('.container li');
    await expect(shipmentInfo.getByText('4,000 lbs')).toBeVisible();
    await expect(shipmentInfo.getByText('90210')).toBeVisible();
    await expect(shipmentInfo.getByText('76127')).toBeVisible();
    await expect(shipmentInfo.getByText('01 Feb 2022')).toBeVisible();
  }

  /**
   */
  async verifyStep5ExistsAndBtnIsDisabled() {
    const stepContainer5 = this.page.locator('[data-testid="stepContainer5"]');
    await expect(stepContainer5.getByRole('button', { name: 'Upload PPM Documents' })).toBeDisabled();
    await expect(
      stepContainer5.locator('p').getByText('After a counselor approves your PPM, you will be able to:'),
    ).toBeVisible();
  }

  /**
   * update the form values by submitting and then return to the
   * page to verify if the values persist and then return to the
   * next page
   *
   */
  async submitAndVerifyUpdateDateAndLocation() {
    await this.page.locator('input[name="pickupPostalCode"]').clear();
    await this.page.locator('input[name="pickupPostalCode"]').type('90210');
    await this.page.locator('input[name="pickupPostalCode"]').blur();

    await this.page.locator('input[name="secondaryPickupPostalCode"]').clear();
    await this.page.locator('input[name="secondaryPickupPostalCode"]').type('90212');
    await this.page.locator('input[name="secondaryPickupPostalCode"]').blur();

    await this.page.locator('input[name="destinationPostalCode"]').clear();
    await this.page.locator('input[name="destinationPostalCode"]').type('76127');
    // TODO: The user has secondary destination zips. We should test clearing this value by selecting the no radio btn. This doesn't work atm
    await this.page.locator('label[for="sitExpectedNo"]').click();

    await this.page.locator('input[name="expectedDepartureDate"]').clear();
    await this.page.locator('input[name="expectedDepartureDate"]').type('01 Feb 2022');
    await this.page.locator('input[name="expectedDepartureDate"]').blur();

    // Change closeout location
    await this.selectDutyLocation('Fort Bragg', 'closeoutOffice');

    await this.navigateFromDateAndLocationPageToEstimatedWeightsPage();

    await this.page.getByRole('button', { name: 'Back' }).click();

    // verify values
    await expect(this.page.locator('input[name="pickupPostalCode"]')).toHaveValue('90210');
    await expect(this.page.locator('label[for="yes-secondary-pickup-postal-code"]')).toBeChecked();
    await expect(this.page.locator('input[name="secondaryPickupPostalCode"]')).toHaveValue('90212');
    await expect(this.page.locator('input[name="destinationPostalCode"]')).toHaveValue('76127');
    await expect(this.page.locator('label[for="hasSecondaryDestinationPostalCodeYes"]')).toBeChecked();
    await expect(this.page.locator('input[name="expectedDepartureDate"]')).toHaveValue('01 Feb 2022');
    await expect(this.page.locator('label[for="sitExpectedNo"]')).toBeChecked();
    await expect(this.page.locator('label[for="sitExpectedNo"]')).toHaveValue('false');

    await this.navigateFromDateAndLocationPageToEstimatedWeightsPage();
  }
}

test.describe('Entire PPM onboarding flow', () => {
  /** @type {CustomerPpmOnboardingPage} */
  let customerPpmOnboardingPage;

  forEachViewport(async ({ isMobile }) => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildDraftMoveWithPPMWithDepartureDate();
      customerPpmOnboardingPage = new CustomerPpmOnboardingPage(customerPpmPage);
      customerPpmOnboardingPage.signInForPPMWithMove(move);
    });

    test('flows through happy path for existing shipment', async () => {
      await customerPpmOnboardingPage.navigateFromHomePageToExistingPPMDateAndLocationPage();
      await customerPpmOnboardingPage.submitsDateAndLocation();
      await customerPpmOnboardingPage.submitsEstimatedWeightsAndProGear();
      await customerPpmOnboardingPage.generalVerifyEstimatedIncentivePage({ isMobile });
      await customerPpmOnboardingPage.submitsAdvancePage({ addAdvance: true, isMobile });
      await customerPpmOnboardingPage.navigateToAgreementAndSign();
      await customerPpmOnboardingPage.submitMove();
      await customerPpmOnboardingPage.verifyStep5ExistsAndBtnIsDisabled();
    });

    test('happy path with edits and backs', async () => {
      await customerPpmOnboardingPage.navigateFromHomePageToExistingPPMDateAndLocationPage();

      await customerPpmOnboardingPage.submitAndVerifyUpdateDateAndLocation();

      await customerPpmOnboardingPage.submitsEstimatedWeightsAndProGear();
      await customerPpmOnboardingPage.verifyEstimatedWeightsAndProGear();

      await customerPpmOnboardingPage.verifyShipmentSpecificInfoOnEstimatedIncentivePage();
      await customerPpmOnboardingPage.generalVerifyEstimatedIncentivePage({ isMobile });

      await customerPpmOnboardingPage.submitsAdvancePage({ addAdvance: true, isMobile });

      await customerPpmOnboardingPage.navigateToAgreementAndSign();

      await customerPpmOnboardingPage.submitMove();
      await customerPpmOnboardingPage.verifyStep5ExistsAndBtnIsDisabled();
    });
  });
});
