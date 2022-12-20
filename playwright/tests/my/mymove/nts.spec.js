const { test, expect } = require('../../utils/customerTest');
const { signInAsExistingCustomer } = require('../../utils/signIn');
const { buildMoveWithOrders } = require('../../utils/testharness');

test('A customer can create, edit, and delete an NTS shipment', async ({ page, request }) => {
  const move = await buildMoveWithOrders(request);
  const userId = move.Orders.ServiceMember.user_id;
  await signInAsExistingCustomer(page, userId);

  // Navigate to create a new shipment
  await page.getByTestId('shipment-selection-btn').click();
  await page.getByTestId('wizardNextButton').click();
  await page.getByText("It's going into storage for months or years (NTS)Movers pack and ship things to ").click();
  await page.getByTestId('wizardNextButton').click();

  // Fill in form to create NTS shipment
  await page.getByLabel('Preferred pickup date').fill('25 Dec 2022');
  await page.getByLabel('Preferred pickup date').blur();
  await page.getByText('Use my current address').click();
  await page.getByTestId('remarks').fill('Grandfather antique clock');
  await page.getByTestId('wizardNextButton').click();

  // Verify that form submitted
  await expect(await page.getByRole('heading', { level: 1 })).toHaveText('Review your details');
  await expect(await page.getByText('Grandfather antique clock')).toBeVisible();

  // Navigate to edit shipment from the review page
  await page.getByTestId('edit-nts-shipment-btn').click();
  await expect(await page.getByRole('heading', { level: 1 })).toHaveText(
    'Where and when should the movers pick up your things going into storage?',
  );

  // Update form (adding releasing agent)
  await page.getByLabel('First name').fill('Grace');
  await page.getByLabel('Last name').fill('Griffin');
  await page.getByLabel('Phone').fill('2025551234');
  await page.getByLabel('Email').fill('grace.griffin@example.com');
  await page.getByTestId('wizardNextButton').click();

  // Verify that form submitted
  await expect(await page.getByRole('heading', { level: 1 })).toHaveText('Review your details');
  await expect(await page.getByText('Grace Griffin')).toBeVisible();

  // Navigate to homepage and delete shipment
  await page.getByTestId('wizardCancelButton').click();
  await expect(await page.getByRole('heading', { name: 'Leo Spacemen', level: 2 })).toBeVisible();
  await page.getByRole('button', { name: 'Delete' }).click();
  await page.getByTestId('modal').getByTestId('button').click();

  await expect(await page.getByText('The shipment was deleted.')).toBeVisible();
  await expect(await page.getByTestId('stepContainer3').getByText('Set up shipments')).toBeVisible();
});
