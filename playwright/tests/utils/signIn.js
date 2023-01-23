/* eslint max-classes-per-file: "off" */
//
// Override max classes per file because only the last class (SignIn)
// is public
//
// @ts-check

// User Types
export const milmoveUserType = 'milmove';
export const TOOOfficeUserType = 'TOO office';
export const TIOOfficeUserType = 'TIO office';
export const QAECSROfficeUserType = 'QAE/CSR office';
export const ServicesCounselorOfficeUserType = 'Services Counselor office';
export const PrimeSimulatorUserType = 'Prime Simulator';

/**
 * App Sign In
 */
class AppSignIn {
  /**
   * Create a base AppSignIn
   * @param {import('@playwright/test').Page} page
   */
  constructor(page) {
    this.page = page;
  }

  /**
   * Sign in as a new user with devlocal
   *
   * @param {string} userType
   */
  async signInAsNewUser(userType) {
    await this.page.goto('/devlocal-auth/login');
    await this.page.locator(`button[data-hook="new-user-login-${userType}"]`).click();
  }

  /**
   * Sign in as existing user with devlocal
   *
   * @param {string} userId
   */
  async signInAsUserWithId(userId) {
    await this.page.goto('/devlocal-auth/login');
    await this.page.locator(`button[value="${userId}"]`).click();
  }
}

/**
 * Admin App Sign In
 * @extends AppSignIn
 */
class AdminAppSignIn extends AppSignIn {
  /**
   * Sign in as new admin user with devlocal
   */
  async newAdminUser() {
    await this.signInAsNewUser('admin');
  }
}

/**
 * Office App Sign In
 * @extends AppSignIn
 */
class OfficeAppSignIn extends AppSignIn {
  /**
   * Sign in as new Service Counselor user with devlocal
   */
  async newServicesCounselorUser() {
    await this.signInAsNewUser(ServicesCounselorOfficeUserType);
  }

  /**
   * Sign in as new TOO user with devlocal
   */
  async newTOOUser() {
    await this.signInAsNewUser(TOOOfficeUserType);
  }

  /**
   * Sign in as new TIO user with devlocal
   */
  async newTIOUser() {
    await this.signInAsNewUser(TIOOfficeUserType);
  }

  /**
   * Sign in as new prime simulator user with devlocal
   */
  async newPrimeSimulatorUser() {
    await this.signInAsNewUser(PrimeSimulatorUserType);
  }

  /**
   * Sign in as new prime simulator user with devlocal
   */
  async newQAECSRUser() {
    await this.signInAsNewUser(QAECSROfficeUserType);
  }

  /**
   * Sign in as existing office user with devlocal
   *
   * @param {string} email
   */
  async existingOfficeUser(email) {
    await this.page.goto('/devlocal-auth/login');
    await this.page.locator('input[name=email]').fill(email);
    await this.page.locator('p', { hasText: 'User Email' }).locator('button').click();
  }
}

/**
 * Customer App Sign In
 * @extends AppSignIn
 */
class CustomerAppSignIn extends AppSignIn {
  /**
   * Sign in as new customer with devlocal
   *
   */
  async newCustomer() {
    await this.signInAsNewUser(milmoveUserType);
  }

  /**
   * Sign in as existing customer with devlocal
   *
   * @param {string} userId
   */
  async existingCustomer(userId) {
    await this.signInAsUserWithId(userId);
  }
}

export class SignIn {
  /**
   * Create a ignIn
   * @param {import('@playwright/test').Page} page
   */
  constructor(page) {
    this.admin = new AdminAppSignIn(page);
    this.office = new OfficeAppSignIn(page);
    this.customer = new CustomerAppSignIn(page);
  }
}

export default SignIn;
