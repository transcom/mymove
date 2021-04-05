import { adminBaseURL } from '../../support/constants';

describe('Office Users List Page', function () {
  before(() => {
    cy.prepareAdminApp();
    cy.logout();
  });

  it('successfully navigates to Office Users page', function () {
    cy.signInAsNewAdminUser();
    // we should be at the office users list page by default,
    // but let's go somewhere else and then come back to make sure the side nav link works
    cy.get('a[href*="system/moves"]').click();
    cy.url().should('eq', adminBaseURL + '/system/moves');

    // now we'll come back to the office users page:
    cy.get('a[href*="system/office_users"]').click();
    cy.url().should('eq', adminBaseURL + '/system/office_users');
    cy.get('header').contains('Office users');

    const columnLabels = ['Id', 'Email', 'First name', 'Last name', 'Transportation Office', 'User Id', 'Active'];
    columnLabels.forEach((label) => {
      cy.get('table').contains(label);
    });
  });
});

describe('Office User Create Page', function () {
  before(() => {
    cy.prepareAdminApp();
    cy.logout();
  });

  it('pulls up create page for an office user', function () {
    cy.signInAsNewAdminUser();
    // we tested the side nav in the previous test,
    // so let's work with the assumption that we were already redirected to this page:
    cy.url().should('eq', adminBaseURL + '/system/office_users');
    cy.get('a[href*="system/office_users/create"]').first().click();
    cy.url().should('eq', adminBaseURL + '/system/office_users/create');

    // we need to add the date to the email so that it is unique every time (only one record per email allowed in db)
    const testEmail = 'cy.office_user.' + Date.now() + '@example.com';

    // create an admin user
    cy.get('input[id="firstName"]').type('Cypress');
    cy.get('input[id="middleInitials"]').type('CY');
    cy.get('input[id="lastName"]').type('Test');
    cy.get('input[id="email"]').type(testEmail);
    cy.get('input[id="telephone"]').type('222-555-1234');
    cy.get('.ra-input-roles input[type="checkbox"]').first().click();
    cy.get('input[id="transportationOfficeId"]').type('JPPSO Testy McTest');
    cy.wait(5000); // we have to wait for the autocomplete to give us something
    cy.get('div[role="tooltip"] div[role="listbox"] li')
      .first()
      .then(($opt) => {
        $opt.trigger('click');
      });
    cy.get('button').contains('Save').click();

    // redirected to edit details page
    cy.get('#id')
      .invoke('val')
      .then((officeUserID) => {
        cy.url().should('contain', officeUserID);
      });

    cy.get('#email').should('have.value', testEmail);
    cy.get('#firstName').should('have.value', 'Cypress');
    cy.get('#middleInitials').should('have.value', 'CY');
    cy.get('#lastName').should('have.value', 'Test');
    cy.get('#telephone').should('have.value', '222-555-1234');
    cy.get('#active').should('contain', 'Yes');
  });
});

describe('Office Users Show Page', function () {
  before(() => {
    cy.prepareAdminApp();
    cy.logout();
  });

  it('pulls up details page for an office user', function () {
    cy.signInAsNewAdminUser();
    cy.url().should('eq', adminBaseURL + '/system/office_users');
    cy.get('tr[resource="office_users"]').first().click();

    // check that the office user's name is shown in the page title
    cy.get('.ra-field-firstName span.MuiTypography-root')
      .invoke('text')
      .then((firstName) => {
        cy.get('.ra-field-lastName span.MuiTypography-root')
          .invoke('text')
          .then((lastName) => {
            cy.get('#react-admin-title').should('contain', firstName + ' ' + lastName);
          });
      });

    const labels = [
      'Id',
      'User Id',
      'Email',
      'First name',
      'Middle initials',
      'Last name',
      'Telephone',
      'Active',
      'Roles',
      'Transportation Office',
      'Created at',
      'Updated at',
    ];
    labels.forEach((label) => {
      cy.get('.MuiCardContent-root label').contains(label);
    });
  });
});

describe('Office Users Edit Page', function () {
  before(() => {
    cy.prepareAdminApp();
    cy.logout();
  });

  it('pulls up edit page for an office user', function () {
    cy.signInAsNewAdminUser();
    cy.url().should('eq', adminBaseURL + '/system/office_users');
    cy.get('tr[resource="office_users"]').first().click();

    // grab the office user's ID to check that the correct value is in the url
    cy.get('.ra-field-id span.MuiTypography-root')
      .invoke('text')
      .as('officeUserID')
      .then((officeUserID) => {
        // continue to the edit page
        cy.get('a').contains('Edit').click();
        cy.url().should('eq', adminBaseURL + '/system/office_users/' + officeUserID);
      });

    const disabledFields = ['id', 'email', 'userId', 'createdAt', 'updatedAt'];
    disabledFields.forEach((label) => {
      cy.get('[id="' + label + '"]').should('be.disabled');
    });

    cy.get('input[id="firstName"]').clear().type('Edit');
    cy.get('input[id="lastName"]').clear().type('CypressUser');

    // set the user to the active status they did NOT have before
    cy.get('#active').click();
    cy.get('ul[aria-labelledby="active-label"] li')
      .not('[aria-selected="true"]')
      .click()
      .then(($selectedOpt) => {
        // grab the value we selected for "Active"
        const newActiveValue = $selectedOpt.attr('data-value');

        cy.get('button').contains('Save').click();

        // back in the move list screen, check that the row for this office user was updated
        cy.url().should('eq', adminBaseURL + '/system/office_users');

        // check that the correct values were saved:
        cy.get('tr')
          .contains(this.officeUserID)
          .parents('tr')
          .find('td.column-active svg') // this column uses icons to indicate true or false
          .should(($activeCol) => {
            expect($activeCol.attr('data-testid')).to.eq(newActiveValue);
          });

        cy.get('tr')
          .contains(this.officeUserID)
          .parents('tr')
          .find('td.column-firstName span.MuiTypography-root')
          .invoke('text')
          .should('eq', 'Edit');

        cy.get('tr')
          .contains(this.officeUserID)
          .parents('tr')
          .find('td.column-lastName span.MuiTypography-root')
          .invoke('text')
          .should('eq', 'CypressUser');
      });
  });
});
