// @ts-nocheck
import './auditUtils';
import { formatRelativeDate, dateInputOperator, stringHelpers as str } from './auditUtils';
import { expect } from './auditTestSetup';
import _ from 'lodash';

const URL_PATHS = {
  ROOT_PATH: '/',
  SIGN_IN_PATH: 'sign-in',
  MOVING_INFO_PATH: 'moves/{0}/moving-info',
  SHIPMENT_TYPE_PATH: 'moves/{0}/shipment-type',
  NEW_HHG_SHIPMENT_PATH: 'moves/{0}/new-shipment',
  REVIEW_PATH: 'moves/{0}/review',
  AGREEMENT_PATH: 'moves/{0}/agreement',
};

const compareUrlPath =
  (check) =>
  ({ pathname }) =>
    str.trimOuterSymbols(pathname, '/') === check;

const RegularPathChecks = {
  ROOT_PATH_CHECK: ({ pathname }) => pathname === URL_PATHS.ROOT_PATH,
  SIGN_IN_PATH_CHECK: compareUrlPath(URL_PATHS.SIGN_IN_PATH),
};

const GetCustomerMoveUrlPaths = (moveId) => ({
  MOVING_INFO_PATH_CHECK: compareUrlPath(str.format(URL_PATHS.MOVING_INFO_PATH, moveId)),
  SHIPMENT_TYPE_PATH_CHECK: compareUrlPath(str.format(URL_PATHS.SHIPMENT_TYPE_PATH, moveId)),
  HHG_NEW_SHIPMENT_PATH_PARAMS_CHECK: ({ pathname, searchParams }) => {
    const checkUrl = str.format(URL_PATHS.NEW_HHG_SHIPMENT_PATH, moveId);
    const queryObject = Object.fromEntries(searchParams.entries());
    const correctUrl = str.trimOuterSymbols(pathname, '/') === checkUrl;
    const correctQueryValue = _.isEqual(queryObject, { type: 'HHG' });
    return correctUrl && correctQueryValue;
  },
  REVIEW_PATH_CHECK: compareUrlPath(str.format(URL_PATHS.REVIEW_PATH, moveId)),
  AGREEMENT_PATH_CHECK: compareUrlPath(str.format(URL_PATHS.AGREEMENT_PATH, moveId)),
});

const GetUserMoveData = async (testHarness) => {
  const move = await testHarness.buildMoveWithOrders();
  const userId = move.Orders.ServiceMember.user_id;
  return { move, userId };
};

const waitForSigninScreen = async ({ page, navToUrl, SIGN_IN_PATH_CHECK }) => {
  const urlWaiter = page.waitForURL(SIGN_IN_PATH_CHECK);
  await page.goto(navToUrl);
  await urlWaiter;
};

const signinAndWaitForCustomerHomeScreen = async ({ page, signInAsExisting, userId, ROOT_PATH_CHECK }) => {
  const urlWaiter = page.waitForURL(ROOT_PATH_CHECK);
  await signInAsExisting(userId);
  await urlWaiter;
};

/** Beginning Shipment Flow Wizard */
const selectShipment = async ({ page, MOVING_URL_CHECK }) => {
  const urlWaiter = page.waitForURL(MOVING_URL_CHECK);
  const SETUP_SHIPMENT_BUTTON = 'shipment-selection-btn';
  await page.getByTestId(SETUP_SHIPMENT_BUTTON).click();
  await urlWaiter;
};

/** Click Next */
const wizardNextToShipmentType = async ({ page, SHIPMENT_TYPE_PATH_CHECK }) => {
  const urlWaiter = page.waitForURL(SHIPMENT_TYPE_PATH_CHECK);
  const NEXT_WIZARD_BUTTON = 'wizardNextButton';
  await page.getByTestId(NEXT_WIZARD_BUTTON).click();
  await urlWaiter;
};

/** select HHG radio button */
const chooseHHG = async ({ page }) => await page.locator('label[for="HHG"]').click();

