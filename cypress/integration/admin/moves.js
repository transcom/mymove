import { adminBaseURL } from '../../support/constants';

describe('Moves Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('successfully navigates to Moves page', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/moves"]').click();
    cy.url().should('eq', adminBaseURL + '/system/moves');
    cy.get('header').contains('Moves');

    const columnLabels = [
      'Id',
      'Order Id',
      'Service Member Id',
      'Locator',
      'Status',
      'Show',
      'Created at',
      'Updated at',
    ];
    columnLabels.forEach((label) => {
      cy.get('table').contains(label);
    });
  });
});

describe('Moves Details Show Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('pulls up details page for a move', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/moves"]').click();
    cy.url().should('eq', adminBaseURL + '/system/moves');
    cy.get('span[reference="moves"]').first().click();

    // check that the move's ID is shown in the page title
    cy.get('.ra-field-id span.MuiTypography-root')
      .invoke('text')
      .then((moveID) => {
        cy.get('#react-admin-title').contains('Move ID: ' + moveID);
      });

    const labels = [
      'Id',
      'Locator',
      'Status',
      'Show',
      'Order Id',
      'Created at',
      'Updated at',
      'User Id',
      'Service member Id',
      'Service member first name',
      'Service member middle name',
      'Service member last name',
    ];
    labels.forEach((label) => {
      cy.get('.MuiCardContent-root label').contains(label);
    });
  });
});

describe('Moves Details Edit Page', function () {
  before(() => {
    cy.prepareAdminApp();
  });

  it('pulls up edit page for a move', function () {
    cy.signInAsNewAdminUser();
    cy.get('a[href*="system/moves"]').click();
    cy.url().should('eq', adminBaseURL + '/system/moves');
    cy.get('span[reference="moves"]').first().click();

    // grab the move's ID to check that the correct value is in the url
    cy.get('.ra-field-id span.MuiTypography-root')
      .invoke('text')
      .as('moveID')
      .then((moveID) => {
        // continue to the edit page
        cy.get('a').contains('Edit').click();
        cy.url().should('eq', adminBaseURL + '/system/moves/' + moveID);
      });

    const disabledFields = [
      'id',
      'locator',
      'status',
      'ordersId',
      'createdAt',
      'updatedAt',
      'serviceMember.userId',
      'serviceMember.id',
      'serviceMember.firstName',
      'serviceMember.middleName',
      'serviceMember.lastName',
    ];
    disabledFields.forEach((label) => {
      cy.get('[id="' + label + '"]').should('be.disabled');
    });

    cy.get('#show').click();
    cy.get('ul[aria-labelledby="show-label"] li')
      .not('[aria-selected="true"]')
      .click()
      .then(($selectedOpt) => {
        // grab the value we selected for "Show"
        const newShowValue = $selectedOpt.attr('data-value');

        cy.get('button').contains('Save').click();

        // back in the move list screen, check that the row for this move was updated
        cy.url().should('eq', adminBaseURL + '/system/moves');
        cy.get('tr')
          .contains(this.moveID)
          .parents('tr')
          .find('td.column-show span.MuiTypography-root')
          .should(($showCol) => {
            expect($showCol.text()).to.eq(newShowValue);
          });
      });
  });
});
