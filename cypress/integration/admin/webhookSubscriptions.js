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

describe('Webhook Subscription Create Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('pulls up create page for a webhook subscription', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/webhook_subscriptions"]').click();
    cy.url().should('eq', adminBaseURL + '/system/webhook_subscriptions');
    cy.get('a[href*="system/webhook_subscriptions/create"]').first().click();
    cy.url().should('eq', adminBaseURL + '/system/webhook_subscriptions/create');

    cy.get('input[id="subscriberId"]').type('5db13bb4-6d29-4bdb-bc81-262f4513ecf6');
    cy.get('input[id="eventKey"]').type('PaymentRequest.Update');
    cy.get('input[id="callbackUrl"]').type('https://test1.com');
    cy.get('div[id="status"]').click();
    cy.get('li[data-value="ACTIVE"]').click();
    cy.get('button').contains('Save').click();

    // redirected to details page
    cy.get('.ra-field-id span.MuiTypography-root')
      .invoke('text')
      .then(() => {
        cy.get('#react-admin-title').contains('Webhook Subscription ID: ');
      });

    cy.get('.ra-field-subscriberId').contains('5db13bb4-6d29-4bdb-bc81-262f4513ecf6');
    cy.get('.ra-field-eventKey').contains('PaymentRequest.Update');
    cy.get('.ra-field-callbackUrl').contains('https://test1.com');
    cy.get('.ra-field-status').contains('ACTIVE');
    cy.get('.ra-field-severity').contains('0');
  });
});
