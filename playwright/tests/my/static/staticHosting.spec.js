// @ts-check
const { test, expect } = require('@playwright/test');

test('pdf download', async ({ page, request }) => {
  await page.goto('/');
  const response = await request.get('/downloads/direct_deposit_form.pdf');
  expect(response.ok()).toBeTruthy();
  const contentType = response.headers()['content-type'];
  expect(contentType).toEqual('application/pdf');
});

test('disallow pdf upload', async ({ page, request }) => {
  await page.goto('/');
  const response = await request.post('/downloads/direct_deposit_form.pdf');

  expect(response.status()).toEqual(405);
});
