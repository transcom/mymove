// @ts-check
import { OfficePage } from '../../utils/office/officeTest';

/**
 * TioFlowPage test fixture
 *
 * The logic in TioFlowPage is only used in this file, so keep the
 * playwright test fixture in this file.
 * @extends OfficePage
 */
export class TioFlowPage extends OfficePage {
  /**
   * @param {OfficePage} officePage
   * @param {Object} move
   * @override
   */
  constructor(officePage, move) {
    super(officePage.page, officePage.request);
    this.move = move;
    this.moveLocator = move.locator;
  }
}

export default TioFlowPage;
