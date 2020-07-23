/* global cy */
import { milmoveAppName, fileUploadTimeout } from '../../support/constants';

function customerFillsInProfileInformation(reloadAfterEveryPage) {
  // dod info
  // does not have welcome message throughout setup
  cy.get('span').contains('Welcome,').should('not.exist');

  // does not have a back button on first flow page
  cy.get('button').contains('Back').should('not.be.visible');

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

  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/transition/);
  });

  // transition
  cy.nextPage();
}

function customerFillsOutOrdersInformation() {
  cy.location().should((loc) => {
    expect(loc.pathname).to.eq('/orders/');
  });

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
}

function customerSetsUpAnHHGMove() {
  cy.get('input[type="radio"]').last().check({ force: true });
  cy.nextPage();
}

describe('The Home Page', function () {
  beforeEach(() => {
    cy.setupBaseUrl(milmoveAppName);
  });
  it('Goes through ', function () {
    cy.signInAsNewMilMoveUser();
    customerFillsInProfileInformation();
    customerFillsOutOrdersInformation();
    customerSetsUpAnHHGMove();
  });
});
