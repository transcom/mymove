import {
  navigateFromCloseoutReviewPageToExpensesPage,
  setMobileViewport,
  signInAndNavigateToPPMReviewPage,
  signInAndNavigateToWeightTicketPage,
} from '../../../support/ppmCustomerShared';

describe('Expenses', function () {
  before(() => {
    cy.prepareCustomerApp();
  });
  beforeEach(() => {
    cy.intercept('GET', '**/internal/moves/**/mto_shipments').as('getShipment');
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

      cy.get('select[name="expenseType"]').as('expenseType').should('have.value', '');
      cy.get('@expenseType').select('Storage');

      cy.get('input[name="description"]').type('Cloud storage');
      cy.get('input[name="paidWithGTCC"][value="true"]').click({ force: true });
      cy.get('input[name="amount"]').type('675.99');

      // TODO: Add receipt document upload when integrated

      cy.get('input[name="sitStartDate"]').type('14 Aug 2022').blur();
      cy.get('input[name="sitEndDate"]').type('20 Aug 2022').blur();

      cy.a11y();
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

      cy.a11y();
    });
  });
});
