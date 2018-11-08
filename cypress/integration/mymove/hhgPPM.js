/* global cy */

const serviceMemberSignsIn = uuid => {
  cy.signInAsUser(uuid);

  cy.location().should(loc => {
    expect(loc.pathname).to.match(/\//);
  });
};

const serviceMemberAddsPPMToHHG = () => {
  cy.get('.step a').contains('Add PPM Shipment');

  // should('have.attr', 'href', /^\/moves\/[^/]+\/hhg-ppm-start/)
};

describe('completing the ppm flow', function() {
  //tear down currently means doing this:
  //update moves set status='DRAFT';
  //delete from personally_procured_moves

  it('progresses thru forms', function() {
    serviceMemberSignsIn('7980f0cf-63e3-4722-b5aa-ba46f8f7ac64');

    serviceMemberAddsPPMToHHG();
  });
});