/** Click Next */
const wizardNextToNewHHGShipment = async ({ page, HHG_NEW_SHIPMENT_PATH_PARAMS_CHECK }) => {
  const urlWaiter = page.waitForURL(HHG_NEW_SHIPMENT_PATH_PARAMS_CHECK);
  const NEXT_WIZARD_BUTTON = 'wizardNextButton';
  await page.getByTestId(NEXT_WIZARD_BUTTON).click();
  await urlWaiter;
};

/** Click Next */
const wizardNextToAgreement = async ({ page, AGREEMENT_PATH_CHECK }) => {
  const urlWaiter = page.waitForURL(AGREEMENT_PATH_CHECK);
  const NEXT_WIZARD_BUTTON = 'wizardNextButton';
  await page.getByTestId(NEXT_WIZARD_BUTTON).click();
  await urlWaiter;
};

const choosePickupDate = async ({ page, dateToSelect, datePickerOperator }) => {
  const SELECTOR = '#requestedPickupDate';
  const dateElement = page.locator(SELECTOR);
  const selectedDate = new Date(dateToSelect);
  await datePickerOperator(page, dateElement, selectedDate);
};

const selectUseMyCurrentAddress = async ({ page }) => await page.getByTestId('checkbox').click();

const submitMainPickupAddress = async (
  { typeInto, selectValue },
  { streetAddress1, streetAddress2, streetAddress3, city, state, postalCode },
) => {
  await typeInto('input[name="pickup.address.streetAddress1"]', streetAddress1);

  if(streetAddress2)
    await typeInto('input[name="pickup.address.streetAddress2"]', streetAddress2);

  if(streetAddress3)
    await typeInto('input[name="pickup.address.streetAddress3"]', streetAddress3);

  await typeInto('input[name="pickup.address.city"]', city);
  await selectValue('select[name="pickup.address.state"]', { label: state });
  await typeInto('input[name="pickup.address.postalCode"]', postalCode);
};

const selectDontUseMyCurrentAddress = async ({ page }) => await page.getByTestId('checkbox').click();

const chooseHasSecondaryPickup = async ({ page }) => await page.locator('label[for="has-secondary-pickup"]').click();
const chooseNoSecondaryPickup = async ({ page }) => await page.locator('label[for="no-secondary-pickup"]').click();

const submitSecondaryPickupAddress = async (
  { typeInto, selectValue },
  { streetAddress1, streetAddress2, streetAddress3, city, state, postalCode },
) => {
  await typeInto('input[name="secondaryPickup.address.streetAddress1"]', streetAddress1);

  if(streetAddress2)
    await typeInto('input[name="secondaryPickup.address.streetAddress2"]', streetAddress2);
  
  if(streetAddress3)
    await typeInto('input[name="secondaryPickup.address.streetAddress3"]', streetAddress3);

  await typeInto('input[name="secondaryPickup.address.city"]', city);
  await selectValue('select[name="secondaryPickup.address.state"]', { label: state });
  await typeInto('input[name="secondaryPickup.address.postalCode"]', postalCode);
};

const submitReleaseAgentFields = async ({
  typeInto,
  releasingAgentDetails: { firstName, lastName, phoneNumber, email },
}) => {
  await typeInto('input[name="pickup.agent.firstName"]', firstName);
  await typeInto('input[name="pickup.agent.lastName"]', lastName);
  await typeInto('input[name="pickup.agent.phone"]', phoneNumber);
  await typeInto('input[name="pickup.agent.email"]', email);
};

const chooseDeliveryDate = async ({ page, dateToSelect, datePickerOperator }) => {
  const SELECTOR = '#requestedDeliveryDate';
  const dateElement = page.locator(SELECTOR);
  const selectedDate = new Date(dateToSelect);
  await datePickerOperator(page, dateElement, selectedDate);
};

const submitReceivingAgentFields = async ({
  typeInto,
  receivingAgentDetails: { firstName, lastName, phoneNumber, email },
}) => {
  await typeInto('input[name="delivery.agent.firstName"]', firstName);
  await typeInto('input[name="delivery.agent.lastName"]', lastName);
  await typeInto('input[name="delivery.agent.phone"]', phoneNumber);
  await typeInto('input[name="delivery.agent.email"]', email);
};

