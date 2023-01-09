// @ts-check
const { test, expect } = require('../utils/officeTest');

const csrfForbiddenMsg = 'Forbidden - CSRF token invalid\n';
const csrfForbiddenRespCode = 403;
const unmaskedCsrfCookieName = '_gorilla_csrf';
const maskedCsrfCookieName = 'masked_gorilla_csrf';

// CSRF protection is turned on for all routes.
// We can test with the local dev login that uses POST
test.describe('testing CSRF protection', () => {
  test('can successfully dev login with both unmasked and masked token', async ({ page, officePage }) => {
    await officePage.signInAsNewTIOUser();
    await expect(page.getByRole('heading', { name: 'Payment requests' })).toBeVisible();
  });

  test('cannot dev login with masked token only', async ({ browser, officePage }) => {
    // create a user we can log in with later
    const user = await officePage.testHarness.buildOfficeUserWithTOOAndTIO();
    // need a browser context to look at cookies
    // this creates a new icognito browser context, which works great
    // for this test
    const context = await browser.newContext();

    // use the new browser context to make an API request using cookies
    const response = await context.request.get('/internal/users/is_logged_in');
    expect(response.ok()).toBeTruthy();

    const { isLoggedIn } = await response.json();
    expect(isLoggedIn).toBeFalsy();

    const allCookies = await context.cookies();
    expect(allCookies.filter((cookie) => cookie.name === unmaskedCsrfCookieName)).toHaveLength(1);

    // add back all cookies but the unmasked csrf one
    context.clearCookies();
    const filteredCookies = allCookies.filter((cookie) => cookie.name !== unmaskedCsrfCookieName);
    context.addCookies(filteredCookies);

    const maskedCsrfCookie = allCookies.filter((cookie) => cookie.name === maskedCsrfCookieName)[0];
    expect(maskedCsrfCookie).toBeDefined();

    const maskedValue = maskedCsrfCookie.value;
    expect(maskedValue).toBeTruthy();

    // make a request with the wrong csrf token in a header
    const failedResponse = await context.request.post('/devlocal-auth/login', {
      form: { id: user.id, userType: 'PPM office' },
      headers: { 'X-CSRF-TOKEN': maskedValue },
    });

    expect(failedResponse.status()).toEqual(csrfForbiddenRespCode);
    expect(await failedResponse.text()).toEqual(csrfForbiddenMsg);
  });

  test('cannot dev login with unmasked token only', async ({ page, officePage }) => {
    // create a user we can log in with later
    const user = await officePage.testHarness.buildOfficeUserWithTOOAndTIO();

    const context = page.context();

    // use the new browser context to make an API request using cookies
    const response = await context.request.get('/internal/users/is_logged_in');
    expect(response.ok()).toBeTruthy();

    const { isLoggedIn } = await response.json();
    expect(isLoggedIn).toBeFalsy();

    const allCookies = await context.cookies();
    expect(allCookies.filter((cookie) => cookie.name === unmaskedCsrfCookieName)).toHaveLength(1);
    expect(allCookies.filter((cookie) => cookie.name === maskedCsrfCookieName)).toHaveLength(1);

    // add back all cookies but the masked csrf one
    context.clearCookies();
    const filteredCookies = allCookies.filter((cookie) => cookie.name !== maskedCsrfCookieName);
    context.addCookies(filteredCookies);

    // make a request with the wrong csrf token in a header
    const failedResponse = await context.request.post('/devlocal-auth/login', {
      form: { id: user.id, userType: 'PPM office' },
      headers: { 'X-CSRF-TOKEN': 'bad' },
    });

    expect(failedResponse.status()).toEqual(csrfForbiddenRespCode);
    expect(await failedResponse.text()).toEqual(csrfForbiddenMsg);

    await context.close();
  });

  test('cannot dev login without unmasked and masked token', async ({ page, officePage }) => {
    // create a user we can log in with later
    const user = await officePage.testHarness.buildOfficeUserWithTOOAndTIO();

    const context = page.context();

    // use the new browser context to make an API request using cookies
    const response = await context.request.get('/internal/users/is_logged_in');
    expect(response.ok()).toBeTruthy();

    const { isLoggedIn } = await response.json();
    expect(isLoggedIn).toBeFalsy();

    const allCookies = await context.cookies();
    expect(allCookies.filter((cookie) => cookie.name === unmaskedCsrfCookieName)).toHaveLength(1);
    expect(allCookies.filter((cookie) => cookie.name === maskedCsrfCookieName)).toHaveLength(1);

    // add back all cookies except the csrf ones
    context.clearCookies();
    const filteredCookies = allCookies
      .filter((cookie) => cookie.name !== maskedCsrfCookieName)
      .filter((cookie) => cookie.name !== unmaskedCsrfCookieName);
    context.addCookies(filteredCookies);

    // make a request with the wrong csrf token in a header
    const failedResponse = await context.request.post('/devlocal-auth/login', {
      form: { id: user.id, userType: 'PPM office' },
      headers: { 'X-CSRF-TOKEN': 'bad' },
    });

    expect(failedResponse.status()).toEqual(csrfForbiddenRespCode);
    expect(await failedResponse.text()).toEqual(csrfForbiddenMsg);

    await context.close();
  });
});

test.describe('testing CSRF protection updating move info', () => {
  test('tests updating user profile with proper tokens', async ({ page, officePage }) => {
    // blargh, too bad
    test.slow();

    const move = await officePage.testHarness.buildHHGMoveWithServiceItemsAndPaymentRequestsAndFiles();
    await officePage.signInAsNewTOOUser();

    // click on newly created move
    await page.getByText(move.locator).click();
    await officePage.waitForLoading();

    await page.getByTestId('edit-orders').click();
    await officePage.waitForLoading();

    await page.getByLabel('Orders number').type('CSRFTEST');

    await page.getByRole('button', { name: 'Save' }).click();
    await officePage.waitForLoading();

    await expect(page.getByText('CSRFTEST')).toBeVisible();
  });

  test('tests updating user profile without masked token', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildHHGMoveWithServiceItemsAndPaymentRequestsAndFiles();
    await officePage.signInAsNewTOOUser();

    // click on newly created move
    await page.getByText(move.locator).click();
    await officePage.waitForLoading();

    await page.getByTestId('edit-orders').click();
    await officePage.waitForLoading();

    await page.getByLabel('Orders number').type('CSRFTEST');

    const context = page.context();

    const allCookies = await context.cookies();
    expect(allCookies.filter((cookie) => cookie.name === unmaskedCsrfCookieName)).toHaveLength(1);
    expect(allCookies.filter((cookie) => cookie.name === maskedCsrfCookieName)).toHaveLength(1);

    // add back all cookies but the masked csrf one
    context.clearCookies();
    const filteredCookies = allCookies.filter((cookie) => cookie.name !== maskedCsrfCookieName);
    context.addCookies(filteredCookies);

    // wait for the response before clicking so we don't miss it
    const values = await Promise.all([
      page.waitForResponse('**/ghc/v1/orders/*'),
      page.getByRole('button', { name: 'Save' }).click(),
    ]);

    expect(values.length).toEqual(2);
    const r = values[0];
    expect(r.status()).toEqual(403);

    await context.close();
  });
});
