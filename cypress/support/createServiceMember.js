/* global cy, Cypress*/
export default function createServiceMember() {
  //make sure state has loaded
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/create/);
  });
  return cy.window().then(window => {
    const state = window._state;
    const serviceMemberId = state.serviceMember.currentServiceMember.id;
    const userId = state.user.userId;
    return cy.fixture('serviceMember.json').then(serviceMember => {
      serviceMember.user_id = userId;
      return cy
        .request('internal/duty_stations?search=ft%20car')
        .then(result => {
          serviceMember.current_station_id = Cypress._.get(
            result,
            'body.[0].id',
          );
          return cy
            .request(
              'PATCH',
              `/internal/service_members/${serviceMemberId}`,
              serviceMember,
            )
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
            );
        });
    });
  });
}
