/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { expect, test, forEachViewport } from './customerPpmTestFixture';

test.describe('Progear', () => {
  test.skip(multiMoveEnabled === 'true', 'Skip if MultiMove workflow is enabled.');
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

test.describe('Progear', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');

  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPM();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.clickOnGoToMoveButton();
      await customerPpmPage.navigateToAboutPage();
      await customerPpmPage.navigateFromWeightTicketPageBadUpload();
      await customerPpmPage.navigateFromCloseoutReviewPageToAddProGearPage();
    });

    test('Test upload of wrong xlsx file to progear page', async ({ customerPpmPage }) => {
      await customerPpmPage.selectMeProGear();
      await customerPpmPage.submitIncorrectXlsxFileForProGear();
    });
  });
});

test.describe('(MultiMove) Progear', () => {
  test.skip(multiMoveEnabled === 'false', 'Skip if MultiMove workflow is not enabled.');
  forEachViewport(async () => {
    test.beforeEach(async ({ customerPpmPage }) => {
      const move = await customerPpmPage.testHarness.buildApprovedMoveWithPPMProgearWeightTicket();
      await customerPpmPage.signInForPPMWithMove(move);
      await customerPpmPage.clickOnGoToMoveButton();
      await customerPpmPage.navigateToPPMReviewPageWithCompletePPM();
      await customerPpmPage.navigateFromCloseoutReviewPageToProGearPage();
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
