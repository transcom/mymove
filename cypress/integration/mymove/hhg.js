/* global cy */
import { fileUploadTimeout } from '../../support/constants';

function customerFillsInProfileInformation(reloadAfterEveryPage) {
  // dod info
  // does not have welcome message throughout setup
  cy.get('span').contains('Welcome,').should('not.exist');
  cy.nextPage();

  cy.get('button.next').should('be.disabled');
  cy.get('select[name="affiliation"]').select('Army');
  cy.get('input[name="edipi"]').type('1234567890');
  cy.get('input[name="social_security_number').type('123456789');
  cy.get('select[name="rank"]').select('E-9');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/name/);
  });
  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  // name
  cy.get('button.next').should('be.disabled');
  cy.get('input[name="first_name"]').type('Jane');
  cy.get('input[name="last_name"]').type('Doe');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/contact-info/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  // contact info
  cy.get('button.next').should('be.disabled');
  cy.get('input[name="telephone"]').type('6784567890');
  cy.get('[type="checkbox"]').not('[disabled]').check({ force: true }).should('be.checked');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/duty-station/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  // duty station
  cy.get('button.next').should('be.disabled');
  cy.selectDutyStation('Fort Carson', 'current_station');

  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/residence-address/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  // residential-address
  cy.get('button.next').should('be.disabled');
  cy.get('input[name="street_address_1"]').type('123 main');
  cy.get('input[name="city"]').type('Anytown');
  cy.get('select[name="state"]').select('CO');
  cy.get('input[name="postal_code"]').clear().type('00001').blur();
  cy.get('#postal_code-error').should('exist');
  cy.get('button.next').should('be.disabled');
  cy.get('input[name="postal_code"]').clear().type('80913');
  cy.get('#postal_code-error').should('not.exist');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/backup-mailing-address/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  // backup address
  cy.get('button.next').should('be.disabled');
  cy.get('input[name="street_address_1"]').type('567 Another St');
  cy.get('input[name="city"]').type('Anytown');
  cy.get('select[name="state"]').select('CO');
  cy.get('input[name="postal_code"]').type('80913');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/backup-contacts/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  // backup contact
  cy.get('input[name="name"]').type('Douglas Glass');
  cy.get('input[name="email"]').type('doug@glass.net');
  cy.nextPage();

  cy.get('h2').contains('Welcome Jane');
  cy.nextPage();
}

function customerFillsOutOrdersInformation() {
  cy.get('select[name="orders_type"]').select('Separation');
  cy.get('select[name="orders_type"]').select('Retirement');
  cy.get('select[name="orders_type"]').select('Permanent Change Of Station (PCS)');

  cy.get('input[name="issue_date"]').first().click();

  cy.get('input[name="issue_date"]').first().type('6/2/2018{enter}').blur();

  cy.get('input[name="report_by_date"]').last().type('8/9/2018{enter}').blur();

  cy.selectDutyStation('NAS Fort Worth JRB', 'new_duty_station');

  cy.nextPage();

  cy.location().should((loc) => {
    expect(loc.pathname).to.eq('/orders/upload');
  });

  cy.upload_file('.filepond--root', 'top-secret.png');
  cy.get('button.next', { timeout: fileUploadTimeout }).should('not.be.disabled').click();
  cy.nextPage();

  cy.get('h1').contains('Moving 101');
  cy.nextPage();
}

function customerSetsUpAnHHGMove() {
  cy.get('input[type="radio"]').last().check({ force: true });
  cy.nextPage();
  cy.get('input[name="requestedPickupDate"]').first().type('08/02/2020').blur();

  // should be empty before using "Use current residence" checkbox
  cy.get(`[data-testid="mailingAddress1"]`).first().should('be.empty');
  cy.get(`[data-testid="city"]`).first().should('be.empty');
  cy.get(`[data-testid="state"]`).first().should('be.empty');
  cy.get(`[data-testid="zip"]`).first().should('be.empty');

  cy.get(`input[name="useCurrentResidence"]`).check({ force: true });

  // releasing agent
  cy.get(`[data-testid="firstName"]`).first().type('John');
  cy.get(`[data-testid="lastName"]`).first().type('Lee');
  cy.get(`[data-testid="phone"]`).first().type('999-999-9999');
  cy.get(`[data-testid="email"]`).first().type('ron@example.com');

  // requested delivery date
  cy.get('input[name="requestedDeliveryDate"]').first().type('09/20/2020').blur();
  // checks has delivery address (default does not have delivery address)
  cy.get('input[type="radio"]').first().check({ force: true });

  // delivery location
  cy.get(`[data-testid="mailingAddress1"]`).last().type('412 Avenue M ');
  cy.get(`[data-testid="mailingAddress2"]`).last().type('#3E');
  cy.get(`[data-testid="city"]`).last().type('Los Angeles');
  cy.get(`[data-testid="state"]`).last().type('CA');
  cy.get(`[data-testid="zip"]`).last().type('90011');

  // releasing agent
  cy.get(`[data-testid="firstName"]`).last().type('John');
  cy.get(`[data-testid="lastName"]`).last().type('Lee');
  cy.get(`[data-testid="phone"]`).last().type('999-999-9999');
  cy.get(`[data-testid="email"]`).last().type('ron@example.com');

  // customer remarks
  cy.get(`[data-testid="remarks"]`).first().type('some customer remark');
  cy.nextPage();

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/moves\/[^/]+\/review/);
  });
}

describe('HHG Setup flow', function () {
  it('Creates a shipment', function () {
    cy.signInAsNewMilMoveUser();
    customerFillsInProfileInformation();
    customerFillsOutOrdersInformation();
    customerSetsUpAnHHGMove();
  });
});
