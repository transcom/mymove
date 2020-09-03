describe('setting up service member profile requiring an access code', function () {
  before(() => {
    cy.prepareCustomerApp();
  });

  beforeEach(() => {
    cy.signInAsNewMilMoveUser();
  });

  it('progresses thru forms', function () {
    cy.get('body').then(($body) => {
      if ($body.find('input[name="claim_access_code"]').length) {
        serviceMemberEntersAccessCode();
      }
    });
    serviceMemberChoosesConusOrOconus();
    serviceMemberProfile();
  });

  it.skip('restarts app after every page', function () {
    serviceMemberProfile(true);
  });
});

function serviceMemberEntersAccessCode() {
  cy.get('input[name="claim_access_code"]').type('PPM-X3FQJK');
  cy.get('button').contains('Continue').click();
}

function serviceMemberChoosesConusOrOconus() {
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/conus-status/);
  });
  cy.get('[data-testid="radio"] label').contains('CONUS');
  cy.get('button.next').click();
}

function serviceMemberProfile(reloadAfterEveryPage) {
  //dod info
  // does not have welcome message throughout setup
  cy.get('span').contains('Welcome,').should('not.exist');

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
  //name
  cy.get('button.next').should('be.disabled');
  cy.get('input[name="first_name"]').type('Jane');
  cy.get('input[name="last_name"]').type('Doe');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/contact-info/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  //contact info
  cy.get('button.next').should('be.disabled');
  cy.get('input[name="telephone"]').type('6784567890');
  cy.get('[type="checkbox"]').not('[disabled]').check({ force: true }).should('be.checked');
  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/duty-station/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  //duty station
  cy.get('button.next').should('be.disabled');
  cy.selectDutyStation('Fort Carson', 'current_station');

  cy.nextPage();
  cy.location().should((loc) => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/residence-address/);
  });

  if (reloadAfterEveryPage) cy.visit('/'); // make sure picks up in right place
  //residential-address
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
  //backup contact
  cy.get('input[name="name"]').type('Douglas Glass');
  cy.get('input[name="email"]').type('doug@glass.net');
  cy.nextPage();

  //transition
  // cy.nextPage();
}
