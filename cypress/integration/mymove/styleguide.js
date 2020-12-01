describe('styleguide', function () {
  it('developer can navigate to styleguide route', function () {
    cy.request('/sm_style_guide').should((response) => {
      expect(response.status).to.eq(200);
    });
  });
});
