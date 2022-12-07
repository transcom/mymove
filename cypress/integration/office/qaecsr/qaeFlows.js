import { QAECSROfficeUserType } from '../../../support/constants';
import { searchForAndNavigateToMove } from './qaeCSRIntegrationUtils';

const moveCode = 'QAEHLP';

before(() => {
  cy.prepareOfficeApp();
});

beforeEach(() => {
  cy.intercept('**/ghc/v1/swagger.yaml').as('getGHCClient');
  cy.intercept('**/ghc/v1/queues/moves?page=1&perPage=20&sort=status&order=asc').as('getSortedOrders');
  cy.intercept('**/ghc/v1/move/**').as('getMoves');
  cy.intercept('GET', '**/ghc/v1/orders/**').as('getOrders');
  cy.intercept('**/ghc/v1/move_task_orders/**/mto_shipments').as('getMTOShipments');
  cy.intercept('**/ghc/v1/moves/**/shipment-evaluation-reports-list').as('getShipmentEvaluationReports');
  cy.intercept('**/ghc/v1/moves/**/counseling-evaluation-reports-list').as('getCounselingEvaluationReports');
  cy.intercept('**/ghc/v1/move_task_orders/**/mto_service_items').as('getMTOServiceItems');
  cy.intercept('GET', '**/ghc/v1/evaluation-reports/**').as('getEvaluationReport');
  cy.intercept('PUT', '**/ghc/v1/evaluation-reports/**').as('saveEvaluationReport');
  cy.intercept('DELETE', '**/ghc/v1/evaluation-reports/**').as('deleteEvaluationReport');
  cy.intercept('POST', '**/ghc/v1/evaluation-reports/**/submit').as('submitEvaluationReport');
  cy.intercept('GET', '**/ghc/v1/report-violations/**').as('getReportViolationsByReportID');
  cy.intercept('POST', '**/ghc/v1/report-violations/**').as('associateReportViolations');
  cy.intercept('**/ghc/v1/moves/search').as('getSearchResults');

  const userId = '2419b1d6-097f-4dc4-8171-8f858967b4db';
  cy.apiSignInAsUser(userId, QAECSROfficeUserType);

  // Go to a move
  searchForAndNavigateToMove();

  // Go to quality assurance tab
  cy.contains('Quality assurance').click();
  cy.url().should('include', `/moves/${moveCode}/evaluation-reports`);
  cy.wait(['@getMoves', '@getMTOShipments', '@getShipmentEvaluationReports', '@getCounselingEvaluationReports']);
  cy.get('[data-testid="evaluationReportTable"]').should('exist');
});

const createShipmentReport = () => {
  // Click the shipment "Create report" button
  cy.get('[data-testid="shipmentEvaluationCreate"]').click();

  // Wait for the fetch and an element of form to load
  cy.wait(['@getEvaluationReport']);
  cy.get('[data-testid="evaluationReportForm"]', { timeout: 10000 }).should('exist');
};

const createCounselingReport = () => {
  // Click the counseling "Create report" button
  cy.get('[data-testid="counselingEvaluationCreate"]').click();

  // Wait for the fetch and an element of form to load
  cy.wait(['@getEvaluationReport']);
  cy.get('[data-testid="evaluationReportForm"]', { timeout: 10000 }).should('exist');
};

const openSubmissionPreview = (isViolationsPage) => {
  // Click review for submit button
  cy.get('[data-testid="reviewAndSubmit"]').click();

  // Wait for the fetch and an element of form to load
  cy.wait(['@saveEvaluationReport', '@getEvaluationReport']);

  // Need to wait for violations associations if on the violations page
  if (isViolationsPage) cy.wait(['@associateReportViolations', '@getReportViolationsByReportID']);

  // Make sure page preview model is open
  cy.get('[data-testid="EvaluationReportPreview"]').should('exist');
  cy.contains('01 Oct 2022').should('exist');
};

