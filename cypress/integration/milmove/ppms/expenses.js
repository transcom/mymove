import {
  navigateFromCloseoutReviewPageToExpensesPage,
  setMobileViewport,
  signInAndNavigateToPPMReviewPage,
  submitExpensePage,
} from '../../../support/ppmCustomerShared';

describe('Expenses', function () {
  before(() => {
    cy.prepareCustomerApp();
  });
  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
    cy.intercept('POST', '**/internal/ppm-shipments/**/uploads**').as('uploadFile');
  });
  const viewportType = [
    { viewport: 'desktop', isMobile: false, userId: '146c2665-5b8a-4653-8434-9a4460de30b5' }, // movingExpensePPM@ppm.approved
    { viewport: 'mobile', isMobile: true, userId: '7d4dbc69-2973-4c8b-bf75-6fb582d7a5f6' }, // movingExpensePPM2@ppm.approved
  ];
  viewportType.forEach(({ viewport, isMobile, userId }) => {
    it(`new expense page loads - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToPPMReviewPage(userId);
      navigateFromCloseoutReviewPageToExpensesPage();
      submitExpensePage();
    });

    it(`edit expense page loads - ${viewport}`, () => {
      if (isMobile) {
        setMobileViewport();
      }

      signInAndNavigateToPPMReviewPage(userId);

      // edit the first expense receipt
      cy.get('.reviewExpenses a').contains('Edit').eq(0).click();
      cy.location().should((loc) => {
        expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/expenses\/[^/]+/);
      });

      cy.get('select[name="expenseType"]').as('expenseType').should('have.value', 'PACKING_MATERIALS');
      cy.get('@expenseType').select('Tolls');

      cy.get('input[name="description"]').as('expenseDescription').should('have.value', 'Packing Peanuts');
      cy.get('@expenseDescription').clear().type('PA Turnpike EZ-Pass');

      cy.get('input[name="paidWithGTCC"][value="true"]').should('be.checked');
      cy.get('input[name="paidWithGTCC"][value="false"]').click({ force: true });

      cy.get('input[name="amount"]').as('expenseAmount').should('have.value', '23.45');
      cy.get('@expenseAmount').clear().type('54.32');

      cy.get('input[name="missingReceipt"]').as('missingReceipt').should('not.be.checked');
      cy.get('@missingReceipt').click({ force: true });

      cy.get('button').contains('Save & Continue').should('be.enabled').click();
      cy.location().should((loc) => {
        expect(loc.pathname).to.match(/^\/moves\/[^/]+\/shipments\/[^/]+\/review/);
      });

      cy.contains('PA Turnpike EZ-Pass');
    });
  });
});
