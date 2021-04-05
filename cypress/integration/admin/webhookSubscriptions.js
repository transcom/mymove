import { adminBaseURL } from '../../support/constants';

const checkLabels = (itemToGet, labels) => {
  labels.forEach((label) => {
    cy.get(itemToGet).contains(label);
  });
};

describe('Webhook Subscriptions', function () {
  before(() => {
    cy.prepareAdminApp();
    cy.logout();
  });

  it('successfully navigates to the webhook subscriptions list page', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/webhook_subscriptions"]').click();
    cy.url().should('eq', adminBaseURL + '/system/webhook_subscriptions');
    cy.get('header').contains('Webhook subscriptions');

    const columnLabels = ['Id', 'Event key', 'Severity', 'Callback url', 'Status', 'Updated at'];
    checkLabels('table', columnLabels);
  });
});

describe('WebhookSubscriptions Details Show Page', function () {
  before(() => {
    cy.prepareAdminApp();
    cy.logout();
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
    checkLabels('.MuiCardContent-root label', labels);
  });
});

describe('WebhookSubscriptions Details Edit Page', function () {
  before(() => {
    cy.prepareAdminApp();
    cy.logout();
  });

  it('pulls up edit page for a webhook subscription', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/webhook_subscriptions"]').click();
    cy.url().should('eq', adminBaseURL + '/system/webhook_subscriptions');
    cy.get('tr[resource="webhook_subscriptions"]').first().click();

    // check that the page pulls up the right webhook subscription
    cy.get('.ra-field-id > div > label')
      .first()
      .next()
      .then(($id) => {
        cy.get('a').contains('Edit').click();
        cy.url().should('eq', adminBaseURL + '/system/webhook_subscriptions/' + $id.text());
      });

    // check labels on edit page
    const labels = ['Id', 'Subscriber', 'Status', 'Event key', 'Callback url', 'Created at', 'Updated at', 'Severity'];
    checkLabels('label', labels);

    // Change webhook subscription status
    cy.get('div[id="status"]').click();
    cy.get('#menu-status ul > li[data-value=DISABLED').click();
    cy.get('button').contains('Save').click();

    // Check that the webhook subscription status was changed for the first webhook subscription in the list
    cy.url().should('eq', adminBaseURL + '/system/webhook_subscriptions');
    cy.get('td.column-status > span').first().should('contain', 'DISABLED');
  });
});

describe('Webhook Subscription Create Page', function () {
  before(() => {
    cy.prepareAdminApp();
    cy.logout();
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
    cy.get('input[id="id"]')
      .invoke('text')
      .then((subID) => {
        cy.get('#react-admin-title').contains('Webhook subscription #' + subID);
      });

    cy.get('input[id="subscriberId"]').type('5db13bb4-6d29-4bdb-bc81-262f4513ecf6');
    cy.get('input[id="eventKey"]').type('PaymentRequest.Update');
    cy.get('input[id="callbackUrl"]').type('https://test1.com');
    cy.get('div[id="status"]').should('contain', 'Active');
    cy.get('div[id="severity"]').should('contain', '0');
  });
});
