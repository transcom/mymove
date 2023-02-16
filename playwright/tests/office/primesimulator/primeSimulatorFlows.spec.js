/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../../utils/officeTest';

// I was trying to avoid importing moment.js, these can definitely be
// improved
/**
 * @param {number} daysInterval
 */
const formatRelativeDate = (daysInterval) => {
  const relativeDate = new Date();
  relativeDate.setDate(relativeDate.getDate() + daysInterval);
  const formattedDay = relativeDate.toLocaleDateString(undefined, { day: '2-digit' });
  const formattedMonth = relativeDate.toLocaleDateString(undefined, {
    month: 'short',
  });
  const formattedYear = relativeDate.toLocaleDateString(undefined, {
    year: 'numeric',
  });
  const formattedDate = `${formattedDay} ${formattedMonth} ${formattedYear}`;

  return {
    relativeDate,
    formattedDate,
  };
};

/**
 * @param {Date} date
 */
const formatNumericDate = (date) => {
  const formattedDay = date.toLocaleDateString(undefined, { day: '2-digit' });
  const formattedMonth = date.toLocaleDateString(undefined, {
    month: '2-digit',
  });
  const formattedYear = date.toLocaleDateString(undefined, {
    year: 'numeric',
  });

  return [formattedYear, formattedMonth, formattedDay].join('-');
};

test.describe('Prime simulator user', () => {
  test('is able to update a shipment', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildPrimeSimulatorMoveNeedsShipmentUpdate();

    await officePage.signInAsNewPrimeSimulatorUser();
    const moveLocator = move.locator;
    const moveID = move.id;

    // wait for the the available moves page to load
    // select the move from the list
    await page.getByText(moveLocator).click();
    await officePage.waitForLoading();
    await expect(page.getByText(moveLocator)).toBeVisible();
    expect(page.url()).toContain(`/simulator/moves/${moveID}/details`);
    // waits for the move details page to load
    await page.getByText('Update Shipment').click();

    // waits for the update shipment page to load
    expect(page.url()).toContain(`/simulator/moves/${moveID}/shipments`);

    const { relativeDate: scheduledDeliveryDate, formattedDate: formattedScheduledDeliveryDate } =
      formatRelativeDate(11);
    await page.locator('input[name="scheduledDeliveryDate"]').type(formattedScheduledDeliveryDate);
    await page.locator('input[name="scheduledDeliveryDate"]').blur();
    const { relativeDate: actualDeliveryDate, formattedDate: formattedActualDeliveryDate } = formatRelativeDate(12);
    await page.locator('input[name="actualDeliveryDate"]').type(formattedActualDeliveryDate);
    await page.locator('input[name="actualDeliveryDate"]').blur();
    // there must be sufficient time prior to the pickup dates to update the estimated weight
    const { relativeDate: scheduledPickupDate, formattedDate: formattedScheduledPickupDate } = formatRelativeDate(11);
    await page.locator('input[name="scheduledPickupDate"]').type(formattedScheduledPickupDate);
    await page.locator('input[name="scheduledPickupDate"]').blur();
    const { relativeDate: actualPickupDate, formattedDate: formattedActualPickupDate } = formatRelativeDate(12);
    await page.locator('input[name="actualPickupDate"]').type(formattedActualPickupDate);
    await page.locator('input[name="actualPickupDate"]').blur();
    // update shipment does not require these fields but we need actual weight to create a payment request, we could
    // perform multiple updates.
    await page.locator('input[name="estimatedWeight"]').type('{backspace}7500');
    await page.locator('input[name="actualWeight"]').type('{backspace}8000');
    await page.locator('input[name="destinationAddress.streetAddress1"]').type('142 E Barrel Hoop Circle');
    await page.locator('input[name="destinationAddress.city"]').type('Joshua Tree');
    await page.locator('select[name="destinationAddress.state"]').selectOption({ label: 'CA' });
    await page.locator('input[name="destinationAddress.postalCode"]').type('92252');

    await page.getByText('Save').click();
    await expect(page.getByText('Successfully updated shipment')).toHaveCount(1);
    expect(page.url()).toContain(`/simulator/moves/${moveID}/details`);
    // If you added another shipment to the move you would want to scope these with within()
    await expect(page.getByText(`Scheduled Pickup Date:${formatNumericDate(scheduledPickupDate)}`)).toBeVisible();
    await expect(page.getByText(`Actual Pickup Date:${formatNumericDate(actualPickupDate)}`)).toBeVisible();
    await expect(page.getByText(`Scheduled Delivery Date:${formatNumericDate(scheduledDeliveryDate)}`)).toBeVisible();
    await expect(page.getByText(`Actual Delivery Date:${formatNumericDate(actualDeliveryDate)}`)).toBeVisible();
    await expect(page.getByText('Estimated Weight:7500')).toBeVisible();
    await expect(page.getByText('Actual Weight:8000')).toBeVisible();
    await expect(page.getByText('Destination Address:142 E Barrel Hoop Circle, Joshua Tree, CA 92252')).toBeVisible();

    // Can only create a payment request if there is a destination

    // waits for the create page to load
    await page.getByText('Create Payment Request').click();
    await expect(page.locator('input[name="serviceItems"]')).not.toHaveCount(0);

    expect(page.url()).toContain(`/simulator/moves/${moveID}/payment-requests/new`);

    // select all of the service items from the move and shipment
    //
    // UGH
    // because of the styling of this input item, we cannot use a
    // css locator for the input item and then click it
    //
    // The styling is very similar to the issue described in
    //
    // https://github.com/microsoft/playwright/issues/3688
    //
    const serviceItems = page.getByText('Add to payment request');
    const serviceItemCount = await serviceItems.count();
    expect(serviceItemCount).toBeGreaterThan(0);
    for (let i = 0; i < serviceItemCount; i += 1) {
      await serviceItems.nth(i).click();
    }

    await page.getByText('Submit Payment Request').click();

    // In CI in particular, this can take longer than 5 seconds
    await expect(page.getByText('Successfully created payment request')).toBeVisible({ timeout: 10000 });

    expect(page.url()).toContain(`/simulator/moves/${moveID}/details`);
    // could also check for a payment request number but we won't know the value ahead of time
  });
});
