// @ts-check
const { test, expect } = require('../../utils/officeTest');

// The logic in QaeFlowPage is only used in this file, so keep the
// playwright test fixture in this file.
class QaeFlowPage {
  /**
   * @param {import('@playwright/test').Page} page
   * @param {import('../../utils/officeTest').OfficePage} officePage
   * @param {string} moveCode
   */
  constructor(page, officePage, moveCode) {
    this.page = page;
    this.officePage = officePage;
    this.moveCode = moveCode;
  }

  /**
   * search for and navigate to move, then QA tab
   */
  async searchForAndNavigateToMoveQATab() {
    await this.officePage.searchForAndNavigateToMove(this.moveCode);

    // Go to quality assurance tab
    await this.page.getByText('Quality assurance').click();
    await this.officePage.waitForLoading();
    expect(this.page.url()).toContain(`/moves/${this.moveCode}/evaluation-reports`);
    const tableCount = await this.page.locator('[data-testid="evaluationReportTable"]').count();
    expect(tableCount).toBeGreaterThan(0);
  }

  /**
   * create shipment report
   */
  async createShipmentReport() {
    // Click the shipment "Create report" button
    await this.page.locator('[data-testid="shipmentEvaluationCreate"]').click();

    // Wait for the fetch and an element of form to load
    await expect(this.page.locator('[data-testid="evaluationReportForm"]')).toHaveCount(1);
  }

  async createCounselingReport() {
    // Click the counseling "Create report" button
    await this.page.locator('[data-testid="counselingEvaluationCreate"]').click();

    await expect(this.page.locator('[data-testid="evaluationReportForm"]')).toHaveCount(1);
  }

  /**
   * open submission preview
   */
  async openSubmissionPreview() {
    // Click review for submit button
    await this.page.locator('[data-testid="reviewAndSubmit"]').click();
    await this.officePage.waitForLoading();

    // Make sure page preview model is open
    await expect(this.page.locator('[data-testid="EvaluationReportPreview"]')).toHaveCount(1);
    await expect(this.page.getByText('01 Oct 2022')).toBeVisible();
  }

  /**
   * submit report from the preview modal and verify
   */
  async submitReportFromPreviewAndVerify() {
    const previewPage = this.page.getByTestId('modal');
    // Click submit button in the modal (waits for button to be attached to DOM before clicking)
    await previewPage.getByRole('button', { name: 'Submit' }).click();
    await this.officePage.waitForLoading();

    // Check for success alert
    await expect(this.page.locator('.usa-alert__text')).toContainText('Your report has been successfully submitted');

    // Make sure thare is at least one report that has been submitted.
    const reportPage = this.page.locator('[data-testid="evaluationReportTable"]').locator('..');

    await expect(reportPage.locator('[data-testid="viewReport"]')).toContainText('View report');
    await expect(reportPage.locator('[data-testid="downloadReport"]')).toContainText('Download');

    // Should be back on the quality assurance tab
    expect(this.page.url()).toContain(`/moves/${this.moveCode}/evaluation-reports`);
    const tableCount = await this.page.locator('[data-testid="evaluationReportTable"]').count();
    expect(tableCount).toBeGreaterThan(0);
  }

  /**
   * save as draft
   */
  async saveAsDraft() {
    // Click save as draft button

    await Promise.all([this.page.waitForNavigation(), this.page.getByRole('button', { name: 'Save draft' }).click()]);

    // Should be back on the quality assurance tab
    expect(this.page.url()).toContain(`/moves/${this.moveCode}/evaluation-reports`);
    const tableCount = await this.page.locator('[data-testid="evaluationReportTable"]').count();
    expect(tableCount).toBeGreaterThan(0);

    // Should show a saved draft alert
    await expect(this.page.locator('.usa-alert__text')).toContainText('Your draft report has been saved');
  }

  // // Fills out the first page of the Evaluation Report Form providing basic content in minimal required fields
  async fillInForm() {
    await this.page.locator('input[name="inspectionDate"]').type('01 Oct 2022');
    await this.page.locator('input[name="inspectionDate"]').blur(); // Date of inspection
    // evaluation start and end times
    await this.page.locator('select[name="evalStartHour"]').selectOption({ label: '04' });
    await this.page.locator('select[name="evalStartMinute"]').selectOption({ label: '25' });
    await this.page.locator('select[name="evalEndHour"]').selectOption({ label: '12' });
    await this.page.locator('select[name="evalEndMinute"]').selectOption({ label: '38' });
    await this.page.locator('[data-testid="radio"] [for="dataReview"]').click(); // Evaluation type
    await this.page.locator('[data-testid="radio"] [for="origin"]').click(); // Evaluation location
    await this.page.locator('[data-testid="radio"] [for="noViolations"]').click(); // Violations observed
    await this.page.locator('textarea[name="remarks"]').type('This is a test evaluation report'); // Evaluation remarks
  }

