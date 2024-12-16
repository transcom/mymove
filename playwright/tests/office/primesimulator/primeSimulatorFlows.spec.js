/**
 * Semi-automated converted from a cypress test, and thus may contain
 * non best-practices, in particular: heavy use of `page.locator`
 * instead of `page.getBy*`.
 */

// @ts-check
import { test, expect } from '../../utils/office/officeTest';

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
    await page.locator('#moveCode').fill(moveLocator);
    await page.locator('#moveCode').press('Enter');
    await page.getByTestId('moveCode-0').click();
    await officePage.waitForLoading();
    await expect(page.getByText(moveLocator)).toBeVisible();
    expect(page.url()).toContain(`/simulator/moves/${moveID}/details`);
    // waits for the move details page to load
    await expect(page.getByText('SUBMITTED')).toHaveCount(1);
    await page.getByRole('link', { name: 'Update Shipment', exact: true }).click();

    // waits for the update shipment page to load
    expect(page.url()).toContain(`/simulator/moves/${moveID}/shipments`);

    const { relativeDate: scheduledDeliveryDate, formattedDate: formattedScheduledDeliveryDate } =
      formatRelativeDate(11);
    await page.locator('input[name="scheduledDeliveryDate"]').fill(formattedScheduledDeliveryDate);
    await page.locator('input[name="scheduledDeliveryDate"]').blur();
    const { relativeDate: actualDeliveryDate, formattedDate: formattedActualDeliveryDate } = formatRelativeDate(12);
    await page.locator('input[name="actualDeliveryDate"]').fill(formattedActualDeliveryDate);
    await page.locator('input[name="actualDeliveryDate"]').blur();
    // there must be sufficient time prior to the pickup dates to update the estimated weight
    const { relativeDate: scheduledPickupDate, formattedDate: formattedScheduledPickupDate } = formatRelativeDate(11);
    await page.locator('input[name="scheduledPickupDate"]').fill(formattedScheduledPickupDate);
    await page.locator('input[name="scheduledPickupDate"]').blur();
    const { relativeDate: actualPickupDate, formattedDate: formattedActualPickupDate } = formatRelativeDate(12);
    await page.locator('input[name="actualPickupDate"]').fill(formattedActualPickupDate);
    await page.locator('input[name="actualPickupDate"]').blur();
    // update shipment does not require these fields but we need actual weight to create a payment request, we could
    // perform multiple updates.
    await page.locator('input[name="estimatedWeight"]').type('{backspace}7500');
    await page.locator('input[name="actualWeight"]').type('{backspace}8000');
    await page.locator('input[name="destinationAddress.streetAddress1"]').fill('142 E Barrel Hoop Circle');
    const locationLookup = 'JOSHUA TREE, CA 92252 (SAN BERNARDINO)';
    await page.locator('input#destinationAddress-location-input').fill('92252');
    await expect(page.getByText(locationLookup, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');

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
    await expect(page.getByText('Delivery Address:142 E Barrel Hoop Circle, Joshua Tree, CA 92252')).toBeVisible();
  });

  test('is able to create payment requests for shipment-level service items', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildPrimeSimulatorMoveNeedsShipmentUpdate();

    await officePage.signInAsNewPrimeSimulatorUser();
    const moveLocator = move.locator;
    const moveID = move.id;

    // wait for the the available moves page to load
    // select the move from the list
    await page.locator('#moveCode').fill(moveLocator);
    await page.locator('#moveCode').press('Enter');
    await page.getByTestId('moveCode-0').click();
    await officePage.waitForLoading();
    await expect(page.getByText(moveLocator)).toBeVisible();
    expect(page.url()).toContain(`/simulator/moves/${moveID}/details`);
    // waits for the move details page to load
    await expect(page.getByText('SUBMITTED')).toHaveCount(1);
    await page.getByRole('link', { name: 'Update Shipment', exact: true }).click();

    // waits for the update shipment page to load
    expect(page.url()).toContain(`/simulator/moves/${moveID}/shipments`);

    const { relativeDate: scheduledDeliveryDate, formattedDate: formattedScheduledDeliveryDate } =
      formatRelativeDate(11);
    await page.locator('input[name="scheduledDeliveryDate"]').fill(formattedScheduledDeliveryDate);
    await page.locator('input[name="scheduledDeliveryDate"]').blur();
    const { relativeDate: actualDeliveryDate, formattedDate: formattedActualDeliveryDate } = formatRelativeDate(12);
    await page.locator('input[name="actualDeliveryDate"]').fill(formattedActualDeliveryDate);
    await page.locator('input[name="actualDeliveryDate"]').blur();
    // there must be sufficient time prior to the pickup dates to update the estimated weight
    const { relativeDate: scheduledPickupDate, formattedDate: formattedScheduledPickupDate } = formatRelativeDate(11);
    await page.locator('input[name="scheduledPickupDate"]').fill(formattedScheduledPickupDate);
    await page.locator('input[name="scheduledPickupDate"]').blur();
    const { relativeDate: actualPickupDate, formattedDate: formattedActualPickupDate } = formatRelativeDate(12);
    await page.locator('input[name="actualPickupDate"]').fill(formattedActualPickupDate);
    await page.locator('input[name="actualPickupDate"]').blur();
    // update shipment does not require these fields but we need actual weight to create a payment request, we could
    // perform multiple updates.
    await page.locator('input[name="estimatedWeight"]').type('{backspace}7500');
    await page.locator('input[name="actualWeight"]').type('{backspace}8000');
    await page.locator('input[name="destinationAddress.streetAddress1"]').fill('142 E Barrel Hoop Circle');
    const locationLookup = 'JOSHUA TREE, CA 92252 (SAN BERNARDINO)';
    await page.locator('input#destinationAddress-location-input').fill('92252');
    await expect(page.getByText(locationLookup, { exact: true })).toBeVisible();
    await page.keyboard.press('Enter');
    await page.locator('select[name="destinationType"]').selectOption({ label: 'Home of record (HOR)' });

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
    await expect(page.getByText('Delivery Address:142 E Barrel Hoop Circle, Joshua Tree, CA 92252')).toBeVisible();

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
      if (
        (await serviceItems.nth(i).locator('..').locator('..').locator('h3').textContent()).includes(
          'Domestic linehaul',
        ) ||
        (await serviceItems.nth(i).locator('..').locator('..').locator('h3').textContent()).includes('Fuel surcharge')
      ) {
        await serviceItems.nth(i).click();
      }
    }

    await page.getByText('Submit Payment Request').click();

    // In CI in particular, this can take longer than 5 seconds
    await expect(page.getByText('Successfully created payment request')).toBeVisible({ timeout: 10000 });

    expect(page.url()).toContain(`/simulator/moves/${moveID}/details`);
    // could also check for a payment request number but we won't know the value ahead of time
  });

  // TODO: Unable to get a shipment to show up for Prime, skipping for now.
  test.skip('is able to see partial status of partial PPM move', async ({ officePage }) => {
    const partialPpmMoveCloseout = await officePage.testHarness.buildPartialPPMMoveReadyForCloseout();
    const partialPpmCloseoutLocator = partialPpmMoveCloseout.locator;

    await officePage.signInAsNewPrimeSimulatorUser();
    await officePage.primeSimulatorNavigateToMove(partialPpmCloseoutLocator);
  });

  test('is able to submit a SIT extension request', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildHHGMoveInSIT();

    await officePage.signInAsNewPrimeSimulatorUser();
    const moveLocator = move.locator;
    const moveID = move.id;

    await page.locator('#moveCode').fill(moveLocator);
    await page.locator('#moveCode').press('Enter');
    await page.getByTestId('moveCode-0').click();
    await officePage.waitForLoading();
    await expect(page.getByText(moveLocator)).toBeVisible();
    expect(page.url()).toContain(`/simulator/moves/${moveID}/details`);

    // Go to Request SIT extension page
    await page.getByText('Request SIT Extension').click();
    expect(page.url()).toContain(`/simulator/moves/${moveID}/shipments`);

    // Check labels and fill out the form
    await expect(page.getByText(`Request Reason`)).toBeVisible();
    await page.locator('select[name="requestReason"]');

    // Expected values for Request Reason
    const expectedValues = [
      'SERIOUS_ILLNESS_MEMBER',
      'SERIOUS_ILLNESS_DEPENDENT',
      'IMPENDING_ASSIGNEMENT',
      'DIRECTED_TEMPORARY_DUTY',
      'NONAVAILABILITY_OF_CIVILIAN_HOUSING',
      'AWAITING_COMPLETION_OF_RESIDENCE',
      'OTHER',
    ];

    // Check each option
    for (const option of expectedValues) {
      await page.locator('select[name="requestReason"]').selectOption({ value: option });
    }

    await expect(page.getByText(`Requested Days`)).toBeVisible();
    await page.locator('input[name="requestedDays"]').fill('12');
    await expect(page.getByText(`Contractor Remarks`)).toBeVisible();
    await page.locator('textarea[name="contractorRemarks"]').fill('Testing contractor remarks');

    // Submit the form
    await page.getByText('Request SIT Extension').click();

    // Get success message
    await expect(page.getByText('Successfully created SIT extension request')).toBeVisible({ timeout: 10000 });
    expect(page.url()).toContain(`/simulator/moves/${moveID}/details`);
  });

  test('is able submit payment request on SIT without destination SIT Out Date', async ({ page, officePage }) => {
    const move = await officePage.testHarness.buildHHGMoveInSITNoDestinationSITOutDate();
    const moveLocator = move.locator;
    const moveID = move.id;
    const items = move.MTOServiceItems;
    const weight = '500';
    let serviceItemID;

    await officePage.signInAsNewPrimeSimulatorUser();
    await page.locator('#moveCode').fill(moveLocator);
    await page.locator('#moveCode').press('Enter');
    await page.getByTestId('moveCode-0').click();
    await page.getByRole('link', { name: 'Create Payment Request' }).click();

    const serviceItemCount = items.length;
    expect(serviceItemCount).toBeGreaterThan(0);
    for (let i = 0; i < serviceItemCount; i += 1) {
      const dddsitIt = items.find((items) => items.ReService.code === 'DDDSIT');
      serviceItemID = dddsitIt.ID;
    }

    await page.locator(`[id="${serviceItemID}-div"] > .usa-checkbox`).click();
    await page.locator(`input[name="params\\.${serviceItemID}.WeightBilled"]`).fill(weight);
    await page.getByText('Submit Payment Request').click();
    await expect(page.getByText('Successfully created payment request')).toBeVisible({ timeout: 10000 });

    await page.getByRole('link', { name: 'Create Payment Request' }).click();
    for (let i = 0; i < serviceItemCount; i += 1) {
      const ddsfsc = items.find((items) => items.ReService.code === 'DDSFSC');
      serviceItemID = ddsfsc.ID;
    }

    await page.locator(`[id="${serviceItemID}-div"] > .usa-checkbox`).click();
    await page.locator(`input[name="params\\.${serviceItemID}.WeightBilled"]`).fill(weight);
    await page.getByText('Submit Payment Request').click();
    await expect(page.getByText('Successfully created payment request')).toBeVisible({ timeout: 10000 });

    expect(page.url()).toContain(`/simulator/moves/${moveID}/details`);
  });
});