const submitReportFromPreview = () => {
  // Click submit button in the modal (waits for button to be attached to DOM before clicking)
  cy.get('[data-testid="modalSubmitButton"]')
    .should(($el) => {
      expect(Cypress.dom.isDetached($el)).to.eq(false);
    })
    .click({ force: true });

  // Wait for the submit
  cy.wait(['@submitEvaluationReport']).then((intercept) => {
    const { statusCode } = intercept.response;
    expect(statusCode).to.eq(204);
  });

  // Should be back on the quality assurance tab
  cy.url().should('include', `/moves/${moveCode}/evaluation-reports`);
  cy.wait(['@getMoves', '@getMTOShipments', '@getShipmentEvaluationReports', '@getCounselingEvaluationReports']);
  cy.get('[data-testid="evaluationReportTable"]').should('exist');
};

const saveAsDraft = () => {
  // Click save as draft button
  cy.get('button').contains('Save draft').click();

  // Wait for the save
  cy.wait(['@saveEvaluationReport']);

  // Should be back on the quality assurance tab
  cy.url().should('include', `/moves/${moveCode}/evaluation-reports`);
  cy.wait(['@getMoves', '@getMTOShipments', '@getShipmentEvaluationReports', '@getCounselingEvaluationReports']);
  cy.get('[data-testid="evaluationReportTable"]').should('exist');

  // Should show a saved draft alert
  cy.get('.usa-alert__text').contains('Your draft report has been saved');
};

// Verify that a report is submitted (kinda hard to do w/out making a race condition. Other tests may create reports on this move and the report IDs are generated on create. Also, we can not delete submitted reports.)
const verifySubmittedReportPresent = () => {
  // Check for success alert
  cy.get('.usa-alert__text').contains('Your report has been successfully submitted');

  // Make sure thare is at least one report that has been submitted.
  cy.get('[data-testid="evaluationReportTable"]')
    .parent()
    .within(() => {
      cy.get('[data-testid="viewReport"]').contains('View report');
      cy.get('[data-testid="downloadReport"]').contains('Download');
    });
};

// Fills out the first page of the Evaluation Report Form providing basic content in minimal required fields
const fillInForm = () => {
  cy.get('input[name="inspectionDate"]').clear().type('01 Oct 2022').blur(); // Date of inspection
  // evaluation start and end times
  cy.get('select[name="evalStartHour"]').select('04').blur();
  cy.get('select[name="evalStartMinute"]').select('25').blur();
  cy.get('select[name="evalEndHour"]').select('12').blur();
  cy.get('select[name="evalEndMinute"]').select('38').blur();
  cy.get('[data-testid="radio"] [for="dataReview"]').click(); // Evaluation type
  cy.get('[data-testid="radio"] [for="origin"]').click(); // Evaluation location
  cy.get('[data-testid="radio"] [for="noViolations"]').click(); // Violations observed
  cy.get('textarea[name="remarks"]').type('This is a test evaluation report'); // Evaluation remarks
};

const goToFormPageTwo = () => {
  // Change violations observed to 'yes'
  cy.get('[data-testid="radio"] [for="yesViolations"]').click();

  // Click 'Next: Select Violations' button
  cy.get('[data-testid="selectViolations"]').click();

  // Wait for form to save and switch over to page 2 of form
  cy.wait(['@saveEvaluationReport', '@getEvaluationReport']);
  cy.get('[data-testid="evaluationViolationsForm"]').should('exist');
};

const selectViolation = (category, violation) => {
  // Expand the category
  cy.get(`[data-testid=accordionButton_${category}Violation]`).click();

  // Select the violation
  cy.get(`input[name="${violation}"]`).click({ force: true });
};

