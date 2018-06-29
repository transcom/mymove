/* global cy, Cypress*/
describe('orders entry', function() {
  beforeEach(() => {
    cy.signInAsNewUser();
  });

  it('will accept orders information', function() {
    createServiceMember().then(() => cy.visit('/'));
    cy.contains('New move from Ft Carson');
    cy.contains('No detail');
    cy.contains('No documents');
    cy.contains('Continue Move Setup').click();

    cy.location().should(loc => {
      expect(loc.pathname).to.eq('/orders/');
    });

    cy.get('select[name="orders_type"]').select('Permanent Change Of Station');
    cy
      .get('input[placeholder="Date"]')
      .first()
      .click();

    cy
      .get('input[placeholder="Date"]')
      .first()
      .type('2018-6-2{enter}')
      .blur();

    cy
      .get('input[placeholder="Date"]')
      .last()
      .type('2018-8-9{enter}')
      .blur();

    cy
      .get('.duty-input-box__input input')
      .first()
      .type('Fort Worth{downarrow}{enter}', { force: true, delay: 150 });

    cy.nextPage();

    cy.location().should(loc => {
      expect(loc.pathname).to.eq('/orders/upload');
    });

    cy.visit('/');
    cy.contains('NAS Fort Worth from Ft Carson');
    cy.get('.whole_box > :nth-child(3) > span').contains('7,000 lbs');
    cy.contains('Continue Move Setup').click();
    cy.location().should(loc => {
      expect(loc.pathname).to.eq('/orders/upload');
    });
  });
});

function createServiceMember() {
  return cy
    .location()
    .should(loc => {
      expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/create/);
    })
    .then(location => {
      const serviceMemberId = location.pathname.match(
        /^\/service-member\/([^/]+)\/create/,
      )[1];
      return serviceMemberId;
    })
    .then(serviceMemberId =>
      cy.fixture('serviceMember.json').then(serviceMember =>
        cy
          .window()
          .then(window => {
            serviceMember.user_id = window._state.user.userId;
            return cy
              .request('internal/duty_stations?search=ft%20car')
              .then(result => {
                serviceMember.current_station_id = Cypress._.get(
                  result,
                  'body.[0].id',
                );
                return cy.request(
                  'PATCH',
                  `/internal/service_members/${serviceMemberId}`,
                  serviceMember,
                );
              });
          })
          .then(() =>
            cy
              .fixture('backupContact.json')
              .then(backupContact =>
                cy.request(
                  'POST',
                  `/internal/service_members/${serviceMemberId}/backup_contacts`,
                  backupContact,
                ),
              ),
          ),
      ),
    );
}
