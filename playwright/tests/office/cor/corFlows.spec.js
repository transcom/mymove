/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect, OfficePage } from '../../utils/office/officeTest';

/**
 * CorFlowPage test fixture
 *
 * The logic in CorFlowPage is only used in this file, so keep the
 * playwright test fixture in this file.
 * @extends OfficePage
 */
class CorFlowPage extends OfficePage {
  /**
   * @param {OfficePage} officePage
   * @param {string} moveLocator
   * @override
   */
  constructor(officePage, moveLocator) {
    super(officePage.page, officePage.request);
    this.moveLocator = moveLocator;
  }

  /**
   * search for and navigate to move
   */
  async searchForAndNavigateToMove() {
    await this.corSearchForAndNavigateToMove(this.moveLocator);
  }
}

test.describe('Contracting Officer Representative', () => {
  /** @type {CorFlowPage} */
  let corFlowPage;
  let move;

  // setup CorFlowPage for each test
  test.beforeEach(async ({ officePage }) => {
    move = await officePage.testHarness.buildHHGMoveWithNTSAndNeedsSC();

    await officePage.signInAsNewCORUser();
    corFlowPage = new CorFlowPage(officePage, move.locator);
    await corFlowPage.searchForAndNavigateToMove();
  });

  test.describe('Happy Paths', () => {
    test('can log in, search for a move, and default to the Move Task Order tab', async ({ page }) => {
      // Make sure we default to where we should as a COR
      await expect(page.url()).toContain(`/moves/${move.locator}/mto`);
    });
  });
});
