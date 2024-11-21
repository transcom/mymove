import { test, expect, OfficePage } from '../../utils/office/officeTest';

/**
 * GsrFlowPage test fixture
 *
 * The logic in GsrFlowPage is only used in this file, so keep the
 * playwright test fixture in this file.
 * @extends OfficePage
 */
class GsrFlowPage extends OfficePage {
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
   * search for and navigate to move, then QA tab
   */
  async searchForAndNavigateToMoveQATab() {
    await this.qaeSearchForAndNavigateToMove(this.moveLocator);
    await this.page.getByText('Quality assurance').click();
    await this.waitForLoading();
    expect(this.page.url()).toContain(`/moves/${this.moveLocator}/evaluation-reports`);
    await expect(this.page.getByTestId('evaluationReportTable').first()).toBeVisible();
  }

  async createShipmentReport() {
    await this.page.locator('[data-testid="shipmentEvaluationCreate"]').click();
    await expect(this.page.locator('[data-testid="evaluationReportForm"]')).toHaveCount(1);
  }

  async openSubmissionPreview() {
    await this.page.locator('[data-testid="reviewAndSubmit"]').click();
    await this.waitForLoading();
    await expect(this.page.locator('[data-testid="EvaluationReportPreview"]')).toHaveCount(1);
    await expect(this.page.getByText('01 Oct 2022')).toBeVisible();
  }

  async submitReportFromPreviewAndVerify() {
    const previewPage = this.page.getByTestId('modal');
    await previewPage.getByRole('button', { name: 'Submit' }).click();
    await this.waitForLoading();
    await expect(this.page.locator('.usa-alert__text')).toContainText('Your report has been successfully submitted');
    const reportPage = this.page.locator('[data-testid="evaluationReportTable"]').locator('..');
    await expect(reportPage.locator('[data-testid="viewReport"]')).toContainText('View report');
    await expect(reportPage.locator('[data-testid="downloadReport"]')).toContainText('Download');
    expect(this.page.url()).toContain(`/moves/${this.moveLocator}/evaluation-reports`);
    await expect(this.page.getByTestId('evaluationReportTable').first()).toBeVisible();
  }

  async fillInForm() {
    await this.page.locator('input[name="inspectionDate"]').fill('01 Oct 2022');
    await this.page.locator('input[name="inspectionDate"]').blur(); // Date of inspection
    await this.page.locator('select[name="evalStartHour"]').selectOption({ label: '04' });
    await this.page.locator('select[name="evalStartMinute"]').selectOption({ label: '25' });
    await this.page.locator('select[name="evalEndHour"]').selectOption({ label: '12' });
    await this.page.locator('select[name="evalEndMinute"]').selectOption({ label: '38' });
    await this.page.locator('[data-testid="radio"] [for="dataReview"]').click(); // Evaluation type
    await this.page.locator('[data-testid="radio"] [for="origin"]').click(); // Evaluation location
    await this.page.locator('[data-testid="radio"] [for="noViolations"]').click(); // Violations observed
    await this.page.locator('textarea[name="remarks"]').fill('This is a test evaluation report'); // Evaluation remarks
  }

  async goToFormPageTwo() {
    await this.page.locator('[data-testid="radio"] [for="yesViolations"]').click();
    await this.page.locator('[data-testid="selectViolations"]').click();
    await expect(this.page.locator('[data-testid="evaluationViolationsForm"]')).toHaveCount(1);
  }

  async selectViolation(category, violation) {
    await this.page.getByTestId(`accordionButton_${category}Violation`).click();
    await this.page.locator(`input[name="${violation}"]`).locator('..').click();
    await this.page.getByTestId(`accordionButton_${category}Violation`).click();
  }
}

const gsrEnabled = process.env.FEATURE_FLAG_GSR_ROLE;

test.describe('GSR User Flow', () => {
  /** @type {GsrFlowPage} */
  let gsrFlowPage;
  let move;

  test.beforeEach(async ({ officePage, page }) => {
    move = await officePage.testHarness.buildHHGMoveWithNTSAndNeedsSC();

    // First we must sign in as a QAE and create a report
    await officePage.signInAsNewQAEUser();
    gsrFlowPage = new GsrFlowPage(officePage, move.locator);
    await gsrFlowPage.searchForAndNavigateToMoveQATab();
    await gsrFlowPage.createShipmentReport();
    await gsrFlowPage.fillInForm();
    await gsrFlowPage.goToFormPageTwo();
    // Selecting a violation
    await gsrFlowPage.selectViolation('Storage', '1.2.6.16 Storage');
    await page.locator('[data-testid="radio"] [for="yes"]').click();
    await page.locator('textarea[name="seriousIncidentDesc"]').fill('This is a serious incident description');
    await gsrFlowPage.openSubmissionPreview();
    await gsrFlowPage.submitReportFromPreviewAndVerify();

    // now we sign in as a GSR
    await page.getByText('Sign out').click();
    await page.waitForURL('**/sign-in');

    await officePage.signInAsNewGSRUser();
    await gsrFlowPage.searchForAndNavigateToMoveQATab();
  });

  test.describe('Leave Appeal Decision', () => {
    test.skip(gsrEnabled === 'false', 'Skip if the GSR flag is off.');
    test('QAE creates an evaluation report and a GSR user can leave remarks', async ({ page }) => {
      await page.getByTestId('viewReport').click();

      // adding appeal to violation
      const addViolationAppealBtn = page.getByTestId('addViolationAppealBtn');
      await expect(addViolationAppealBtn).toBeVisible();
      await addViolationAppealBtn.click();
      await expect(page.getByRole('heading', { name: 'Leave Appeal Decision' })).toBeVisible();
      const remarksInput = page.getByTestId('addAppealRemarks');
      await expect(remarksInput).toBeVisible();
      await page.locator('textarea[name="remarks"]').fill('These are some appeal remarks for a violation');
      await expect(page.getByText('Sustained')).toBeVisible();
      await expect(page.getByText('Rejected')).toBeVisible();
      await page.getByText('Sustained').click();
      await page.getByRole('button', { name: 'Save' }).click();

      const showAppealsBtn = page.getByRole('button', { name: 'Show appeals' });
      await expect(showAppealsBtn).toBeVisible();
      await showAppealsBtn.click();

      await expect(page.getByText('Sustained')).toBeVisible();
      await expect(page.getByText('These are some appeal remarks for a violation')).toBeVisible();

      // adding serious incident appeal
      const addSeriousIncidentAppealBtn = page.getByTestId('addSeriousIncidentAppealBtn');
      await expect(addSeriousIncidentAppealBtn).toBeVisible();
      await addSeriousIncidentAppealBtn.click();
      await expect(page.getByRole('heading', { name: 'Leave Appeal Decision' })).toBeVisible();
      await expect(remarksInput).toBeVisible();
      await page.locator('textarea[name="remarks"]').fill('These are some appeal remarks for a serious incident');
      await page.getByText('Rejected').click();
      await page.getByRole('button', { name: 'Save' }).click();

      await expect(showAppealsBtn).toBeVisible();
      await showAppealsBtn.click();

      await expect(page.getByText('Rejected')).toBeVisible();
      await expect(page.getByText('These are some appeal remarks for a serious incident')).toBeVisible();
    });
  });
});
