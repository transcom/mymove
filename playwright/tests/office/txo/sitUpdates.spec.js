import { test, expect } from '../../utils/office/officeTest';

import { TooFlowPage } from './tooTestFixture';

test.describe('TOO user', () => {
  /** @type {TooFlowPage} */
  let tooFlowPage;

  test.describe('updating a move shipment in SIT', () => {
    test.beforeEach(async ({ officePage }) => {
      // build move in SIT with 90 days authorized and without pending extension requests
      const move = await officePage.testHarness.buildHHGMoveInSIT();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await officePage.tooNavigateToMove(move.locator);
    });

    test('is able to increase a SIT authorization', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // increase SIT authorization to 100 days
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Edit SIT authorization' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('100');
      await page.getByTestId('dropdown').selectOption('AWAITING_COMPLETION_OF_RESIDENCE');
      await page.getByTestId('officeRemarks').fill('residence under construction');
      await expect(page.getByTestId('form').getByTestId('button')).toBeEnabled();
      await page.getByTestId('form').getByTestId('button').click();

      // assert that days authorization is now 100
      await expect(page.getByTestId('sitStatusTable').getByText('100', { exact: true }).first()).toBeVisible();
    });

    test('is able to decrease a SIT authorization', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // decrease SIT authorization to 80 days
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Edit SIT authorization' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('80');
      await page.getByTestId('dropdown').selectOption('AWAITING_COMPLETION_OF_RESIDENCE');
      await page.getByTestId('officeRemarks').fill('residence under construction');
      await page.getByTestId('form').getByTestId('button').click();

      // assert that days authorization is now 80
      await expect(page.getByTestId('sitStatusTable').getByText('80', { exact: true }).first()).toBeVisible();
    });

    test('is unable to decrease the SIT authorization below the number of days already used', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // try to decrease SIT authorization to 1 day
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Edit SIT authorization' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('1');
      await page.getByTestId('dropdown').selectOption('AWAITING_COMPLETION_OF_RESIDENCE');
      await page.getByTestId('officeRemarks').fill('residence under construction');

      // assert that save button is disabled and error messages are present
      await expect(page.getByTestId('form').getByTestId('button')).not.toBeEnabled();
      await expect(page.getByTestId('form').getByTestId('sitStatusTable').getByTestId('errorMessage')).toBeVisible();
      await expect(
        page.getByText('The end date must occur after the start date. Please select a new date.', { exact: true }),
      ).toBeVisible();
    });
  });

  test.describe('updating a move shipment in SIT with a SIT extension request', () => {
    test.beforeEach(async ({ officePage }) => {
      // build move in SIT with 90 days authorized and with one pending extension request
      const move = await officePage.testHarness.buildHHGMoveInSITWithPendingExtension();
      await officePage.signInAsNewTOOUser();
      tooFlowPage = new TooFlowPage(officePage, move);
      await officePage.tooNavigateToMove(move.locator);
    });

    test('is able to approve the SIT extension request', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // assert that there is a pending SIT extension request
      await expect(page.getByText('Additional days requested')).toBeVisible();

      // approve SIT extension with an adjusted approved days value of 100 days and change the extension reason
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Review additional days requested' })).toBeVisible();
      await page.getByTestId('daysApproved').clear();
      await page.getByTestId('daysApproved').fill('100');
      await page.getByText('Yes', { exact: true }).click();
      await page.getByTestId('dropdown').selectOption('OTHER');
      await page.getByTestId('officeRemarks').fill('allowance increased by 20 days instead of the requested 45 days');
      await page.getByTestId('form').getByTestId('button').click();

      // assert that there is no pending SIT extension request and the days authorization is now 100
      await expect(page.getByText('Additional days requested')).toBeHidden();
      await expect(page.getByTestId('sitStatusTable').getByText('100', { exact: true }).first()).toBeVisible();
    });

    test('is able to deny the SIT extension request', async ({ page }) => {
      // navigate to MTO tab
      await page.getByTestId('MoveTaskOrder-Tab').click();
      await tooFlowPage.waitForPage.moveTaskOrder();

      // assert that there is a pending SIT extension request
      await expect(page.getByText('Additional days requested')).toBeVisible();

      // deny SIT extension
      await page.getByTestId('sitExtensions').getByTestId('button').click();
      await expect(page.getByRole('heading', { name: 'Review additional days requested' })).toBeVisible();
      await page.getByText('No', { exact: true }).click();
      await page.getByTestId('officeRemarks').fill('extension request denied');
      await page.getByTestId('form').getByTestId('button').click();

      // assert that there is no pending SIT extension request and the days authorization is still 90
      await expect(page.getByText('Additional days requested')).toBeHidden();
      await expect(page.getByTestId('sitStatusTable').getByText('90', { exact: true }).first()).toBeVisible();
    });
  });
});

