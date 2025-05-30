/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../utils/office/officeTest';
import { roleTypes } from '../../../src/constants/userRoles';

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
    await officePage.signInAsNewTOOUser();
    await expect(page.getByText('All moves')).toBeVisible();
    expect(new URL(page.url()).pathname).toBe('/move-queue');
  });

  test('redirects TIO to TIO homepage', async ({ page, officePage }) => {
    await officePage.signInAsNewTIOUser();
    await expect(page.getByRole('heading', { name: 'Payment requests' })).toBeVisible();
    expect(new URL(page.url()).pathname).toBe('/payment-requests');
  });

  test('redirects Services Counselor to Services Counselor homepage', async ({ page, officePage }) => {
    await officePage.signInAsNewServicesCounselorUser();
    await expect(page.getByText('Moves')).toBeVisible();
    await expect(page.getByRole('link', { name: 'Counseling' })).toBeVisible();
    expect(new URL(page.url()).pathname).toBe('/counseling');
  });

  test.describe('multiple role selection', () => {
    const landingChecks = {
      [roleTypes.TOO]: { heading: 'All moves' },
      [roleTypes.TIO]: { heading: 'Payment requests' },
      [roleTypes.SERVICES_COUNSELOR]: { heading: 'Moves' },
      [roleTypes.PRIME_SIMULATOR]: { heading: 'Moves available to Prime' },
      [roleTypes.QAE]: { heading: 'Search for a move' },
      [roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE]: { heading: 'Search for a move' },
      [roleTypes.HQ]: { heading: 'All moves' },
      [roleTypes.CONTRACTING_OFFICER]: { heading: 'Search for a move' },
    };

    test('can switch between Multirole roles', async ({ page, officePage }) => {
      test.slow();
      await officePage.signInAsNewMultiRoleUser();
      await officePage.waitForLoading();
      for (const [role, check] of Object.entries(landingChecks)) {
        await test.step(`switch to ${role}`, async () => {
          // open select app
          await page.getByText('Change user role').click();
          await expect(page).toHaveURL(/\/select-application/);

          // change roles
          await page.getByText(`Select ${role}`).click();
          await page.waitForTimeout(100); // SC component is flaky, let it load up here
          await officePage.waitForLoading();

          // assert landing page
          await expect(page.getByRole('heading', { name: check.heading })).toBeVisible();
        });
      }
    });
    test('can switch between TOO & TIO roles', async ({ page, officePage }) => {
      test.slow();
      await officePage.signInAsNewTIOAndTOOUser();
      await page.getByText('Change user role').click();
      expect(page.url()).toContain('/select-application');
      await page.getByText('Select task_invoicing_officer').click();
      await officePage.waitForLoading();
      await expect(page.getByRole('heading', { name: 'Payment requests' })).toBeVisible();

      await page.getByText('Change user role').click();
      expect(page.url()).toContain('/select-application');
      await page.getByText('Select task_ordering_officer').click();
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
