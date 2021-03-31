import { fileUploadTimeout } from '../../../support/constants';

export function customerFillsInProfileInformation(reloadAfterEveryPage) {
  // dod info
  // does not have welcome message throughout setup
  cy.wait('@createServiceMember');
  cy.get('span').contains('Welcome,').should('not.exist');

  // CONUS OR OCONUS
  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/conus-oconus/);
  });
  cy.get('[data-testid="radio"] label').contains('CONUS').click();
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/dod-info/);
  });

  // DOD INFO
  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');
  cy.get('select[name="affiliation"]').select('Army');
  cy.get('input[name="edipi"]').type('1234567890');
  cy.get('select[name="rank"]').select('E-9');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/name/);
  });
  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place

  // name
  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');
  cy.get('input[name="first_name"]').type('Jane');
  cy.get('input[name="last_name"]').type('Doe');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/contact-info/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  // contact info
  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');
  cy.get('input[name="telephone"]').type('6784567890');
  cy.get('[type="checkbox"]').not('[disabled]').check({ force: true }).should('be.checked');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/current-duty/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  // duty station
  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');
  cy.selectDutyStation('Fort Carson', 'current_station');

  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/current-address/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  // residential-address
  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');
  cy.get('input[name="current_residence.street_address_1"]').type('123 main');
  cy.get('input[name="current_residence.city"]').type('Anytown');
  cy.get('select[name="current_residence.state"]').select('CO');
  cy.get('input[name="current_residence.postal_code"]').clear().type('00001').blur();
  cy.get('#postal_code-error').should('exist');
  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');
  cy.get('input[name="postal_code"]').clear().type('80913');
  cy.get('span[data-testid="errorMessage"]').should('not.exist');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/backup-address/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  // backup address
  cy.get('button[data-testid="wizardNextButton"]').should('be.disabled');
  cy.get('input[name="backup_mailing_address.street_address_1"]').type('567 Another St');
  cy.get('input[name="backup_mailing_address.city"]').type('Anytown');
  cy.get('select[name="backup_mailing_address.state"]').select('CO');
  cy.get('input[name="backup_mailing_address.postal_code"]').type('80913');
  cy.get('span[data-testid="errorMessage"]').should('not.exist');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/backup-contact/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  // backup contact
  cy.get('input[name="name"]').type('Douglas Glass');
  cy.get('input[name="email"]').type('doug@glass.net');
  cy.nextPage();

  cy.get('[data-testid="customer-header"]').contains('Jane Doe');
  cy.get('p').contains("You're leaving Fort Carson");
  cy.get('[data-testid="stepContainer2"]').get('[data-testid="button"]').contains('Add orders').click();
}

export function customerFillsOutOrdersInformation() {
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
  cy.get('button[data-testid="wizardNextButton"]', { timeout: fileUploadTimeout }).should('not.be.disabled').click();

  cy.get('h1').contains('Tips for planning your shipments');
  cy.nextPage();

  cy.visit('/');
  cy.get('[data-testid="doc-list-container"]').contains('top-secret.png');
  cy.get('button').contains('Plan your shipments').click();
}
