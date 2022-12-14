// @ts-check
const { test, expect } = require('../utils/officeTest');

test.describe('The document viewer', () => {
  test.describe('When not logged in', () => {
    test('shows page not found', async ({ page }) => {
      await page.goto('/moves/foo/documents');
      await expect(page.getByText('Welcome')).toBeVisible();
      // sign in button not in header
      await expect(page.locator('#main').getByRole('button', { name: 'Sign in' })).toBeVisible();
    });
  });
});

test.describe('When user is logged in', () => {
  test('produces error when move cannot be found', async ({ page, officePage }) => {
    await officePage.gotoAndWaitForLoading('/moves/ffffffff-ffff-ffff-ffff-ffffffffffffffff/documents');

    // todo: we want better messages when we are making custom call
    await expect(page.getByText('An error occurred')).toBeVisible();
  });

  test('loads basic information about the move', async ({ page, officePage }) => {
    const move = await officePage.buildInProgressPPMMove();
    const moveId = move.id;
    officePage.gotoAndWaitForLoading(`/moves/${moveId}/documents`);
    await expect(page.getByText('In Progress, PPM')).toBeVisible();
    //       cy.contains('GBXYUI');
    //       cy.contains('1617033988');
  });
  //     test('can upload a new document', async ({page, officePage}) => {
  //       cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents');
  //       cy.get('[data-testid="document-upload-link"]')
  //         .find('a')
  //         .should('have.attr', 'href')
  //         .and('contain', '/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
  //       cy.get('[data-testid="document-upload-link"]');
  //       cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');

  //       cy.contains('Upload a new document');
  //       cy.get('button.submit').should('be.disabled');
  //       cy.get('input[name="title"]').type('super secret info document');
  //       cy.get('select[name="move_document_type"]').select('Other document type');
  //       cy.get('input[name="notes"]').type('burn after reading');
  //       cy.get('button.submit').should('be.disabled');

  //       cy.intercept('/internal/uploads').as('uploadFile');
  //       cy.upload_file('.filepond--root', 'sample-orders.png');
  //       cy.wait('@uploadFile');

  //       cy.intercept('/internal/moves/*/move_documents').as('postDocument');
  //       cy.get('button.submit', { timeout: fileUploadTimeout }).should('not.be.disabled').click();
  //       cy.wait('@postDocument');
  //     });
  //     test('can upload a weight ticket set', async ({page, officePage}) => {
  //       cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
  //       cy.contains('Upload a new document');
  //       cy.get('button.submit').should('be.disabled');
  //       cy.get('input[name="title"]').type('Weight ticket document');
  //       cy.get('select[name="move_document_type"]').select('Weight ticket set');
  //       cy.get('input[name="notes"]').type('burn after reading');
  //       cy.get('button.submit').should('be.disabled');

  //       cy.intercept('/internal/uploads').as('uploadFile');
  //       cy.upload_file('.filepond--root', 'sample-orders.png');
  //       cy.wait('@uploadFile');

  //       cy.intercept('/internal/moves/*/move_documents').as('postDocument');
  //       cy.get('button.submit', { timeout: fileUploadTimeout }).should('not.be.disabled').click();
  //       cy.wait('@postDocument');
  //     });
  //   });
});
