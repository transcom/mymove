import { QAECSROfficeUserType } from '../../../support/constants';

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

    // This user has multiple roles, which is the kind of user we use to test in staging.
    // By using this type of user, we can catch bugs like the one fixed in PR 6706.
    const userId = 'b264abd6-52fc-4e42-9e0f-173f7d217bc5';
    cy.apiSignInAsUser(userId, QAECSROfficeUserType);
  });

  it('is able to create and save a draft shipment evaluation report', () => {
    const moveLocator = 'TEST12';
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    // Move Details page
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.url().should('include', `/moves/${moveLocator}/details`);
    // go to the QAE report section
    cy.contains('Quality assurance').click();
    cy.wait(['@getMTOShipments', '@getShipmentEvaluationReports', '@getCounselingEvaluationReports']);

    // create report object to edit
    cy.get('[data-testid="shipmentEvaluationCreate"]').click();
    cy.wait(['@getEvaluationReport']);
    cy.get('[data-testid="textarea"]').type('this is a remark');
    cy.contains('Save draft').click();

    cy.url().should('include', `/moves/${moveLocator}/evaluation-reports`);
    cy.wait(['@getMTOShipments', '@getShipmentEvaluationReports', '@getCounselingEvaluationReports']);
    cy.contains('Your draft report has been saved');
  });

  it('does not prompt to delete report after first save', () => {
    const moveLocator = 'TEST12';
    cy.wait(['@getSortedOrders']);
    cy.contains(moveLocator).click();
    // Move Details page
    cy.wait(['@getMoves', '@getOrders', '@getMTOShipments', '@getMTOServiceItems']);
    cy.url().should('include', `/moves/${moveLocator}/details`);
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
    cy.contains('View report').click();
    cy.wait(['@getEvaluationReport']);
    cy.get('[data-testid="cancelForUpdated"]').click();
    cy.url().should('include', `/moves/${moveLocator}/evaluation-reports`);
    cy.contains('View report');
  });
});
