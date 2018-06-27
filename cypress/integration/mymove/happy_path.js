/* global cy */
function signInAsNewUser() {
  cy.contains('Local Sign In').click();
  cy.contains('Login as New User').click();
}
function next() {
  cy
    .get('button.next')
    .should('be.enabled')
    .click();
}

function serviceMemberProfile(reload) {
  //dod info
  cy.get('button.next').should('be.disabled');
  cy.get('select[name="affiliation"]').select('Army');
  cy.get('input[name="edipi"]').type('1234567890');
  cy.get('input[name="social_security_number').type('123456789');
  cy.get('select[name="rank"]').select('E-9');
  next();

  if (reload) cy.visit('/'); // make sure picks up in right place
  //name
  cy.get('button.next').should('be.disabled');
  cy.get('input[name="first_name"]').type('Jane');
  cy.get('input[name="last_name"]').type('Doe');
  next();

  if (reload) cy.visit('/'); // make sure picks up in right place
  //contact info
  cy.get('button.next').should('be.disabled');
  cy.get('input[name="telephone"]').type('6784567890');
  cy
    .get('[type="checkbox"]')
    .not('[disabled]')
    .check({ force: true })
    .should('be.checked');
  next();

  if (reload) cy.visit('/'); // make sure picks up in right place
  //duty station
  cy.get('button.next').should('be.disabled');
  cy
    .get('.duty-input-box__input input')
    .first()
    .type('Ft Carson{downarrow}{enter}', { force: true, delay: 150 });

  next();

  if (reload) cy.visit('/'); // make sure picks up in right place
  //residential-address
  cy.get('button.next').should('be.disabled');
  cy.get('input[name="street_address_1"]').type('123 main');
  cy.get('input[name="city"]').type('Anytown');
  cy.get('select[name="state"]').select('CO');
  cy.get('input[name="postal_code"]').type('80913');
  next();

  if (reload) cy.visit('/'); // make sure picks up in right place
  // backup address
  cy.get('button.next').should('be.disabled');
  cy.get('input[name="street_address_1"]').type('567 Another St');
  cy.get('input[name="city"]').type('Anytown');
  cy.get('select[name="state"]').select('CO');
  cy.get('input[name="postal_code"]').type('80913');
  next();

  if (reload) cy.visit('/'); // make sure picks up in right place
  //backup contact
  cy.get('input[name="name"]').type('Douglas Glass');
  cy.get('input[name="email"]').type('doug@glass.net');
  next();

  //transition
  // next();
}
describe('The Home Page', function() {
  it('successfully loads', function() {
    cy.visit('/');
  });
});
describe('setting up service member profile', function() {
  beforeEach(() => {
    signInAsNewUser();
  });
  afterEach(() => {
    cy.contains('Sign Out').click();
  });
  it('progresses thru forms', function() {
    serviceMemberProfile();
  });
  it('starts finishing the service member profile', function() {
    serviceMemberProfile(true);
  });
});
