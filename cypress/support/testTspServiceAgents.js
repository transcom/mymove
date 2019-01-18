/* global cy, Cypress*/

export function getFixture(role) {
  return {
    Origin: {
      Company: 'ACME Movers',
      Email: 'acme@example.com',
      Phone: '303-867-5309',
      Role: 'origin_service_agent',
    },
    OriginUpdate: {
      Company: 'ACME Movers',
      Email: 'acmemovers@example.com',
      Phone: '303-867-5308',
      Role: 'origin_service_agent',
    },
    Destination: {
      Company: 'ACE Movers',
      Email: 'acmemoving@example.com',
      Phone: '303-867-5310',
      Role: 'destination_service_agent',
    },
    DestinationUpdate: {
      Company: 'ACE Moving Company',
      Email: 'moveme@example.com',
      Phone: '303-867-5311',
      Role: 'destination_service_agent',
    },
  }[role];
}

export function userVerifiesTspAssigned() {
  const tspFields = cy
    .get('.editable-panel-3-column')
    .contains('TSP')
    .parent()
    .within(() => {
      cy
        .get('.panel-field')
        .contains('Name')
        .parent()
        .should('contain', 'Truss Transport LLC (J12K)');
      cy
        .get('.panel-field')
        .contains('Email')
        .parent()
        .should('contain', 'joey.j@example.com');
      cy
        .get('.panel-field')
        .contains('Phone number')
        .parent()
        .should('contain', '(555) 101-0101');
    });
}
export function userInputsServiceAgent(role) {
  const fixture = getFixture(role);

  // Enter details in form
  cy
    .get('input[name="' + fixture.Role + '.company"]')
    .first()
    .type(fixture.Company)
    .blur();
  cy
    .get('input[name="' + fixture.Role + '.email"]')
    .first()
    .type(fixture.Email)
    .blur();
  cy
    .get('input[name="' + fixture.Role + '.phone_number"]')
    .first()
    .type(fixture.Phone)
    .blur();
}

export function userClearsServiceAgent(role) {
  const fixture = getFixture(role);

  // Clear details in form
  cy
    .get('input[name="' + fixture.Role + '.company"]')
    .clear()
    .blur();
  cy
    .get('input[name="' + fixture.Role + '.email"]')
    .clear()
    .blur();
  cy
    .get('input[name="' + fixture.Role + '.phone_number"]')
    .clear()
    .blur();
}

export function userCancelsServiceAgent(role) {
  cy
    .get('button')
    .contains('Cancel')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Cancel')
    .click();
}

export function userSavesServiceAgent(role) {
  const fixture = getFixture(role);

  cy
    .get('button')
    .contains('Save')
    .should('be.enabled');

  cy
    .get('button')
    .contains('Save')
    .click();

  // Verify data has been saved in the UI
  cy
    .get('div.company')
    .get('span')
    .contains(fixture.Company);
  cy
    .get('div.email')
    .get('span')
    .contains(fixture.Email);
  cy
    .get('div.phone_number')
    .get('span')
    .contains(fixture.Phone);

  // Refresh browser and make sure changes persist
  cy.patientReload();

  cy
    .get('div.company')
    .get('span')
    .contains(fixture.Company);
  cy
    .get('div.email')
    .get('span')
    .contains(fixture.Email);
  cy
    .get('div.phone_number')
    .get('span')
    .contains(fixture.Phone);
}
