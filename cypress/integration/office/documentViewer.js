import { fileUploadTimeout } from '../../support/constants';

describe('The document viewer', function () {
  describe('When not logged in', function () {
    beforeEach(() => {
      cy.prepareOfficeApp();
      cy.logout();
    });

    it('shows page not found', function () {
      cy.patientVisit('/moves/foo/documents');
      cy.contains('Welcome');
      cy.contains('Sign in');
    });
  });
});
