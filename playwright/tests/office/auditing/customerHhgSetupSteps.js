// @ts-nocheck
import { milmoveUserType } from '../../utils/my/customerTest';
import './auditUtils';
import { dateInputOperator } from './auditUtils';
import { expect } from './auditTestSetup';
import _ from 'lodash';

const GetUserMoveData = async (testHarness) => {
  const move = await testHarness.buildMoveWithOrders();
  const userId = move.Orders.ServiceMember.user_id;
  return { move, userId };
};

/**
 * Create a Household Good Shipment as a move owner (service member/customer)
 */
const testExecution = async ({
  page,
  helpers,
  baseURLS,
  testHarness,
  signInAsExisting,
  signInAsNew,
  signOut,
  waitForLoading,
  stringHelpers,
  utils,
  fast: { clickTextAsync, typeInto, typeAndBlur, selectValue },
}) => {
  const fillHHGForm = async (page, { formatRelativeDate }, moveId) => {
    const hhgFormData = {
      preferredPickupDate: formatRelativeDate(5),
      preferredDeliveryDate: formatRelativeDate(20),
      useExistingAddress: true,
      releasingAgentDetails: {
        firstName: 'larry',
        lastName: 'thombson',
        phoneNumber: '555-555-5555',
        email: 'lthompson@test.test',
      },
      receivingAgentDetails: {
        firstName: 'larrisson',
        lastName: 'thornbsonson',
        phoneNumber: '555-555-5555',
        email: 'lthornbsonson@test.test',
      },
      destination: {
        street: '142 E Barrel Hoop Circle',
        city: 'Joshua Tree',
        state: 'CA',
        postalCode: '92252',
      },
    };

    const REQUESTED_PICKUP_DATE_SELECTOR = '#requestedPickupDate';
    const pickupDateElement = page.locator(REQUESTED_PICKUP_DATE_SELECTOR);
    //await expect(pickupDateElement).toBeVisible();

    const pickupDate = new Date(hhgFormData.preferredPickupDate.relativeDate);
    await dateInputOperator(page, pickupDateElement, pickupDate);

    //click 'use my current address'
    await page.getByTestId('checkbox').click();

    //submit releasing agent fields
    await typeInto('input[name="pickup.agent.firstName"]', hhgFormData.releasingAgentDetails.firstName);
    await typeInto('input[name="pickup.agent.lastName"]', hhgFormData.releasingAgentDetails.lastName);
    await typeInto('input[name="pickup.agent.phone"]', hhgFormData.releasingAgentDetails.phoneNumber);
    await typeInto('input[name="pickup.agent.email"]', hhgFormData.releasingAgentDetails.email);
    //submit delivery date
    
    const REQUESTED_DELIVERY_SELECTOR = '#requestedDeliveryDate';
    const deliveryDateElement = page.locator(REQUESTED_DELIVERY_SELECTOR);
    //await expect(pickupDateElement).toBeVisible();

    const deliveryDate = new Date(hhgFormData.preferredDeliveryDate.relativeDate);
    await dateInputOperator(page, deliveryDateElement, deliveryDate);

    //submit receiving agent fields
    await typeInto('input[name="delivery.agent.firstName"]', hhgFormData.receivingAgentDetails.firstName);
    await typeInto('input[name="delivery.agent.lastName"]', hhgFormData.receivingAgentDetails.lastName);
    await typeInto('input[name="delivery.agent.phone"]', hhgFormData.receivingAgentDetails.phoneNumber);
    await typeInto('input[name="delivery.agent.email"]', hhgFormData.receivingAgentDetails.email);

    const REVIEW_PATH = stringHelpers.format('moves/{0}/review', moveId);
    const reviewUrlCheck = ({ pathname }) => stringHelpers.trimOuterSymbols(pathname, '/') === REVIEW_PATH;

    const waitForReviewPath = page.waitForURL(reviewUrlCheck);
    const NEXT_WIZARD_BUTTON = 'wizardNextButton';
    await page.getByTestId(NEXT_WIZARD_BUTTON).click();
    await waitForReviewPath;
  };

  const navigateToHHGForm = async (page, moveId) => {
    const MOVING_INFO_PATH = stringHelpers.format('moves/{0}/moving-info', moveId);
    const movingInfoUrlCheck = ({ pathname }) => stringHelpers.trimOuterSymbols(pathname, '/') === MOVING_INFO_PATH;

    const SHIPMENT_TYPE_PATH = stringHelpers.format('moves/{0}/shipment-type', moveId);
    const shipmentTypeUrlCheck = ({ pathname }) => stringHelpers.trimOuterSymbols(pathname, '/') === SHIPMENT_TYPE_PATH;

    const NEW_HHG_SHIPMENT_PATH = stringHelpers.format('moves/{0}/new-shipment', moveId);
    const hhgNewShipmentUrlCheck = ({ pathname, searchParams }) => {
      const queryObject = Object.fromEntries(searchParams.entries());
      const correctUrl = stringHelpers.trimOuterSymbols(pathname, '/') === NEW_HHG_SHIPMENT_PATH;
      const correctQueryValue = _.isEqual(queryObject, { type: 'HHG' });
      return correctUrl && correctQueryValue;
    };

    const waitForMovingInfoPath = page.waitForURL(movingInfoUrlCheck);
    const SETUP_SHIPMENT_BUTTON = 'shipment-selection-btn';
    await page.getByTestId(SETUP_SHIPMENT_BUTTON).click();
    await waitForMovingInfoPath;

    const waitForShipmentTypePath = page.waitForURL(shipmentTypeUrlCheck);
    const NEXT_WIZARD_BUTTON = 'wizardNextButton';
    await page.getByTestId(NEXT_WIZARD_BUTTON).click();
    await waitForShipmentTypePath;

    // select HHG radio button
    await page.locator('label[for="HHG"]').click();

    const waitForNewHHGShipmentPath = page.waitForURL(hhgNewShipmentUrlCheck);
    await page.getByTestId(NEXT_WIZARD_BUTTON).click();
    await waitForNewHHGShipmentPath;
  };

  // ------[START OF STEPS]------ //

  const { move, userId } = await GetUserMoveData(testHarness);

  await page.goto(baseURLS.my);
  await signInAsExisting(userId);
  //start seting up the shipments...
  await navigateToHHGForm(page, move.id);
  await fillHHGForm(page, utils, move.id);

  //wait for url when clicking next to proceed...
  const AGREEMENT_PATH = stringHelpers.format('moves/{0}/agreement', move.id);
  const agreementUrlCheck = ({ pathname }) => stringHelpers.trimOuterSymbols(pathname, '/') === AGREEMENT_PATH;
  
  //wait for url when clicking next to proceed...
  const waitForAgreementPath = page.waitForURL(agreementUrlCheck);
  const NEXT_WIZARD_BUTTON = 'wizardNextButton';
  await page.getByTestId(NEXT_WIZARD_BUTTON).click();
  await waitForAgreementPath;

  //sign agreement
  await typeInto('input[name="signature"]', 'first last signature');

  //wait for url when clicking complete to proceed...
  const waitForCustomerScreen = page.waitForURL(baseURLS.my);
  await page.getByTestId('wizardCompleteButton').click();
  await waitForCustomerScreen;

  // ------[END OF STEPS]------ //
};

export const RUN_MOVER_OWNER_SETS_UP_HHG = {
  run: (input, getHelpers) => testExecution({ ...input, ...getHelpers(input) }),
};
