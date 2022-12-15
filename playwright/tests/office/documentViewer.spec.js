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
  // skip tests logged in as PPM Office Users
  test.skip('produces error when move cannot be found', async ({ page, officePage }) => {
    await officePage.signInAsNewPPMUser();
    await officePage.gotoAndWaitForLoading('/moves/ffffffff-ffff-ffff-ffff-ffffffffffffffff/documents');

    // todo: we want better messages when we are making custom call
    await expect(page.getByText('An error occurred')).toBeVisible();
  });

  // skip tests logged in as PPM Office Users
  test.skip('loads basic information about the move', async ({ page, officePage }) => {
    await officePage.signInAsNewPPMUser();
    const move = await officePage.buildInProgressPPMMove();
    const moveId = move.id;
    const firstName = move.Orders.ServiceMember.first_name;
    const lastName = move.Orders.ServiceMember.last_name;
    const { edipi } = move.Orders.ServiceMember;
    await officePage.gotoAndWaitForLoading(`/moves/${moveId}/documents`);
    const lastFirst = `${lastName}, ${firstName}`;
    await expect(page.getByText(lastFirst)).toBeVisible();
    await expect(page.getByText(move.locator)).toBeVisible();
    await expect(page.getByText(edipi)).toBeVisible();
  });

  // skip tests logged in as PPM Office Users
  test.skip('can upload a new document', async ({ page, officePage }) => {
    const move = await officePage.buildInProgressPPMMove();
    const moveId = move.id;
    await officePage.gotoAndWaitForLoading(`/moves/${moveId}/documents`);

    expect(page.url()).toContain('/documents/new');
    await expect(page.getByRole('button', { name: 'Save' })).toBeDisabled();
    await page.locator('input[name="title"]').type('super secret info document');
    await page.locator('select[name="move_document_type"]').selectOption({ label: 'Other document type' });
    await page.locator('input[name="notes"]').type('burn after reading');
    await expect(page.getByRole('button', { name: 'Save' })).toBeDisabled();

    // path is relative to where tests are running. Not sure if there
    // is a better way to configure this
    const fixtureFile = 'playwright/tests/fixtures/sample-orders.png';
    await page.locator('input[name="filepond"]').setInputFiles(fixtureFile);
    await page.waitForLoadState('networkidle');
    await expect(page.getByRole('button', { name: 'Save' })).not.toBeDisabled();
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page.getByText('All Documents (1)')).toHaveCount(1);
  });

  //     test('can upload a weight ticket set', async ({page, officePage}) => {
  //       cy.patientVisit('/moves/c9df71f2-334f-4f0e-b2e7-050ddb22efa1/documents/new');
  //       cy.contains('Upload a new document');
  //       await page.locator('button.submit').should('be.disabled');
  //       await page.locator('input[name="title"]').type('Weight ticket document');
  //       await page.locator('select[name="move_document_type"]').select('Weight ticket set');
  //       await page.locator('input[name="notes"]').type('burn after reading');
  //       await page.locator('button.submit').should('be.disabled');

  //       cy.intercept('/internal/uploads').as('uploadFile');
  //       cy.upload_file('.filepond--root', 'sample-orders.png');
  //       cy.wait('@uploadFile');

  //       cy.intercept('/internal/moves/*/move_documents').as('postDocument');
  //       await page.locator('button.submit', { timeout: fileUploadTimeout }).should('not.be.disabled').click();
  //       cy.wait('@postDocument');
  //     });
  //   });
});