  /**
   * change violations to yes and go to page two
   */
  async goToFormPageTwo() {
    // Change violations observed to 'yes'
    await this.page.locator('[data-testid="radio"] [for="yesViolations"]').click();

    // Click 'Next: Select Violations' button
    await this.page.locator('[data-testid="selectViolations"]').click();

    // Wait for form to save and switch over to page 2 of form
    await expect(this.page.locator('[data-testid="evaluationViolationsForm"]')).toHaveCount(1);
  }

  /**
   * Select the violoations
   * @param {string} category
   * @param {string} violation
   */
  async selectViolation(category, violation) {
    // Expand the category
    await this.page.getByTestId(`accordionButton_${category}Violation`).click();

    // Select the violation
    await this.page.locator(`input[name="${violation}"]`).locator('..').click();

    // Collapse the category so the page is back to original state
    await this.page.getByTestId(`accordionButton_${category}Violation`).click();
  }
}

test.describe('Quality Evaluation Report', () => {
  /** @type {QaeFlowPage} */
  let qaeFlowPage;
  let move;

  // setup qaeFlowpage for each test
  test.beforeEach(async ({ page, officePage }) => {
    move = await officePage.buildHHGMoveWithNTSAndNeedsSC();

    await officePage.signInAsNewQAECSRUser();
    qaeFlowPage = new QaeFlowPage(page, officePage, move.locator);
    await qaeFlowPage.searchForAndNavigateToMoveQATab();
  });

  test.describe('Happy Paths', () => {
    // Create new shipment report, fill out minimal fields, submit from 1st page (no violations)
    test('can complete a minimal shipment evaluation report from creation through submission', async ({ page }) => {
      // Create a new shipment report
      await qaeFlowPage.createShipmentReport();

      // Fill in the eval report form with minimal info required to submit
      await qaeFlowPage.fillInForm();

      // Review report for submission
      await qaeFlowPage.openSubmissionPreview();

      // Verify preview has correct sections/headers
      const previewPage = page.getByTestId('EvaluationReportPreview');
      await expect(previewPage.locator('h1')).toContainText('Shipment report');
      await expect(previewPage.locator('h2').getByText('Move information', { exact: true })).toBeVisible();
      await expect(previewPage.locator('h2').getByText('Evaluation report', { exact: true })).toBeVisible();
      await expect(previewPage.locator('h3').getByText('Violations', { exact: true })).toBeVisible();
      await expect(previewPage.locator('h3').getByText('QAE remarks', { exact: true })).toBeVisible();

      // Verify preview displays saved report data
      const evalReportPage = previewPage.locator('h2').getByText('Evaluation report', { exact: true }).locator('..');
      await expect(evalReportPage.locator('td').getByText('01 Oct 2022')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('Data review')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('Origin')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('No')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('This is a test evaluation report')).toBeVisible();

      // Submit report
      await qaeFlowPage.submitReportFromPreviewAndVerify();
    });

    // Create new counseling report, fill out minimal fields, submit
    // from 1st page (no violations)
    test('can complete a minimal counseling evaluation report from creation through submission', async ({ page }) => {
      // Create a new counseling report
      await qaeFlowPage.createCounselingReport();

      // Fill out the eval report form with minimal info required to submit
      await qaeFlowPage.fillInForm();

      // Review report for submission (saves draft)
      await qaeFlowPage.openSubmissionPreview();

      // Verify preview has correct sections/headers
      const previewPage = page.getByTestId('EvaluationReportPreview');
      await expect(previewPage.locator('h1')).toContainText('Counseling report');
      await expect(previewPage.locator('h2').getByText('Move information', { exact: true })).toBeVisible();
      await expect(previewPage.locator('h2').getByText('Evaluation report', { exact: true })).toBeVisible();
      await expect(previewPage.locator('h3').getByText('Violations', { exact: true })).toBeVisible();
      await expect(previewPage.locator('h3').getByText('QAE remarks', { exact: true })).toBeVisible();

      // Verify preview displays saved report data
      const evalReportPage = previewPage.locator('h2').getByText('Evaluation report', { exact: true }).locator('..');
      await expect(evalReportPage.locator('td').getByText('01 Oct 2022')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('No')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('This is a test evaluation report')).toBeVisible();

      // Submit report
      await qaeFlowPage.submitReportFromPreviewAndVerify();
    });

    // Create new report, fill out minimal fields with violations
    // observed = true, select violations on second page, submit
    test('can complete a minimal evaluation report with violations from creation through submission', async ({
      page,
    }) => {
      // Create a new shipment report
      await qaeFlowPage.createShipmentReport();

      // Fill out the eval report form with minimal info required to submit
      await qaeFlowPage.fillInForm();

      // Go to next page of form to select violations
      await qaeFlowPage.goToFormPageTwo();

      // Select a violation
      await qaeFlowPage.selectViolation('Storage', '1.2.6.16 Storage');

      // No serious violations
      await page.locator('[data-testid="radio"] [for="no"]').click();

      // Review report for submission
      await qaeFlowPage.openSubmissionPreview();

      // Verify preview has correct sections/headers
      const previewPage = page.getByTestId('EvaluationReportPreview');
      await expect(previewPage.getByTestId('noViolationsObserved')).not.toBeVisible();
      await expect(previewPage.getByText('1.2.6.16 Storage')).toBeVisible();

      await expect(previewPage.locator('h1')).toContainText('Shipment report');
      await expect(previewPage.locator('h2').getByText('Move information', { exact: true })).toBeVisible();
      await expect(previewPage.locator('h2').getByText('Evaluation report', { exact: true })).toBeVisible();
      await expect(previewPage.locator('h3').getByText('Violations', { exact: true })).toBeVisible();
      await expect(previewPage.locator('h3').getByText('QAE remarks', { exact: true })).toBeVisible();

      // Verify preview displays saved report data
      const evalReportPage = previewPage.locator('h2').getByText('Evaluation report', { exact: true }).locator('..');
      await expect(evalReportPage.locator('small').getByText('Provide adequate storage facilities')).toBeVisible();
      await expect(evalReportPage.locator('td').getByText('01 Oct 2022')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('Data review')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('Origin')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('No')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('This is a test evaluation report')).toBeVisible();

      // Submit report
      await qaeFlowPage.submitReportFromPreviewAndVerify();
    });

    // Create new report, fill out all conditional fields on 1st page
    // fields, on second page select correct violations to
    // display/fill out date fields, serious violations = true, submit
    test('can complete an evaluation report with all fields populated, including conditionally displayed fields', async ({
      page,
    }) => {
      // Create a new shipment report
      await qaeFlowPage.createShipmentReport();

      // Fill out the eval report form with minimal info required to submit
      await qaeFlowPage.fillInForm();

      await page.locator('[data-testid="radio"] [for="physical"]').click();
      await page.locator('[data-testid="radio"] [for="origin"]').click();

      // Time departed needed for physical inspection
      await page.locator('select[name="timeDepartHour"]').selectOption({ label: '02' });
      await page.locator('select[name="timeDepartMinute"]').selectOption({ label: '15' });

      // Evaluation location, has up to 3 conditional fields displayed dependent upon selection
      await page.locator('[data-testid="radio"] [for="destination"]').click();
      await page.locator('input[name="observedShipmentDeliveryDate"]').type('02 Oct 2022');
      await page.locator('input[name="observedShipmentDeliveryDate"]').blur(); // Observed delivery date
      await page.locator('[data-testid="radio"] [for="origin"]').click(); // Evaluation location
      await page.locator('input[name="observedShipmentPhysicalPickupDate"]').type('02 Oct 2022');
      await page.locator('input[name="observedShipmentPhysicalPickupDate"]').blur(); // Observed pickup date
      await page.locator('[data-testid="radio"] [for="other"]').click(); // Evaluation location
      await page.locator('textarea[name="otherEvaluationLocation"]').type('This is a test other location text');

      // Go to next page of form to select violations
      await qaeFlowPage.goToFormPageTwo();

      // Expand violation categories that have violaions with KPIs
      await qaeFlowPage.selectViolation('Counseling', '1.2.5.3 Scheduling');
      await qaeFlowPage.selectViolation('ShipmentSchedule', '1.2.6.7 Pickup');
      await qaeFlowPage.selectViolation('LossDamage', '1.2.7.2.2 Claims Settlement');
      await qaeFlowPage.selectViolation('ShipmentSchedule', '1.2.6.15 Delivery');

      // Fill out date fields for violations with KPIs
      await page.locator('input[name="observedClaimsResponseDate"]').type('03 Oct 2022');
      await page.locator('input[name="observedClaimsResponseDate"]').blur(); // Observed claims response date
      await page.locator('input[name="observedPickupDate"]').type('04 Oct 2022');
      await page.locator('input[name="observedPickupDate"]').blur(); // Observed pickup date
      await page.locator('input[name="observedDeliveryDate"]').type('05 Oct 2022');
      await page.locator('input[name="observedDeliveryDate"]').blur(); // Observed delivery date
      await page.locator('input[name="observedPickupSpreadStartDate"]').type('06 Oct 2022');
      await page.locator('input[name="observedPickupSpreadStartDate"]').blur(); // Observed pickup spread start date
      await page.locator('input[name="observedPickupSpreadEndDate"]').type('07 Oct 2022');
      await page.locator('input[name="observedPickupSpreadEndDate"]').blur(); // Observed pickup spread end date

      // Serious violations
      await page.locator('[data-testid="radio"] [for="yes"]').click();
      await page.locator('[data-testid="textarea"]').type('This is a test serious violation text');

      // Review report for submission
      await qaeFlowPage.openSubmissionPreview();

      const previewPage = page.getByTestId('EvaluationReportPreview');
      await expect(previewPage.getByText('Move information')).toBeVisible();
      await expect(previewPage.getByText('Yes')).toBeVisible();

      await expect(previewPage.locator('h1')).toContainText('Shipment report');
      await expect(previewPage.locator('h2').getByText('Evaluation report', { exact: true })).toBeVisible();
      await expect(previewPage.locator('h3').getByText('Violations', { exact: true })).toBeVisible();
      await expect(previewPage.locator('h3').getByText('QAE remarks', { exact: true })).toBeVisible();

      // Verify preview displays saved report data
      const evalReportPage = previewPage.locator('h2').getByText('Evaluation report', { exact: true }).locator('..');
      await expect(evalReportPage.locator('td').getByText('01 Oct 2022')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('Physical')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('Other')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('This is a test other location text')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('This is a test evaluation report')).toBeVisible();

      //  kpi dates
      await expect(evalReportPage.locator('dd').getByText('03 Oct 2022')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('04 Oct 2022')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('05 Oct 2022')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('06 Oct 2022')).toBeVisible();
      await expect(evalReportPage.locator('dd').getByText('07 Oct 2022')).toBeVisible();

      // Submit report
      await qaeFlowPage.submitReportFromPreviewAndVerify();
    });

    test('can edit a saved draft', async ({ page }) => {
      // Create a new shipment report
      await qaeFlowPage.createShipmentReport();

      // Fill out the eval report form with minimal info
      await qaeFlowPage.fillInForm();

      // Save draft
      await qaeFlowPage.saveAsDraft();

      // Try to edit a draft
      await page.locator('[data-testid="editReport"]').first().click();

      // Verify the form to edit is displayed
      await expect(page.locator('[data-testid="evaluationReportForm"]')).toBeVisible();
    });

    test('can delete a draft', async ({ page }) => {
      // Create a new shipment report
      await qaeFlowPage.createShipmentReport();

      // Save draft (reroutes back to reports list)
      await qaeFlowPage.saveAsDraft();

      // Delete draft
      await page.locator('[data-testid="deleteReport"]').first().click();

      // Confirm delete modal

      // ahobson 2023-02-04 - somehow there are two modals on the
      // page? confirmed with developer tools. We need to interact
      // with the 2nd one. Not sure why
      const modal = page.getByTestId('modal').nth(1);
      await expect(modal.getByText('Are you sure you want to delete this report?')).toBeVisible();
      await modal.getByText('Yes, delete').click();
      await expect(page.getByTestId('modal')).toHaveCount(0);

      // Verify deletion
      await expect(page.getByText('Your draft report has been saved')).toHaveCount(0);
      await expect(page.locator('.usa-alert__text')).toContainText('Your report has been deleted');
    });

    test('can view a submitted draft', async ({ page }) => {
      // Create and submit a new report
      await qaeFlowPage.createShipmentReport();
      await qaeFlowPage.fillInForm();
      await qaeFlowPage.openSubmissionPreview();
      await qaeFlowPage.submitReportFromPreviewAndVerify();

      // View saved report
      await page.locator('[data-testid="viewReport"]').first().click();

      // Verify the 'view report' modal is shown
      await expect(page.getByText('Evaluation report', { exact: true })).toBeVisible();
      await expect(page.getByTestId('EvaluationReportPreview')).toBeVisible();
    });

    test('has option to download a submitted draft', async ({ page }) => {
      // Create and submit a new report
      await qaeFlowPage.createShipmentReport();
      await qaeFlowPage.fillInForm();
      await qaeFlowPage.openSubmissionPreview();
      await qaeFlowPage.submitReportFromPreviewAndVerify();

      // Verify download button is displayed
      await expect(page.getByTestId('downloadReport')).toBeVisible();
    });
  });

  test.describe('Save Draft Behavior', () => {
    test('can save report as draft from the first form page (no violations)', async ({ page }) => {
      // Create a new shipment report
      await qaeFlowPage.createShipmentReport();
      // Fill out the eval report form with minimal info
      await qaeFlowPage.fillInForm();
      // Save draft
      await qaeFlowPage.saveAsDraft();
      // Verify draft saved
      await expect(page.getByText('Your draft report has been saved')).toBeVisible();
    });

    test('can save report as draft from the second form page (has violations)', async ({ page }) => {
      // Create a new shipment report
      await qaeFlowPage.createShipmentReport();
      // Fill out the eval report form with minimal info
      await qaeFlowPage.fillInForm();
      // Go to next page of form to select violations
      await qaeFlowPage.goToFormPageTwo();
      // Select a violation
      await qaeFlowPage.selectViolation('Storage', '1.2.6.16 Storage');
      // Select no serious violations
      await page.locator('[data-testid="radio"] [for="no"]').click();
      // Save draft
      await qaeFlowPage.saveAsDraft();
      // Verify draft saved
      await expect(page.getByText('Your draft report has been saved')).toBeVisible();
    });
  });

  test.describe('Cancel Behavior', () => {
    // Create a report, click cancel before saving, verify/confim
    // modal, verify delete worked upon reroute
    test('prompts to delete if the report has not been saved since creation', async ({ page }) => {
      // Create a new shipment report
      await qaeFlowPage.createShipmentReport();
      // Cancel report
      await page.locator('[data-testid="cancelReport"]').click();
      // Verify cancel modal
      await expect(page.getByText('Are you sure you want to cancel this report?')).toBeVisible();
      await page.getByText('Yes, cancel').click();
      // Verify deletion
      await expect(page.locator('.usa-alert__text')).toContainText('Your report has been canceled');
    });

    // Create a report, save draft, click cancel, verify reroute w/out delete
    test('does not prompt to delete report after first save', async ({ page, officePage }) => {
      // Create a new shipment report
      await qaeFlowPage.createShipmentReport();
      // Fill out the eval report form with minimal info
      await qaeFlowPage.fillInForm();
      // Save draft
      await qaeFlowPage.saveAsDraft();
      // Edit draft report
      await page.getByTestId('editReport').last().click();
      await officePage.waitForLoading();

      const formCount = await page.getByTestId('evaluationReportForm').count();
      expect(formCount).toBeGreaterThan(0);
      // Click cancel
      await page.getByTestId('cancelReport').click();
      // No cancel report since not deleting
      await expect(page.getByText('Are you sure you want to cancel this report?')).not.toBeVisible();
      await officePage.waitForLoading();

      // Should have been rerouted to reports list
      expect(page.url()).toContain(`/moves/${qaeFlowPage.moveCode}/evaluation-reports`);
      const tableCount = await page.getByTestId('evaluationReportTable').count();
      expect(tableCount).toBeGreaterThan(0);
      // Verify no delete modal
      await expect(page.locator('.usa-alert__text')).not.toBeVisible();
    });

    // Create a report through 2nd page (violations),  click cancel, verify reroute w/out delete
    test('returns to reports list without deletion when canceled from the 2nd form page', async ({
      page,
      officePage,
    }) => {
      // Create a new shipment report
      await qaeFlowPage.createShipmentReport();
      // Fill out the eval report form with minimal info
      await qaeFlowPage.fillInForm();
      // Go to next page of form to select violations
      await qaeFlowPage.goToFormPageTwo();
      // Click cancel
      await page.getByTestId('cancelReport').click();
      // No cancel report since not deleting
      expect(page.getByText('Are you sure you want to cancel this report?')).not.toBeVisible();
      await officePage.waitForLoading();
      // Should have been rerouted to reports list
      expect(page.url()).toContain(`/moves/${qaeFlowPage.moveCode}/evaluation-reports`);
      const tableCount = await page.getByTestId('evaluationReportTable').count();
      expect(tableCount).toBeGreaterThan(0);
      // Verify no delete modal
      await expect(page.locator('.usa-alert__text')).not.toBeVisible();
    });
  });
});
