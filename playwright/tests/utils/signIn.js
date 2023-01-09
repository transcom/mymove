// @ts-check

// User Types
export const milmoveUserType = 'milmove';
export const PPMOfficeUserType = 'PPM office';
export const TOOOfficeUserType = 'TOO office';
export const TIOOfficeUserType = 'TIO office';
export const QAECSROfficeUserType = 'QAE/CSR office';
export const ServicesCounselorOfficeUserType = 'Services Counselor office';
export const PrimeSimulatorUserType = 'Prime Simulator';

/**
 * Sign In
 *
 * @param {import('@playwright/test').Page} page
 */
export function newSignIn(page) {
  /**
   * Sign in as a new user with devlocal
   *
   * @param {string} userType
   */
  const signInAsNewUser = async (userType) => {
    await page.goto('/devlocal-auth/login');
    await page.locator(`button[data-hook="new-user-login-${userType}"]`).click();
  };

  /**
   * Sign in as existing user with devlocal
   *
   * @param {string} userId
   */
  const signInAsUserWithId = async (userId) => {
    await page.goto('/devlocal-auth/login');
    await page.locator(`button[value="${userId}"]`).click();
  };

  return {
    admin: {
      /**
       * Sign in as new admin user with devlocal
       */
      async newAdminUser() {
        await signInAsNewUser('admin');
      },
    },
    office: {
      /**
       * Sign in as new Service Counselor user with devlocal
       */
      async newServicesCounselorUser() {
        await signInAsNewUser(ServicesCounselorOfficeUserType);
      },

      /**
       * Sign in as new TOO user with devlocal
       */
      async newTOOUser() {
        await signInAsNewUser(TOOOfficeUserType);
      },

      /**
       * Sign in as new TIO user with devlocal
       */
      async newTIOUser() {
        await signInAsNewUser(TIOOfficeUserType);
      },

      /**
       * Sign in as new prime simulator user with devlocal
       */
      async newPrimeSimulatorUser() {
        await signInAsNewUser(PrimeSimulatorUserType);
      },
      /**
       * Sign in as existing office user with devlocal
       *
       * @param {string} email
       */
      async existingOfficeUser(email) {
        await page.goto('/devlocal-auth/login');
        await page.locator('input[name=email]').fill(email);
        await page.locator('p', { hasText: 'User Email' }).locator('button').click();
      },
    },
    customer: {
      /**
       * Sign in as new customer with devlocal
       *
       */
      async newCustomer() {
        await signInAsNewUser(milmoveUserType);
      },
      /**
       * Sign in as existing customer with devlocal
       *
       * @param {string} userId
       */
      async existingCustomer(userId) {
        await signInAsUserWithId(userId);
      },
    },
  };
}
export default newSignIn;
