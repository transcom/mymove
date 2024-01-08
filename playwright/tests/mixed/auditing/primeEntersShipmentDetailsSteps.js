import { expect } from './auditTestSetup';

export const testExecution = async ({
  page,
  helpers,
  baseURLS,
  testHarness,
  waitForLoading,
  signInAsNewPrimeSimulatorUser,
  signOut,
  stringHelpers,
  utils,
  fast: { clickTextAsync, typeInto, typeAndBlur, selectValue },
}) => {
  const FillDestinationDetailFields = async ({
    estimatedWeight,
    actualWeight,
    address: { street, city, state, postalCode },
  }) => {
    const selectors = {
      estimatedWeight: 'input[name="estimatedWeight"]',
      actualWeight: 'input[name="actualWeight"]',
      address: {
        street: 'input[name="destinationAddress.streetAddress1"]',
        city: 'input[name="destinationAddress.city"]',
        state: 'select[name="destinationAddress.state"]',
        postalCode: 'input[name="destinationAddress.postalCode"]',
      },
    };

    await typeInto(selectors.estimatedWeight, `{backspace}${estimatedWeight}`);
    await typeInto(selectors.actualWeight, `{backspace}${actualWeight}`);
    await typeInto(selectors.address.street, `${street}`);
    await typeInto(selectors.address.city, `${city}`);
    await selectValue(selectors.address.state, { label: state });
    await typeInto(selectors.address.postalCode, `${postalCode}`);
  };

  const checkShipmentWeights = async ({
    scheduledPickupDate: sPickupDate,
    actualPickupDate: aPickupDate,
    scheduledDeliveryDate: sDeliveryDate,
    actualDeliveryDate: aDeliveryDate,
    estimatedWeight,
    actualWeight,
    address: { street, city, state, postalCode },
  }) => {
    const SCHEDULED_PICKUP_DATE_LABEL = 'Scheduled Pickup Date:';
    const SCHEDULED_PICKUP_DATE_VALUE = utils.formatNumericDate(sPickupDate);
    const SCHEDULED_PICKUP_DATE_CONTENT = `${SCHEDULED_PICKUP_DATE_LABEL}${SCHEDULED_PICKUP_DATE_VALUE}`;

    const ACTUAL_PICKUP_DATE_LABEL = 'Actual Pickup Date:';
    const ACTUAL_PICKUP_DATE_VALUE = utils.formatNumericDate(aPickupDate);
    const ACTUAL_PICKUP_DATE_CONTENT = `${ACTUAL_PICKUP_DATE_LABEL}${ACTUAL_PICKUP_DATE_VALUE}`;

    const SCHEDULED_DELIVERY_DATE_LABEL = 'Scheduled Delivery Date:';
    const SCHEDULED_DELIVERY_DATE_VALUE = utils.formatNumericDate(sDeliveryDate);
    const SCHEDULED_DELIVERY_DATE_CONTENT = `${SCHEDULED_DELIVERY_DATE_LABEL}${SCHEDULED_DELIVERY_DATE_VALUE}`;

    const ACTUAL_DELIVERY_DATE_LABEL = 'Actual Delivery Date:';
    const ACTUAL_DELIVERY_DATE_VALUE = utils.formatNumericDate(aDeliveryDate);
    const ACTUAL_DELIVERY_DATE_CONTENT = `${ACTUAL_DELIVERY_DATE_LABEL}${ACTUAL_DELIVERY_DATE_VALUE}`;

    const ESTIMATED_WEIGHT_LABEL = 'Estimated Weight:';
    const ESTIMATED_WEIGHT_VALUE = estimatedWeight;
    const ESTIMATED_WEIGHT_CONTENT = utils.textWithNoTrailingNumbers(
      `${ESTIMATED_WEIGHT_LABEL}${ESTIMATED_WEIGHT_VALUE}`,
    );

    const ACTUAL_WEIGHT_LABEL = 'Actual Weight:';
    const ACTUAL_WEIGHT_VALUE = actualWeight;
    const ACTUAL_WEIGHT_CONTENT = utils.textWithNoTrailingNumbers(`${ACTUAL_WEIGHT_LABEL}${ACTUAL_WEIGHT_VALUE}`);

    const DESTINATION_ADDRESS_LABEL = 'Destination Address:';
    const DESTINATION_ADDRESS_VALUE = `${street}, ${city}, ${state} ${postalCode}`;
    const DESTINATION_ADDRESS_CONTENT = `${DESTINATION_ADDRESS_LABEL}${DESTINATION_ADDRESS_VALUE}`;

    const textToCheck = [
      SCHEDULED_PICKUP_DATE_CONTENT,
      ACTUAL_PICKUP_DATE_CONTENT,
      SCHEDULED_DELIVERY_DATE_CONTENT,
      ACTUAL_DELIVERY_DATE_CONTENT,
      ESTIMATED_WEIGHT_CONTENT,
      ACTUAL_WEIGHT_CONTENT,
      DESTINATION_ADDRESS_CONTENT,
    ].map(async (content) => expect(page.getByText(content)).toBeVisible());

    await textToCheck.reduce((chain, waiter) => chain.then(() => waiter), Promise.resolve());
  };

  const MOVE_CODE_COLUMN_SELECTOR = 'moveCode-0';

  /**
   * select the move from the list
   * wait for the the available moves page to load
   * waits for the move details page to load
   * */
  const NavigateToMoveDetailsScreen = async (move_selector) => {
    await page.locator(move_selector).fill(moveLocator);
    await page.locator(move_selector).press('Enter');

    // TODO TestIds indicate a dimension of significance.
    // naming the attribute to something with respect to the data might be of use.
    // having a centralized approach to deriving attributes across the presentation layer
    // would help when it comes to element associations
    await page.getByTestId(MOVE_CODE_COLUMN_SELECTOR).click();
    await waitForLoading();
    await expect(page.getByText(moveLocator)).toBeVisible();
  };

  const DAY_11TH = utils.formatRelativeDate(11);
  const DAY_12TH = utils.formatRelativeDate(12);

  const [
    { relativeDate: scheduledDeliveryDate, formattedDate: formattedScheduledDeliveryDate },
    { relativeDate: actualDeliveryDate, formattedDate: formattedActualDeliveryDate },
    { relativeDate: scheduledPickupDate, formattedDate: formattedScheduledPickupDate },
    { relativeDate: actualPickupDate, formattedDate: formattedActualPickupDate },
  ] = [DAY_11TH, DAY_12TH, DAY_11TH, DAY_12TH];

  const SCHEDULED_DELIVERY_DATE_INPUT_SELECTOR = 'input[name="scheduledDeliveryDate"]';
  const ACTUAL_DELIVERY_DATE_INPUT_SELECTOR = 'input[name="actualDeliveryDate"]';
  const SCHEDULED_PICKUP_DATE_INPUT_SELECTOR = 'input[name="scheduledPickupDate"]';
  const ACTUAL_PICKUP_DATE_INPUT_SELECTOR = 'input[name="actualPickupDate"]';

  const EXPECT_URL_MOVE_DETAILS_FORMAT = '/simulator/moves/{0}/details';
  const EXPECT_URL_SHIPMENTS_VALUE_FORMAT = '/simulator/moves/{0}/shipments';

  const INPUT_ACTION_ENTRIES = [
    { selector: SCHEDULED_DELIVERY_DATE_INPUT_SELECTOR, value: formattedScheduledDeliveryDate },
    { selector: ACTUAL_DELIVERY_DATE_INPUT_SELECTOR, value: formattedActualDeliveryDate },
    { selector: SCHEDULED_PICKUP_DATE_INPUT_SELECTOR, value: formattedScheduledPickupDate },
    { selector: ACTUAL_PICKUP_DATE_INPUT_SELECTOR, value: formattedActualPickupDate },
  ];

  const MOVE_CODE_SELECTOR = '#moveCode';

  const destinationFieldValues = {
    estimatedWeight: 7500,
    actualWeight: 8000,
    address: {
      street: '142 E Barrel Hoop Circle',
      city: 'Joshua Tree',
      state: 'CA',
      postalCode: '92252',
    },
  };

  //validate sit values against the UI
  const sitValues = {
    ...destinationFieldValues,
    scheduledDeliveryDate,
    actualDeliveryDate,
    scheduledPickupDate,
    actualPickupDate,
  };

  // --------------[BEGIN TEST]---------------- //

  const move = await testHarness.buildPrimeSimulatorMoveNeedsShipmentUpdate();
  const moveLocator = move.locator;
  const moveID = move.id;
  const EXPECT_MOVE_DETAILS_URL_VALUE = stringHelpers.format(EXPECT_URL_MOVE_DETAILS_FORMAT, moveID);
  const EXPECT_SHIPMENT_URL_VALUE = stringHelpers.format(EXPECT_URL_SHIPMENTS_VALUE_FORMAT, moveID);

  await page.goto(baseURLS.office);

  await signInAsNewPrimeSimulatorUser();
  await NavigateToMoveDetailsScreen(MOVE_CODE_SELECTOR);

  expect(page.url()).toContain(EXPECT_MOVE_DETAILS_URL_VALUE);

  await page.getByText('Update Shipment').click();

  expect(page.url()).toContain(EXPECT_SHIPMENT_URL_VALUE);

  await INPUT_ACTION_ENTRIES.reduce(
    async (chain, { selector, value }) => chain.then(() => typeAndBlur(selector, value)),
    Promise.resolve(),
  );

  await FillDestinationDetailFields(destinationFieldValues);

  const SAVE_BUTTON_TEXT = 'Save';
  await clickTextAsync(SAVE_BUTTON_TEXT);

  const EXPECT_SUCCESSFUL_SHIPMENT_UPDATE = 'Successfully updated shipment';
  await expect(page.getByText(EXPECT_SUCCESSFUL_SHIPMENT_UPDATE)).toHaveCount(1);

  expect(page.url()).toContain(EXPECT_MOVE_DETAILS_URL_VALUE);

  await checkShipmentWeights(sitValues);

  await signOut();
  // --------------[END TEST]---------------- //
};

export const RUN_PRIME_SETS_UP_SHIPMENT = ({ run: (input, getHelpers) => testExecution({ ...input, ...getHelpers(input) }) });