/** Click Next */
const wizardNextToReview = async ({ page, REVIEW_PATH_CHECK }) => {
  const urlWaiter = page.waitForURL(REVIEW_PATH_CHECK);
  const NEXT_WIZARD_BUTTON = 'wizardNextButton';
  await page.getByTestId(NEXT_WIZARD_BUTTON).click();
  await urlWaiter;
};

/** Fill in signature */
const submitSignature = async ({ signature, typeInto }) => await typeInto('input[name="signature"]', signature);

const wizardCompleteToCustomerHome = async ({ page, ROOT_PATH_CHECK }) => {
  const urlWaiter = page.waitForURL(ROOT_PATH_CHECK);
  const COMPLETE_WIZARD_BUTTON = 'wizardCompleteButton';
  await page.getByTestId(COMPLETE_WIZARD_BUTTON).click();
  await urlWaiter;
};

/** Customer HHG Setup Flow
 * Create a Household Good Shipment as a move owner (service member/customer)
 */
const PerformHHGSetup = async (
  { page, baseURLS, testHarness, signInAsExisting, fast: { typeInto, selectValue } },
  { preferredPickupDate, preferredDeliveryDate, releasingAgentDetails, receivingAgentDetails, signature },
) => {
  const { move, userId } = await GetUserMoveData(testHarness);

  const {
    ROOT_PATH_CHECK,
    MOVING_INFO_PATH_CHECK,
    SHIPMENT_TYPE_PATH_CHECK,
    HHG_NEW_SHIPMENT_PATH_PARAMS_CHECK,
    REVIEW_PATH_CHECK,
    AGREEMENT_PATH_CHECK,
    SIGN_IN_PATH_CHECK,
  } = { ...GetCustomerMoveUrlPaths(move.id), ...RegularPathChecks };

  await waitForSigninScreen({ page, navToUrl: baseURLS.my, SIGN_IN_PATH_CHECK });

  await signinAndWaitForCustomerHomeScreen({ page, signInAsExisting, userId, ROOT_PATH_CHECK });
  await selectShipment({ page, MOVING_INFO_PATH_CHECK });
  await wizardNextToShipmentType({ page, SHIPMENT_TYPE_PATH_CHECK });
  await chooseHHG({ page });
  await wizardNextToNewHHGShipment({ page, HHG_NEW_SHIPMENT_PATH_PARAMS_CHECK });
  //fill form start
  await choosePickupDate({
    page,
    dateToSelect: preferredPickupDate.relativeDate,
    datePickerOperator: dateInputOperator,
  });
  await selectUseMyCurrentAddress({ page });
  await chooseHasSecondaryPickup({page });

  await submitSecondaryPickupAddress({ typeInto, selectValue }, {
    streetAddress1: 'test',
    city: 'test',
    state: 'CA',
    postalCode: '90210',
  })

  await submitReleaseAgentFields({ typeInto, releasingAgentDetails });
  await chooseDeliveryDate({
    page,
    dateToSelect: preferredDeliveryDate.relativeDate,
    datePickerOperator: dateInputOperator,
  });
  await submitReceivingAgentFields({ typeInto, receivingAgentDetails });
  await wizardNextToReview({ page, REVIEW_PATH_CHECK });
  //fil form end
  //agreement
  await wizardNextToAgreement({ page, AGREEMENT_PATH_CHECK });
  await submitSignature({ signature, typeInto });
  await wizardCompleteToCustomerHome({ page, ROOT_PATH_CHECK });
};

export const RUN_MOVER_OWNER_SETS_UP_HHG = {
  run: async (input, getHelpers) => {
    const todaysDate = new Date();
    todaysDate.setMonth(10);
    const hhgFormData = {
      preferredPickupDate: formatRelativeDate(5),
      preferredDeliveryDate: {relativeDate: todaysDate},
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
      signature: 'mr service member',
    };

    await PerformHHGSetup({ ...input, ...getHelpers(input) }, hhgFormData);
  },
};
