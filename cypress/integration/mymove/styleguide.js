/* global cy*/

describe('styleguide', function() {
  it('developer can navigate to styleguide route', function() {
    userVisitsStyleguideRoute();
  });
});

function userVisitsStyleguideRoute() {
  cy.visit('/sm_style_guide');
  cy.contains('Placeholder style guide');
}