test.describe('During ZIP update test', () => {
  test('the TIO was able to verify the correct ZIP', async ({ officePage }) => {
    const move = await officePage.testHarness.buildHHGMoveInSIT();

    // TOO Approves Address Update in updated SIT item
    await officePage.signInAsNewMultiRoleUser();
    await officePage.page.getByText('Change user role').click();
    await officePage.page.getByRole('button').getByText('transportation_ordering_officer').click();
    await officePage.page.getByRole('heading').getByText('All moves').waitFor();
    await officePage.page.locator('[name="locator"]').fill(move.locator);
    await officePage.page.getByRole('heading').getByText('All moves').click();
    await officePage.page.locator('[data-testid^="locator-"]').click();
    await officePage.page.getByRole('heading').getByText('Move details').waitFor();
    await officePage.page.getByTestId('MoveTaskOrder-Tab').click();
    await officePage.page.getByRole('heading').getByText('Move task order').waitFor();
    await officePage.page.locator('[data-testid="reviewRequestTextButton"]').click();
    await officePage.page.getByTestId('modal').getByText('Review request: service item update').waitFor();
    const addressLines = await officePage.page.$$('[data-testid="AddressLine"]');
    expect(addressLines.length === 6).toBeTruthy();
    const zipNew = await addressLines[5].textContent();
    await officePage.page.getByTestId('reviewSITAddressUpdateForm').getByTestId('radio').getByText('Yes').click();
    await officePage.page.getByTestId('reviewSITAddressUpdateForm').getByTestId('radio').getByText('Yes').blur();
    await officePage.page.getByTestId('officeRemarks').fill('approved the address change');
    expect(await officePage.page.getByRole('button').getByText('Save').isEnabled());
    await officePage.page.getByRole('button').getByText('Save').click();
    const modalClostButton = await officePage.page.getByTestId('modal').getByTestId('modalCloseButton');
    if (await modalClostButton) {
      await modalClostButton.click();
    }

    // Prime user submits payment request
    await officePage.page.getByText('Change user role').click();
    await officePage.page.getByRole('button').getByText('prime_simulator').click();
    await officePage.page.getByRole('heading').getByText('Moves available to Prime').waitFor();
    await officePage.page.locator('[name="moveCode"]').fill(move.locator);
    await officePage.page.getByRole('heading').getByText('Moves available to Prime').click();
    await officePage.page.locator('[data-testid^="moveCode-"]').click();
    await officePage.page
      .getByTestId('move-details-flash-grid-container')
      .getByText('Create Payment Request')
      .waitFor();
    await officePage.page.getByTestId('move-details-flash-grid-container').getByText('Create Payment Request').click();
    await officePage.page.getByRole('heading').getByText('Domestic destination SIT fuel surcharge').waitFor();
    await officePage.page.waitForLoadState('load');
    const paymentClickers = await officePage.page.getByText('Add to payment request').all();
    for (const clicker of paymentClickers) {
      const h3TextContent = await clicker
        .locator('..')
        .locator('..')
        .locator('[class="descriptionList_descriptionList__v+wtB"] h3')
        .textContent();
      if (
        h3TextContent.includes('Domestic destination SIT fuel surcharge') ||
        h3TextContent.includes('Domestic destination SIT delivery')
      ) {
        await clicker.click();
        await clicker.blur();
      }
    }
    await officePage.page.getByRole('button').getByText('Submit Payment Request').click();

    // TIO user views the request
    await officePage.page.getByText('Change user role').click();
    await officePage.page.getByRole('button').getByText('transportation_invoicing_officer').click();
    await officePage.page.waitForLoadState('load');
    await officePage.page.locator('[name="locator"]').fill(move.locator);
    await officePage.page.getByRole('heading').getByText('Payment requests').click();
    if (await officePage.page.getByTestId('locator-0').isVisible()) {
      await officePage.page.getByTestId('locator-0').click();
    } else {
      await officePage.page.getByText('Change user role').click();
      await officePage.page.getByRole('button').getByText('transportation_invoicing_officer').click();
      await officePage.page.waitForLoadState('load');
      await officePage.page.locator('[name="locator"]').fill(move.locator);
      await officePage.page.getByRole('heading').getByText('Payment requests').click();
      await officePage.page.getByTestId('locator-0').click();
    }
    await officePage.page.locator(`a[href="/moves/${move.locator}/payment-requests"]`).click();
    const reviewButton = await officePage.page.getByTestId('reviewBtn');
    if (await reviewButton.isDisabled()) {
      const reviewWeightsButton = await officePage.page
        .locator('button[class="usa-button"]')
        .getByText('Review weights');
      await reviewWeightsButton.click();
      await officePage.page.getByRole('button').getByText('Review shipment weights').click();
      const editButton = await officePage.page.getByRole('button').getByText('Edit');
      await editButton.click();
      const estimatedWeight = await officePage.page
        .locator('div[class="ShipmentCard_weights__TobbB"]')
        .getByTestId('estimatedWeightContainer')
        .locator('strong')
        .textContent();
      await officePage.page
        .locator('fieldset[class="usa-fieldset EditBillableWeight_fieldset__VATbs"]')
        .locator('div[class="usa-form-group"]')
        .locator('input[class="usa-input EditBillableWeight_maxBillableWeight__hV+1b"]')
        .fill(estimatedWeight);
      const shipmentWeightsRemarks = await officePage.page
        .locator('fieldset[class="usa-fieldset EditBillableWeight_fieldset__VATbs"]')
        .locator('div')
        .locator('textarea[data-testid="remarks"]');
      await shipmentWeightsRemarks.fill('test remarks about weight');
      await shipmentWeightsRemarks.blur();
      await officePage.page.getByTestId('button').getByText('Save changes').click();
      await officePage.page.locator(`a[href="/moves/${move.locator}/payment-requests"]`).click();
    }
    await reviewButton.click();
    await officePage.page.locator('h2').getByText('Review service items').waitFor();
    const itemNumHeader = await officePage.page
      .locator('div[data-testid="ReviewServiceItems"]')
      .locator('div[class="ReviewServiceItems_top__g1U7J"]')
      .locator('div[data-testid="itemCount"]')
      .textContent();
    const itemNumMax = parseInt(itemNumHeader.substring(1), 10);
    await Promise.all(
      Array.from({ length: itemNumMax }, async () => {
        const serviceItemName = await officePage.page
          .getByTestId('ShipmentContainer')
          .locator('dd[data-testid="serviceItemName"]')
          .textContent();
        if (serviceItemName.includes('Domestic destination SIT')) {
          await officePage.page
            .getByTestId('ShipmentContainer')
            .locator('button[class="usa-button usa-button--unstyled ServiceItemCard_toggleCalculations__D54xv"]')
            .getByText('Show calculations')
            .click();
          const serviceItemZipSectionString = await officePage.page
            .locator('div[data-testid="ServiceItemCalculations"]')
            .locator('div[class="ServiceItemCalculations_col__MePhS"]')
            .locator('div')
            .locator('ul[data-testid="details"]')
            .locator('li')
            .locator('p')
            .locator('small')
            .textContent();
          const serviceItemZips = {
            from: await serviceItemZipSectionString.substring(0, 5),
            to: await serviceItemZipSectionString.substring(serviceItemZipSectionString.length - 6),
          };
          expect(serviceItemZips.to).toEqual(zipNew);
        }
        await officePage.page.getByTestId('ShipmentContainer').getByTestId('approveRadio').check();
        await officePage.page.getByTestId('ShipmentContainer').getByTestId('nextServiceItem').click();
      }),
    );
  });
});
