/* global cy, */
describe('orders entry', function() {
  beforeEach(() => {
    cy.signInAsNewUser();
  });

  it('progresses thru forms', function() {
    orders();
  });
  //   it.skip('restarts app after every page', function() {
  //     orders(true);
  //   });
});

function createServiceMember() {
  return cy
    .location()
    .should(loc => {
      expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/create/);
    })
    .then(location => {
      ///service-member/cecee251-1ded-47bb-b7cd-a6cf59c832f2/create
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
            return cy.request(
              'PATCH',
              `/internal/service_members/${serviceMemberId}`,
              serviceMember,
            );
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
function orders(restartAfterEachPage) {
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
  // .type('20181020')
  // .blur();
  //   cy
  //     .get('input[placeholder="Date"]')
  //     .last()
  //     .type('20181120')
  //     .blur();

  //   cy
  //     .get('.duty-input-box__input input')
  //     .first()
  //     .type('Fort Worth{downarrow}{enter}', { force: true, delay: 150 });
}
