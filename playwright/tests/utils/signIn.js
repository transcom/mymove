export async function signInAsNewUser(page, userType) {
  await page.goto('/devlocal-auth/login');
  await page.locator(`button[data-hook=new-user-login-${userType}]`).click();
}

export async function signInAsNewAdminUser(page) {
  await signInAsNewUser(page, 'admin');
}
