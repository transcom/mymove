import { adminBaseURL } from '../../support/constants';

describe('Webhook Subscriptions', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('successfully navigates to the webhook subscriptions list page', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/webhook_subscriptions"]').click();
    cy.url().should('eq', adminBaseURL + '/system/webhook_subscriptions');
    cy.get('header').contains('Webhook subscriptions');

    const columnLabels = ['Id', 'Event key', 'Severity', 'Callback url', 'Status', 'Updated at'];
    columnLabels.forEach((label) => {
      cy.get('.RaDatagrid-table-65').contains(label);
    });
  });
});
