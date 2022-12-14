async function signInAsNewUser(page, userType) {
  await page.goto('/devlocal-auth/login');
  await page.locator(`button[data-hook=new-user-login-${userType}]`).click();
}

export async function signInAsNewAdminUser(page) {
  await signInAsNewUser(page, 'admin');
}

// User Types
export const milmoveUserType = 'milmove';
export const PPMOfficeUserType = 'PPM office';
export const TOOOfficeUserType = 'TOO office';
export const TIOOfficeUserType = 'TIO office';
export const QAECSROfficeUserType = 'QAE/CSR office';
export const ServicesCounselorOfficeUserType = 'Services Counselor office';
export const PrimeSimulatorUserType = 'Prime Simulator';

async function signInAsUserWithIdAndType(page, userId) {
  await page.goto('/devlocal-auth/login');
  await page.locator(`button[value="${userId}"]`).click();
}

export async function signInAsExistingCustomer(page, userId) {
  return signInAsUserWithIdAndType(page, userId);
}
