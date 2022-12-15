// @ts-check
const { test, expect } = require('../utils/officeTest');

test.describe('Office Home Page', () => {
  test('successfully loads when not logged in', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByText('office.move.mil')).toBeVisible();
    await expect(page.locator('button').getByText('Sign in', { exact: true })).toHaveCount(2);
  });

  // skip tests logged in as PPM Office Users
  test.skip('open accepted shipments queue and see moves', async ({ page }) => {
    await page.goto('/');
    // cy.signInAsNewPPMOfficeUser();
    // cy.patientVisit('/queues/all');
    // cy.location().should((loc) => {
    //   expect(loc.pathname).to.match(/^\/queues\/all/);
    // });
    // cy.get('[data-testid=locator]').contains('NOSHOW').should('not.exist');
  });

  // skip tests logged in as PPM Office Users
  test.skip('office user can use a single click to view move info', async ({ page }) => {
    await page.goto('/');
    // cy.waitForReactTableLoad();
    // cy.get('[data-testid=queueTableRow]:first').click();
    // cy.url().should('include', '/moves/');
  });
});

test.describe('Office authorization', () => {
  test('redirects TOO to TOO homepage', async ({ page, officePage }) => {
    officePage.signInAsNewTOOUser();
    await expect(page.getByText('All moves')).toBeVisible();
    expect(new URL(page.url()).pathname).toBe('/');
  });

  test('redirects TIO to TIO homepage', async ({ page, officePage }) => {
    officePage.signInAsNewTIOUser();
    await expect(page.getByRole('heading', { name: 'Payment requests' })).toBeVisible();
    expect(new URL(page.url()).pathname).toBe('/');
  });

  test('redirects Services Counselor to Services Counselor homepage', async ({ page, officePage }) => {
    officePage.signInAsNewServicesCounselorUser();
    await expect(page.getByText('Moves')).toBeVisible();
    await expect(page.getByRole('link', { name: 'Counseling' })).toBeVisible();
    expect(new URL(page.url()).pathname).toBe('/');
  });

  // test('redirects PPM office user to old office queue', async ({page}) => {
  //   cy.signInAsNewPPMOfficeUser();
  //   await expect(page.getByText('New moves')).toBeVisible();
  //   cy.url().should('eq', officeBaseURL + '/');
  // });

  test.describe('multiple role selection', () => {
    test('can switch between TOO & TIO roles', async ({ page, officePage }) => {
      await officePage.signInAsNewTIOAndTOOUser();
      await expect(page.getByText('All moves')).toBeVisible(); // TOO home
      await page.getByText('Change user role').click();
      expect(page.url()).toContain('/select-application');
      await page.getByText('Select transportation_invoicing_officer').click();
      await officePage.waitForLoading();
      await expect(page.getByRole('heading', { name: 'Payment requests' })).toBeVisible();

      await page.getByText('Change user role').click();
      expect(page.url()).toContain('/select-application');
      await page.getByText('Select transportation_ordering_officer').click();
      await officePage.waitForLoading();
      await expect(page.getByText('All moves')).toBeVisible();
    });
  });
});

test.describe('Queue staleness indicator', () => {
  // skip tests logged in as PPM Office Users
  test.skip('displays the correct time ago text', async ({ page }) => {
    await page.goto('/');
    // cy.clock();
    // cy.signInAsNewPPMOfficeUser();
    // cy.patientVisit('/queues/all');
    // cy.get('[data-testid=staleness-indicator]').contains('Last updated a few seconds ago');
    // cy.tick(120000);
    // cy.get('[data-testid=staleness-indicator]').contains('Last updated 2 mins ago');
  });
});