describe('Quality Evaluation Report', () => {
  describe('Happy Paths', () => {
    // Create new shipment report, fill out minimal fields, submit from 1st page (no violations)
    it('can complete a minimal shipment evaluation report from creation through submission', () => {
      // Create a new shipment report
      createShipmentReport();

      // Fill in the eval report form with minimal info required to submit
      fillInForm();

      // Review report for submission
      openSubmissionPreview();

      // Veryify preview has correct sections/headers
      cy.get('h1').contains('Shipment report');
      cy.get('h2').contains('Move information');
      cy.get('h2').contains('Evaluation report').parent().as('report');
      cy.get('h3').contains('Violations');
      cy.get('h3').contains('QAE remarks');

      // Verify preview displays saved report data
      cy.get('@report').within(() => {
        cy.get('td').contains('01 Oct 2022');
        cy.get('dd').contains('Data review');
        cy.get('dd').contains('Origin');
        cy.get('dd').contains('No');
        cy.get('dd').contains('This is a test evaluation report');
      });

      // Submit report
      submitReportFromPreview();

      // Verify submission
      verifySubmittedReportPresent();
    });

    // Create new counseling report, fill out minimal fields, submit from 1st page (no violations)
    it('can complete a minimal counseling evaluation report from creation through submission', () => {
      // Create a new counseling report
      createCounselingReport();

      // Fill out the eval report form with minimal info required to submit
      fillInForm();

      // Review report for submission (saves draft)
      openSubmissionPreview();

      // Veryify preview has correct sections/headers
      cy.get('h1').contains('Counseling report');
      cy.get('h2').contains('Move information');
      cy.get('h2').contains('Evaluation report').parent().as('report');
      cy.get('h3').contains('Violations');
      cy.get('h3').contains('QAE remarks');

      // Verify preview displays saved report data
      cy.get('@report').within(() => {
        cy.get('td').contains('01 Oct 2022');
        cy.get('dd').contains('No');
        cy.get('dd').contains('This is a test evaluation report');
      });

      // Submit report
      submitReportFromPreview();

      // Verify submission
      verifySubmittedReportPresent();
    });

    // Create new report, fill out minimal fields with violations observed = true, select violations on second page, submit
    it('can complete a minimal evaluation report with violations from creation through submission', () => {
      // Create a new shipment report
      createShipmentReport();

      // Fill out the eval report form with minimal info required to submit
      fillInForm();

      // Go to next page of form to select violations
      goToFormPageTwo();

      // Select a violation
      selectViolation('Storage', '1.2.6.16 Storage');

      // No serious violations
      cy.get('[data-testid="radio"] [for="no"]').click();

      // Review report for submission
      openSubmissionPreview();

      // Verify preview
      cy.get("[data-testid='noViolationsObserved").should('not.exist');
      cy.contains('1.2.6.16 Storage').should('exist');

      // Veryify preview has correct sections/headers
      cy.get('h1').contains('Shipment report');
      cy.get('h2').contains('Move information');
      cy.get('h2').contains('Evaluation report').parent().as('report');
      cy.get('h3').contains('Violations');
      cy.get('h3').contains('QAE remarks');

      // Verify preview displays saved report data
      cy.get('@report').within(() => {
        cy.get('dt').contains('Violations observed');
        cy.get('small').contains('Provide adequate storage facilities');
        cy.get('td').contains('01 Oct 2022');
        cy.get('dd').contains('Data review');
        cy.get('dd').contains('Origin');
        cy.get('dd').contains('No');
        cy.get('dd').contains('This is a test evaluation report');
      });

      // Submit report
      submitReportFromPreview();

      verifySubmittedReportPresent();
    });

    // Create new report, fill out all conditional fields on 1st page fields, on second page select correct violations to display/fill out date fields, serious violations = true, submit
    it('can complete an evaluation report with all fields populated, including conditionally displayed fields', () => {
      // Create a new shipment report
      createShipmentReport();

      // Fill out the eval report form with minimal info required to submit
      fillInForm();

      cy.get('[data-testid="radio"] [for="physical"]').click();
      cy.get('[data-testid="radio"] [for="origin"]').click();

      // Time departed needed for physical inspection
      cy.get('select[name="timeDepartHour"]').select('02').blur();
      cy.get('select[name="timeDepartMinute"]').select('15').blur();

      // Evaluation location, has up to 3 conditional fields displayed dependent upon selection
      cy.get('[data-testid="radio"] [for="destination"]').click();
      cy.get('input[name="observedShipmentDeliveryDate"]').clear().type('02 Oct 2022').blur(); // Observed delivery date
      cy.get('[data-testid="radio"] [for="origin"]').click(); // Evaluation location
      cy.get('input[name="observedShipmentPhysicalPickupDate"]').clear().type('02 Oct 2022').blur(); // Observed pickup date
      cy.get('[data-testid="radio"] [for="other"]').click(); // Evaluation location
      cy.get('textarea[name="otherEvaluationLocation"]').type('This is a test other location text');

      // Go to next page of form to select violations
      goToFormPageTwo();

      // Expand violation categories that have violaions with KPIs
      selectViolation('Counseling', '1.2.5.3 Scheduling');
      selectViolation('ShipmentSchedule', '1.2.6.7 Pickup');
      selectViolation('LossDamage', '1.2.7.2.2 Claims Settlement');
      selectViolation('ShipmentSchedule', '1.2.6.15 Delivery');

      // Fill out date fields for violations with KPIs
      cy.get('input[name="observedClaimsResponseDate"]').clear().type('03 Oct 2022').blur(); // Observed claims response date
      cy.get('input[name="observedPickupDate"]').clear().type('04 Oct 2022').blur(); // Observed pickup date
      cy.get('input[name="observedDeliveryDate"]').clear().type('05 Oct 2022').blur(); // Observed delivery date
      cy.get('input[name="observedPickupSpreadStartDate"]').clear().type('06 Oct 2022').blur(); // Observed pickup spread start date
      cy.get('input[name="observedPickupSpreadEndDate"]').clear().type('07 Oct 2022').blur(); // Observed pickup spread end date

      // Serious violations
      cy.get('[data-testid="radio"] [for="yes"]').click();
      cy.get('[data-testid="textarea"]').type('This is a test serious violation text');

      // Review report for submission
      openSubmissionPreview(true);

      // Veryify preview has correct sections/headers
      cy.contains('Move information').should('exist');
      cy.contains('Yes').should('exist');

      cy.get('h1').contains('Shipment report');
      cy.get('h2').contains('Evaluation report').parent().as('report');
      cy.get('h3').contains('Violations');
      cy.get('h3').contains('QAE remarks');

      // Verify preview has correct sections/headers/content
      cy.get('@report').within(() => {
        cy.contains('td', '01 Oct 2022');
        cy.contains('dd', 'Physical');
        cy.contains('dd', 'Other');
        cy.contains('dd', 'This is a test other location text');
        cy.contains('dd', 'This is a test evaluation report');

        //  kpi dates
        cy.contains('dd', '03 Oct 2022');
        cy.contains('dd', '04 Oct 2022');
        cy.contains('dd', '05 Oct 2022');
        cy.contains('dd', '06 Oct 2022');
        cy.contains('dd', '07 Oct 2022');
      });

      // Submit report
      submitReportFromPreview();

      verifySubmittedReportPresent();
    });

    it('can edit a saved draft', () => {
      // Create a new shipment report
      createShipmentReport();

      // Fill out the eval report form with minimal info
      fillInForm();

      // Save draft
      saveAsDraft();

      // Try to edit a draft
      cy.get('[data-testid="editReport"]').first().click();

      // Verify the form to edit is displayed
      cy.wait(['@getEvaluationReport']);
      cy.get('[data-testid="evaluationReportForm"]', { timeout: 10000 }).should('exist');
    });

    it('can delete a draft', () => {
      // Create a new shipment report
      createShipmentReport();

      // Save draft (reroutes back to reports list)
      saveAsDraft();

      // Delete draft
      cy.get('[data-testid="deleteReport"]').first().click();

      // Confirm delete modal
      cy.contains('Are you sure you want to delete this report?').should('exist');
      cy.contains('Yes, delete').click({ force: true });
      cy.wait(['@deleteEvaluationReport', '@getShipmentEvaluationReports', '@getCounselingEvaluationReports']);

      // Verify deletion
      cy.contains('Your draft report has been saved').should('not.exist'); // wait for alert to replace previous one
      cy.get('.usa-alert__text').contains('Your report has been deleted');
    });

    it('can view a submitted draft', () => {
      // Create and submit a new report
      createShipmentReport();
      fillInForm();
      openSubmissionPreview();
      submitReportFromPreview();

      // View saved report
      cy.get('[data-testid="viewReport"]').first().click();

      // Verify the 'view report' modal is shown
      cy.contains('Evaluation report').should('exist');
      cy.get('[data-testid="EvaluationReportPreview"]').should('be.visible');
    });

    it('has option to download a submitted draft', () => {
      // Create and submit a new report
      createShipmentReport();
      fillInForm();
      openSubmissionPreview();
      submitReportFromPreview();

      // Verify download button is displayed
      cy.get('[data-testid="downloadReport"]').should('be.visible');
    });
  });

  describe('Save Draft Behavior', () => {
    it('can save report as draft from the first form page (no violations)', () => {
      // Create a new shipment report
      createShipmentReport();

      // Fill out the eval report form with minimal info
      fillInForm();

      // Save draft
      saveAsDraft();

      // Verify draft saved
      cy.contains('Your draft report has been saved');
    });

    it('can save report as draft from the second form page (has violations)', () => {
      // Create a new shipment report
      createShipmentReport();

      // Fill out the eval report form with minimal info
      fillInForm();

      // Go to next page of form to select violations
      goToFormPageTwo();

      // Select a violation
      selectViolation('Storage', '1.2.6.16 Storage');

      // Select no serious violations
      cy.get('[data-testid="radio"] [for="no"]').click();

      // Save draft
      saveAsDraft();

      // Verify draft saved
      cy.contains('Your draft report has been saved');
    });
  });

  describe('Cancel Behavior', () => {
    // Create a report, click cancel before saving, verify/confim modal, verify delete worked upon reroute
    it('prompts to delete if the report has not been saved since creation', () => {
      // Create a new shipment report
      createShipmentReport();

      // Cancel report
      cy.get('[data-testid="cancelReport"]').click();

      // Verify cancel modal
      cy.contains('Are you sure you want to cancel this report?').should('exist');
      cy.contains('Yes, cancel').click({ force: true });

      // Verify deletion
      cy.get('.usa-alert__text').contains('Your report has been canceled');
    });

    // Create a report, save draft, click cancel, verify reroute w/out delete
    it('does not prompt to delete report after first save', () => {
      // Create a new shipment report
      createShipmentReport();
      // Fill out the eval report form with minimal info
      fillInForm();

      // Save draft
      saveAsDraft();

      // Edit draft report
      cy.get('[data-testid="editReport"]').last().click();

      // Verify the form to edit is displayed
      cy.wait(['@getEvaluationReport']);
      cy.get('[data-testid="evaluationReportForm"]', { timeout: 10000 }).should('exist');

      // Click cancel
      cy.get('[data-testid="cancelReport"]').click();

      // No cancel report since not deleting
      cy.contains('Are you sure you want to cancel this report?').should('not.exist');

      // Should have been rerouted to reports list
      cy.url().should('include', `/moves/${moveCode}/evaluation-reports`);
      cy.wait(['@getMoves', '@getMTOShipments']);
      cy.get('[data-testid="evaluationReportTable"]').should('exist');

      // Verify no delete modal
      cy.get('.usa-alert__text').should('not.exist');
    });

    // Create a report through 2nd page (violations),  click cancel, verify reroute w/out delete
    it('returns to reports list without deletion when canceled from the 2nd form page', () => {
      // Create a new shipment report
      createShipmentReport();

      // Fill out the eval report form with minimal info
      fillInForm();

      // Go to next page of form to select violations
      goToFormPageTwo();

      // Click cancel
      cy.get('[data-testid="cancelReport"]').click();

      // No cancel report since not deleting
      cy.contains('Are you sure you want to cancel this report?').should('not.exist');

      // Should have been rerouted to reports list
      cy.url().should('include', `/moves/${moveCode}/evaluation-reports`);
      cy.wait(['@getMoves', '@getMTOShipments']);
      cy.get('[data-testid="evaluationReportTable"]').should('exist');

      // Verify no delete modal
      cy.get('.usa-alert__text').should('not.exist');
    });
  });
});
