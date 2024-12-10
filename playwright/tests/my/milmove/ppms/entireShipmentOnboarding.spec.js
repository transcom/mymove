/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { expect, test, forEachViewport, CustomerPpmPage } from './customerPpmTestFixture';

const multiMoveEnabled = process.env.FEATURE_FLAG_MULTI_MOVE;

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
  async verifyManagePPMStepExistsAndBtnIsDisabled() {
    const stepContainer = this.page.locator('[data-testid="stepContainer6"]');

    if (stepContainer == null) {
      this.page.locator('[data-testid="stepContainer5"]');
    }

    await expect(stepContainer.getByRole('button', { name: 'Upload PPM Documents' })).toBeDisabled();
    await expect(
      stepContainer.locator('p').getByText('After a counselor approves your PPM, you will be able to:'),
    ).toBeVisible();
  }

  /**
   * update the form values by submitting and then return to the
   * page to verify if the values persist and then return to the
   * next page
   *
   */
  async submitAndVerifyUpdateDateAndLocation() {
    const pickupLocation = 'BEVERLY HILLS, CA 90212 (LOS ANGELES)';
    await this.page.locator('label[for="yes-secondary-pickup-address"]').click();

    await this.page.locator('input[name="secondaryPickupAddress.address.streetAddress1"]').fill('1234 Street');
    await this.page.locator('input[id="secondaryPickupAddress.address-location-input"]').fill('90212');
    await expect(this.page.getByText(pickupLocation, { exact: true })).toBeVisible();
    await this.page.keyboard.press('Enter');

    // TODO: The user has secondary destination zips. We should test clearing this value by selecting the no radio btn. This doesn't work atm
    await this.page.locator('label[for="sitExpectedNo"]').click();

    await this.page.locator('input[name="expectedDepartureDate"]').clear();
    await this.page.locator('input[name="expectedDepartureDate"]').fill('01 Feb 2022');
    await this.page.locator('input[name="expectedDepartureDate"]').blur();

    // Change closeout location
    await this.selectDutyLocation('Fort Bragg', 'closeoutOffice');
    await this.page.keyboard.press('Backspace'); // tests if backspace clears the duty location field
    await expect(this.page.getByLabel('Which closeout office should review your PPM?')).toBeEmpty();
    await this.selectDutyLocation('Fort Bragg', 'closeoutOffice');
    await this.navigateFromDateAndLocationPageToEstimatedWeightsPage();

    await this.page.getByRole('button', { name: 'Back' }).click();

    await this.navigateFromDateAndLocationPageToEstimatedWeightsPage();
  }
}

test.describe('About Form Date flow', () => {
  /** @type {CustomerPpmOnboardingPage} */
  let customerPpmOnboardingPage;

  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPM();
      customerPpmOnboardingPage = new CustomerPpmOnboardingPage(customerPpmPage);
      await customerPpmOnboardingPage.signInForPPMWithMove(move);
    });

    test('Fill out About Form Date', async () => {
      await customerPpmOnboardingPage.navigateToAboutPageAndFillOutAboutFormDate();
    });
  });
});

test.describe('Entire PPM onboarding flow', () => {
  /** @type {CustomerPpmOnboardingPage} */
  let customerPpmOnboardingPage;

  forEachViewport(async ({ isMobile }) => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildDraftMoveWithPPMWithDepartureDate();
      customerPpmOnboardingPage = new CustomerPpmOnboardingPage(customerPpmPage);
      await customerPpmOnboardingPage.signInForPPMWithMove(move);
    });

    test('flows through happy path for existing shipment', async () => {
      await customerPpmOnboardingPage.navigateFromHomePageToExistingPPMDateAndLocationPage();
      await customerPpmOnboardingPage.submitsDateAndLocation();
      await customerPpmOnboardingPage.submitsEstimatedWeightsAndProGear();
      await customerPpmOnboardingPage.generalVerifyEstimatedIncentivePage({ isMobile });
      await customerPpmOnboardingPage.submitsAdvancePage({ addAdvance: true, isMobile });
      await customerPpmOnboardingPage.navigateToAgreementAndSign();
      await customerPpmOnboardingPage.submitMove();
      await customerPpmOnboardingPage.verifyManagePPMStepExistsAndBtnIsDisabled();
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
      await customerPpmOnboardingPage.verifyManagePPMStepExistsAndBtnIsDisabled();
    });
  });
});

test.describe('(MultiMove) Entire PPM onboarding flow', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  /** @type {CustomerPpmOnboardingPage} */
  let customerPpmOnboardingPage;

  forEachViewport(async ({ isMobile }) => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildDraftMoveWithPPMWithDepartureDate();
      customerPpmOnboardingPage = new CustomerPpmOnboardingPage(customerPpmPage);
      await customerPpmOnboardingPage.signInForPPMWithMove(move);
    });

    test('flows through happy path for existing shipment', async () => {
      await customerPpmOnboardingPage.navigateFromHomePageToExistingPPMDateAndLocationPage();
      await customerPpmOnboardingPage.submitsDateAndLocation();
      await customerPpmOnboardingPage.submitsEstimatedWeightsAndProGear();
      await customerPpmOnboardingPage.generalVerifyEstimatedIncentivePage({ isMobile });
      await customerPpmOnboardingPage.submitsAdvancePage({ addAdvance: true, isMobile });
      await customerPpmOnboardingPage.navigateToAgreementAndSign();
      await customerPpmOnboardingPage.submitMove();
      await customerPpmOnboardingPage.verifyManagePPMStepExistsAndBtnIsDisabled();
    });

    test('happy path with edits and backs', async () => {
      test.skip(true, 'Test fails at navigateFromHomePageToExistingPPMDateAndLocationPage()');
      await customerPpmOnboardingPage.navigateFromHomePageToExistingPPMDateAndLocationPage();

      await customerPpmOnboardingPage.submitAndVerifyUpdateDateAndLocation();

      await customerPpmOnboardingPage.submitsEstimatedWeightsAndProGear();
      await customerPpmOnboardingPage.verifyEstimatedWeightsAndProGear();

      await customerPpmOnboardingPage.verifyShipmentSpecificInfoOnEstimatedIncentivePage();
      await customerPpmOnboardingPage.generalVerifyEstimatedIncentivePage({ isMobile });

      await customerPpmOnboardingPage.submitsAdvancePage({ addAdvance: true, isMobile });

      await customerPpmOnboardingPage.navigateToAgreementAndSign();

      await customerPpmOnboardingPage.submitMove();
      await customerPpmOnboardingPage.verifyManagePPMStepExistsAndBtnIsDisabled();
    });
  });
});
