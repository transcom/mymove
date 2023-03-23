/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { expect, test } from './customerPpmTestFixture';

test.describe('PPM Onboarding - Add Estimated  Weight and Pro-gear', () => {
  test.beforeEach(async ({ customerPpmPage }) => {
    const move = await customerPpmPage.testHarness.buildUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights();
    await customerPpmPage.signInForPPMWithMove(move);
    await customerPpmPage.navigateFromHomePageToExistingPPMDateAndLocationPage();
    await customerPpmPage.navigateFromDateAndLocationPageToEstimatedWeightsPage();
  });

  test('doesn’t allow SM to progress if form is in an invalid state', async ({ page }) => {
    await expect(page.getByRole('heading', { name: 'Estimated weight' })).toBeVisible();
    await expect(page).toHaveURL(/\/estimated-weight/);
    await expect(page.locator('p[class="usa-alert__text"]')).toContainText(
      'Total weight allowance for your move: 5,000 lbs',
    );

    // missing required weight
    await page.locator('input[name="estimatedWeight"]').clear();
    const errorMessage = page.locator('[class="usa-error-message"]');
    await expect(errorMessage).toContainText('Required');
    await expect(page.locator('[class="usa-error-message"] + div').locator('input')).toHaveAttribute(
      'id',
      'estimatedWeight',
    );

    // estimated weight violates min
    await page.locator('input[name="estimatedWeight"]').type('0');
    await page.locator('input[name="estimatedWeight"]').blur();
    await expect(errorMessage).toContainText('Enter a weight greater than 0 lbs');
    await expect(page.locator('[class="usa-error-message"] + div').locator('input')).toHaveAttribute(
      'id',
      'estimatedWeight',
    );

    await page.locator('input[name="estimatedWeight"]').type('500');
    await page.locator('input[name="estimatedWeight"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // a warning is displayed when estimated weight is greater than the SM's weight allowance
    await page.locator('input[name="estimatedWeight"]').clear();
    await page.locator('input[name="estimatedWeight"]').type('17000');
    const warningMessage = page.locator('[data-testid="textInputWarning"]');
    await expect(warningMessage).toContainText(
      'This weight is more than your weight allowance. Talk to your counselor about what that could mean for your move.',
    );
    await expect(page.locator('[data-testid="textInputWarning"] + div').locator('input')).toHaveAttribute(
      'id',
      'estimatedWeight',
    );
    await page.locator('input[name="estimatedWeight"]').clear();
    await page.locator('input[name="estimatedWeight"]').type('500');
    await expect(warningMessage).not.toBeVisible();

    // pro gear violates max
    await page.locator('label[for="hasProGearYes"]').click();

    await page.locator('input[name="proGearWeight"]').clear();
    await page.locator('input[name="proGearWeight"]').type('5000');
    await page.locator('input[name="proGearWeight"]').blur();
    await expect(errorMessage).toContainText('Enter a weight 2,000 lbs or less');
    await page.locator('input[name="proGearWeight"]').clear();
    await page.locator('input[name="proGearWeight"]').type('500');
    await page.locator('input[name="proGearWeight"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // When hasProGear is true show error if either personal or spouse pro gear isn't specified
    await page.locator('input[name="proGearWeight"]').clear();
    await page.locator('input[name="proGearWeight"]').blur();
    await expect(errorMessage).toContainText(
      "Enter a weight into at least one pro-gear field. If you won't have pro-gear, select No above.",
    );
    await page.locator('input[name="proGearWeight"]').clear();
    await page.locator('input[name="proGearWeight"]').type('500');
    await page.locator('input[name="proGearWeight"]').blur();
    await expect(errorMessage).not.toBeVisible();

    // spouse pro gear max violation
    await page.locator('input[name="spouseProGearWeight"]').clear();
    await page.locator('input[name="spouseProGearWeight"]').type('1000');
    await page.locator('input[name="spouseProGearWeight"]').blur();
    await expect(errorMessage).toContainText('Enter a weight 500 lbs or less');

    await page.locator('input[name="spouseProGearWeight"]').clear();
    await page.locator('input[name="spouseProGearWeight"]').type('100');
    await page.locator('input[name="spouseProGearWeight"]').blur();
    await expect(errorMessage).not.toBeVisible();
  });

  test('can continue to next page', async ({ customerPpmPage }) => {
    await customerPpmPage.submitsEstimatedWeights();
  });

  test('can continue to next page with progear added', async ({ customerPpmPage }) => {
    await customerPpmPage.submitsEstimatedWeightsAndProGear();
  });
});
