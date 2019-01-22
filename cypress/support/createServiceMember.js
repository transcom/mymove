/* global cy, Cypress*/
export default function createServiceMember() {
  //make sure state has loaded
  cy.location().should(loc => {
    expect(loc.pathname).to.match(/^\/service-member\/[^/]+\/create/);
  });

  return cy.request('/internal/users/logged_in').then(result => {
    const userId = Cypress._.get(result, 'body.id');
    const serviceMemberId = Cypress._.get(result, 'body.service_member.id');
    return cy.fixture('serviceMember.json').then(serviceMember => {
      serviceMember.user_id = userId;
      return cy.request('internal/duty_stations?search=ft%20car').then(result => {
        serviceMember.current_station_id = Cypress._.get(result, 'body.[0].id');
        return cy.getCookie('masked_gorilla_csrf').then(token => {
          return cy
            .request({
              method: 'PATCH',
              url: `/internal/service_members/${serviceMemberId}`,
              headers: { 'X-CSRF-Token': token.value },
              body: serviceMember,
            })
            .then(() =>
              cy.fixture('backupContact.json').then(backupContact =>
                cy.request({
                  method: 'POST',
                  url: `/internal/service_members/${serviceMemberId}/backup_contacts`,
                  headers: { 'X-CSRF-Token': token.value },
                  body: backupContact,
                }),
              ),
            );
        });
      });
    });
  });
}
