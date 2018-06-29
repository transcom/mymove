/* global cy, Cypress*/
export default function createServiceMember() {
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
