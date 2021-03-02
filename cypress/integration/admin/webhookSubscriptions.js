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
      cy.get('table').contains(label);
    });
  });
});

describe('WebhookSubscriptions Details Show Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('pulls up details page for a webhook subscription', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/webhook_subscriptions"]').click();
    cy.url().should('eq', adminBaseURL + '/system/webhook_subscriptions');
    cy.get('tr[resource="webhook_subscriptions"]').first().click();

    // check that the webhookSubscription's ID is shown in the page title
    cy.get('.ra-field-id span.MuiTypography-root')
      .invoke('text')
      .then((webhookSubscriptionID) => {
        cy.get('#react-admin-title').contains('Webhook Subscription ID: ' + webhookSubscriptionID);
      });

    const labels = ['Id', 'Subscriber', 'Status', 'Event key', 'Callback url', 'Created at', 'Updated at', 'Severity'];
    labels.forEach((label) => {
      cy.get('.MuiCardContent-root label').contains(label);
    });
  });
});
