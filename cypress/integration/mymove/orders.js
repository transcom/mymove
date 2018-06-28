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
  cy
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
        cy.window().then(window => {
          serviceMember.user_id = window._state.user.userId;
          cy.request(
            'PATCH',
            `/internal/service_members/${serviceMemberId}`,
            serviceMember,
          );
        }),
      ),
    );
}
function orders(restartAfterEachPage) {
  createServiceMember();
}
