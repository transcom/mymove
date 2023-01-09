// @ts-check
const { test, expect } = require('../../utils/customerTest');

/**
 * @param {import('@playwright/test').Page} page
 */
async function customerChoosesAnHHGMove(page) {
  await page.getByRole('button', { name: 'Set up your shipments' }).click();
  await page.getByRole('button', { name: 'Next' }).click();

  await page.getByText('HHG').click();
  await page.getByRole('button', { name: 'Next' }).click();
}

/**
 * @param {import('@playwright/test').Page} page
 */
async function customerSetsUpAnHHGMove(page) {
  await expect(page.getByRole('button', { name: 'Next' })).toBeDisabled();

  await page.getByLabel('Preferred pickup date').focus();
  await page.getByLabel('Preferred pickup date').blur();
  await expect(page.locator('.usa-error-message')).toContainText('Required');
  await page.locator('input[name="pickup.requestedDate"]').type('08/02/2020');
  await page.locator('input[name="pickup.requestedDate"]').blur();
  await expect(page.locator('.usa-error-message')).toHaveCount(0);

  await expect(page.locator('button[data-testid="wizardNextButton"]')).toBeDisabled();

  // should be empty before using "Use current residence" checkbox
  await expect(page.locator('input[name="pickup.address.streetAddress1"]')).toBeEmpty();
  await expect(page.locator('input[name="pickup.address.city"]')).toBeEmpty();
  await expect(page.locator('input[name="pickup.address.postalCode"]')).toBeEmpty();

  // should have expected "Required" error for required fields
  await page.locator('input[name="pickup.address.streetAddress1"]').focus();
  await page.locator('input[name="pickup.address.streetAddress1"]').blur();
  await expect(page.locator('.usa-error-message')).toContainText('Required');
  await page.locator('input[name="pickup.address.streetAddress1"]').type('Some address');
  await expect(page.locator('.usa-error-message')).toHaveCount(0);

  await page.locator('input[name="pickup.address.city"]').focus();
  await page.locator('input[name="pickup.address.city"]').blur();
  await expect(page.locator('.usa-error-message')).toContainText('Required');
  await page.locator('input[name="pickup.address.city"]').type('Some city');
  await expect(page.locator('.usa-error-message')).toHaveCount(0);

  await page.locator('select[name="pickup.address.state"]').focus();
  await page.locator('select[name="pickup.address.state"]').blur();
  await expect(page.locator('.usa-error-message')).toContainText('Required');
  await page.locator('select[name="pickup.address.state"]').selectOption('CA');
  await expect(page.locator('.usa-error-message')).toHaveCount(0);

  await page.locator('input[name="pickup.address.postalCode"]').focus();
  await page.locator('input[name="pickup.address.postalCode"]').blur();
  await expect(page.locator('.usa-error-message')).toContainText('Required');
  await page.locator('input[name="pickup.address.postalCode"]').type('9');
  await page.locator('input[name="pickup.address.postalCode"]').blur();
  await expect(page.locator('.usa-error-message')).toContainText('Must be valid zip code');
  await page.locator('input[name="pickup.address.postalCode"]').type('1111');
  await page.locator('input[name="pickup.address.postalCode"]').blur();
  await expect(page.locator('.usa-error-message')).toHaveCount(0);

  // Next button disabled
  await expect(page.locator('button[data-testid="wizardNextButton"]')).toBeDisabled();

  // overwrites data typed from above
  //
  // because of the styling of this input item, we cannot use a
  // css locator for the input item and then click it
  //
  // The styling is very similar to the issue described in
  //
  // https://github.com/microsoft/playwright/issues/3688
  //
  await page.getByText('Use my current address').click();

  // secondary pickup location
  // same click issue as above
  await page.getByRole('group', { name: 'Pickup location' }).getByText('Yes').click();

  await page.locator('input[name="secondaryPickup.address.streetAddress1"]').type('123 Some address');
  await page.locator('input[name="secondaryPickup.address.city"]').type('Some city');
  await page.locator('select[name="secondaryPickup.address.state"]').selectOption('CA');
  await page.locator('input[name="secondaryPickup.address.postalCode"]').type('90210');
  await page.locator('input[name="secondaryPickup.address.postalCode"]').blur();

  // releasing agent
  await page.locator('input[name="pickup.agent.firstName"]').type('John');
  await page.locator('input[name="pickup.agent.lastName"]').type('Lee');
  await page.locator('input[name="pickup.agent.phone"]').type('9999999999');
  await page.locator('input[name="pickup.agent.email"]').type('ron');
  await page.locator('input[name="pickup.agent.email"]').blur();
  await expect(page.locator('.usa-error-message')).toContainText('Must be valid email');
  // playwright overwrites the text field when using .type(). We could
  // use .fill(oldValue + newValue), but that doesn't seem worthwhile
  // here
  //
  // See https://github.com/microsoft/playwright/issues/10753
  await page.locator('input[name="pickup.agent.email"]').type('ron@example.com');
  await expect(page.locator('.usa-error-message')).toHaveCount(0);

  await expect(page.locator('button[data-testid="wizardNextButton"]')).toBeDisabled();

  // requested delivery date
  await page.locator('input[name="delivery.requestedDate"]').first().type('09/20/2020');
  await page.locator('input[name="delivery.requestedDate"]').first().blur();

  // checks has delivery address (default does not have delivery
  // address)
  // same click issue as above
  await page.getByRole('group', { name: 'Delivery location' }).getByText('Yes').click();
  // delivery location
  await page.locator('input[name="delivery.address.streetAddress1"]').type('412 Avenue M');
  await page.locator('input[name="delivery.address.streetAddress2"]').type('#3E');
  await page.locator('input[name="delivery.address.city"]').type('Los Angeles');
  await page.locator('select[name="delivery.address.state"]').selectOption('CA');
  await page.locator('input[name="delivery.address.postalCode"]').type('91111');
  await page.locator('input[name="delivery.address.postalCode"]').blur();

  // secondary delivery location
  // same click issue as above
  await page
    .getByRole('group', { name: 'Delivery location' })
    .locator('fieldset:has-text("Second delivery location")')
    .getByText('Yes')
    .click();

  await page.locator('input[name="secondaryDelivery.address.streetAddress1"]').type('123 Oak Street');
  await page.locator('input[name="secondaryDelivery.address.streetAddress2"]').type('5A');
  await page.locator('input[name="secondaryDelivery.address.city"]').type('San Diego');
  await page.locator('select[name="secondaryDelivery.address.state"]').selectOption('CA');
  await page.locator('input[name="secondaryDelivery.address.postalCode"]').type('91111');
  await page.locator('input[name="secondaryDelivery.address.postalCode"]').blur();

  // releasing agent
  await page.locator('input[name="delivery.agent.firstName"]').type('John');
  await page.locator('input[name="delivery.agent.lastName"]').type('Lee');
  await page.locator('input[name="delivery.agent.phone"]').type('9999999999');
  await page.locator('input[name="delivery.agent.email"]').type('ron');
  await page.locator('input[name="delivery.agent.email"]').blur();
  await expect(page.locator('.usa-error-message')).toContainText('Must be valid email');
  // same type/fill problem as above
  await page.locator('input[name="delivery.agent.email"]').type('ron@example.com');
  await expect(page.locator('.usa-error-message')).toHaveCount(0);

  // customer remarks
  await page.locator('[data-testid="remarks"]').first().type('some customer remark');
  await page.getByRole('button', { name: 'Next' }).click();
}

test('A customer following HHG Setup flow', async ({ page, customerPage }) => {
  const move = await customerPage.testHarness.buildSpouseProGearMove();
  const userId = move.Orders.ServiceMember.user_id;
  await customerPage.signIn.customer.existingCustomer(userId);

  await customerChoosesAnHHGMove(page);
  await customerSetsUpAnHHGMove(page);

  // beforeEach(() => {
  //   cy.intercept('POST', '**/internal/service_members').as('createServiceMember');
  //   cy.intercept('PATCH', '**/internal/mto-shipments/**').as('patchShipment');
  //   cy.intercept('**/internal/moves/**/mto_shipments').as('getMTOShipments');
  //   cy.intercept('**/internal/users/logged_in').as('getLoggedInUser');
  // });

  // it('can create an HHG shipment, review and edit details, and submit their move', function () {
  //   // profile@comple.te
  //   const userId = '3b9360a3-3304-4c60-90f4-83d687884077';
  //   cy.apiSignInAsUser(userId);
  //   customerChoosesAnHHGMove();
  //   customerSetsUpAnHHGMove();
  //   customerReviewsMoveDetailsAndEditsHHG();
  //   customerSubmitsMove();
  // });
});
