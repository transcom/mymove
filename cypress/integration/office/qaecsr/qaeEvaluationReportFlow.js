import { QAECSROfficeUserType } from '../../../support/constants';
import { searchForAndNavigateToMove } from './qaeCSRIntegrationUtils';

describe('Quality Evaluation Report Flows', () => {
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
    cy.intercept('**/ghc/v1/evaluation-reports/**').as('getEvaluationReport');
    cy.intercept('**/ghc/v1/moves/search').as('getSearchResults');

    const userId = '2419b1d6-097f-4dc4-8171-8f858967b4db';
    cy.apiSignInAsUser(userId, QAECSROfficeUserType);
  });

  /* This test is being temporarily skipped until flakiness issues can be resolved. */
  it.skip('is able to create and save a draft shipment evaluation report', () => {
    const moveLocator = 'TEST12';
    // Navigate to the move
    searchForAndNavigateToMove(moveLocator);
    // Go to quality assurance tab
    cy.contains('Quality assurance').click();
    // create report object to edit
    cy.get('[data-testid="shipmentEvaluationCreate"]').click();
    cy.wait(['@getEvaluationReport']);
    cy.get('[data-testid="textarea"]').type('this is a remark');
    cy.contains('Save draft').click();

    cy.url().should('include', `/moves/${moveLocator}/evaluation-reports`);
    cy.wait(['@getMTOShipments', '@getShipmentEvaluationReports', '@getCounselingEvaluationReports']);
    cy.contains('Your draft report has been saved');
  });

  /* This test is being temporarily skipped until flakiness issues can be resolved. */
  it.skip('does not prompt to delete report after first save', () => {
    const moveLocator = 'TEST12';
    // Navigate to the move
    searchForAndNavigateToMove(moveLocator);
    // go to the QAE report section
    cy.contains('Quality assurance').click();
    cy.wait(['@getMTOShipments', '@getShipmentEvaluationReports', '@getCounselingEvaluationReports']);
    cy.get('[data-testid="shipmentEvaluationCreate"]').click();
    cy.wait(['@getEvaluationReport']);
    cy.get('[data-testid="textarea"]').type('this is a remark');
    cy.contains('Save draft').click();
    cy.url().should('include', `/moves/${moveLocator}/evaluation-reports`);
    cy.wait(['@getMTOShipments', '@getShipmentEvaluationReports', '@getCounselingEvaluationReports']);
    cy.contains('Your draft report has been saved');
    // On this run through the cancel button should kick us back to the reports list view.
    cy.contains('Edit report').click();
    cy.wait(['@getEvaluationReport']);
    cy.get('[data-testid="cancelForUpdated"]').click();
    cy.url().should('include', `/moves/${moveLocator}/evaluation-reports`);
    cy.contains('Edit report');
  });

  /* This test is being temporarily skipped until flakiness issues can be resolved. */
  it.skip('does prompt to delete if the report has not been saved since creation', () => {
    const moveLocator = 'TEST12';
    // Navigate to the move
    searchForAndNavigateToMove(moveLocator);
    // go to the QAE report section
    cy.contains('Quality assurance').click();
    cy.wait(['@getMTOShipments', '@getShipmentEvaluationReports', '@getCounselingEvaluationReports']);
    cy.get('[data-testid="shipmentEvaluationCreate"]').click();
    cy.wait(['@getEvaluationReport']);
    cy.contains('Cancel').click();
    cy.contains('Yes, Cancel').click();
    cy.url().should('include', `/moves/${moveLocator}/evaluation-reports`);
    cy.contains('Your report has been canceled');
  });
});
