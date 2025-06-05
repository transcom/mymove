/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { expect, test, forEachViewport } from './customerPpmTestFixture';

test.describe('Workflow About Your PPM', () => {
  test.beforeEach(async ({ customerPpmPage }) => {
    const move = await customerPpmPage.testHarness.buildUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights();
    await customerPpmPage.signInForPPMWithMove(move);
    await customerPpmPage.clickOnGoToMoveButton();
  });

  test('does not allow SM to progress if form is in an invalid state', async ({ page }) => {
    await page.getByTestId('editShipmentButton').click();
    await page.waitForURL(/\/moves\/[\d|a-z|-]+\/shipments\/[\d|a-z|-]+\/.*/);
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await page.waitForURL(/\/moves\/[\d|a-z|-]+\/shipments\/[\d|a-z|-]+\/estimated-weight$/);
    await page.getByRole('button', { name: 'Save & Continue' }).click();
    await page.waitForURL(/\/moves\/[\d|a-z|-]+\/shipments\/[\d|a-z|-]+\/estimated-incentive$/);
    await page.getByRole('button', { name: 'Next' }).click();

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
    await advanceInput.fill('1');
    await advanceInput.blur();
    await expect(page.locator('[class="usa-error-message"]')).not.toBeVisible();

    // advance violates min
    await advanceInput.clear();
    await advanceInput.fill('0');
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
    await advanceInput.fill('1');
    await advanceInput.blur();
    await expect(errorMessage).not.toBeVisible();

    // advance violates max (over 60% of incentive)
    await advanceInput.clear();
    await advanceInput.fill('600001');
    await advanceInput.blur();
    await expect(errorMessage).toContainText('Enter an amount');
    await expect(page.locator('[class="usa-error-message"] + div').locator('input')).toHaveAttribute(
      'id',
      'advanceAmountRequested',
    );
    await expect(saveButton).toBeDisabled();
    await advanceInput.clear();
    await advanceInput.fill('1');
    await advanceInput.blur();
    await expect(errorMessage).not.toBeVisible();
  });

  forEachViewport(async ({ isMobile }) => {
    [true, false].forEach((addAdvance) => {
      const advanceText = addAdvance ? 'request' : 'opt to not receive';
      test(`can ${advanceText} an advance`, async ({ customerPpmPage }) => {
        await customerPpmPage.navigateFromMoveHomeToAdvances();
        await customerPpmPage.submitsAdvancePage({ addAdvance, isMobile });
      });
    });
  });
});
