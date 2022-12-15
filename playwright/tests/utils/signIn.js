// User Types
export const milmoveUserType = 'milmove';
export const PPMOfficeUserType = 'PPM office';
export const TOOOfficeUserType = 'TOO office';
export const TIOOfficeUserType = 'TIO office';
export const QAECSROfficeUserType = 'QAE/CSR office';
export const ServicesCounselorOfficeUserType = 'Services Counselor office';
export const PrimeSimulatorUserType = 'Prime Simulator';

async function signInAsNewUser(page, userType) {
  await page.goto('/devlocal-auth/login');
  await page.locator(`button[data-hook="new-user-login-${userType}"]`).click();
}

export async function signIntoAdminAsNewAdminUser(page) {
  await signInAsNewUser(page, 'admin');
}

export async function signIntoOfficeAsNewPPMUser(page) {
  await signInAsNewUser(page, PPMOfficeUserType);
}

export async function signIntoOfficeAsNewServicesCounselorUser(page) {
  await signInAsNewUser(page, ServicesCounselorOfficeUserType);
}

export async function signIntoOfficeAsNewTOOUser(page) {
  await signInAsNewUser(page, TOOOfficeUserType);
}

export async function signIntoOfficeAsNewTIOUser(page) {
  await signInAsNewUser(page, TIOOfficeUserType);
}

async function signInAsUserWithIdAndType(page, userId) {
  await page.goto('/devlocal-auth/login');
  await page.locator(`button[value="${userId}"]`).click();
}

export async function signInAsExistingCustomer(page, userId) {
  return signInAsUserWithIdAndType(page, userId);
}

export async function signInAsExistingOfficeUser(page, email) {
  await page.goto('/devlocal-auth/login');
  await page.locator('input[name=email]').fill(email);
  await page.locator('p', { hasText: 'User Email' }).locator('button').click();
}
