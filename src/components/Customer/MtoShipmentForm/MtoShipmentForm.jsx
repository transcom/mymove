import React, { Component, useState, useEffect } from 'react';
import { bool, func, string } from 'prop-types';
import { Field, Formik } from 'formik';
import { generatePath } from 'react-router-dom';
import {
  Alert,
  Checkbox,
  Fieldset,
  FormGroup,
  Grid,
  GridContainer,
  Label,
  Radio,
  Textarea,
  Button,
} from '@trussworks/react-uswds';

import boatShipmentstyles from '../BoatShipment/BoatShipmentForm/BoatShipmentForm.module.scss';

import getShipmentOptions from './getShipmentOptions';
import styles from './MtoShipmentForm.module.scss';

import { RouterShape } from 'types';
import Callout from 'components/Callout';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { DatePickerInput } from 'components/form/fields';
import { Form } from 'components/form/Form';
import Hint from 'components/Hint/index';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { customerRoutes } from 'constants/routes';
import { roleTypes } from 'constants/userRoles';
import { shipmentForm } from 'content/shipments';
import {
  createMTOShipment,
  getResponseError,
  patchMTOShipment,
  dateSelectionIsWeekendHoliday,
} from 'services/internalApi';
import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
import formStyles from 'styles/form.module.scss';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { OrdersShape } from 'types/customerShapes';
import { ShipmentShape } from 'types/shipment';
import { formatMtoShipmentForAPI, formatMtoShipmentForDisplay } from 'utils/formatMtoShipment';
import { formatUBAllowanceWeight, formatWeight } from 'utils/formatters';
import { validateDate } from 'utils/validation';
import withRouter from 'utils/routing';
import { ORDERS_TYPE } from 'constants/orders';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { dateSelectionWeekendHolidayCheck } from 'utils/calendar';

const blankAddress = {
  address: {
    streetAddress1: '',
    streetAddress2: '',
    city: '',
    state: '',
    postalCode: '',
  },
};

class MtoShipmentForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMessage: null,
      isTertiaryAddressEnabled: false,
    };
  }

  componentDidMount() {
    isBooleanFlagEnabled('third_address_available').then((enabled) => {
      this.setState({
        isTertiaryAddressEnabled: enabled,
      });
    });
  }

  submitMTOShipment = ({
    pickup,
    hasDeliveryAddress,
    delivery,
    customerRemarks,
    hasSecondaryPickup,
    secondaryPickup,
    hasSecondaryDelivery,
    secondaryDelivery,
    hasTertiaryDelivery,
    hasTertiaryPickup,
    tertiaryDelivery,
    tertiaryPickup,
  }) => {
    const {
      router: { navigate, params },
      shipmentType,
      isCreatePage,
      mtoShipment,
      updateMTOShipment,
    } = this.props;

    const { moveId } = params;

    const isNTSR = shipmentType === SHIPMENT_OPTIONS.NTSR;
    const saveDeliveryAddress = hasDeliveryAddress === 'yes' || isNTSR;

    const preformattedMtoShipment = {
      shipmentType,
      moveId,
      customerRemarks,
      pickup,
      delivery: {
        ...delivery,
        address: saveDeliveryAddress ? delivery.address : undefined,
      },
      hasSecondaryPickup: hasSecondaryPickup === 'yes',
      secondaryPickup: hasSecondaryPickup === 'yes' ? secondaryPickup : {},
      hasSecondaryDelivery: hasSecondaryDelivery === 'yes',
      secondaryDelivery: hasSecondaryDelivery === 'yes' ? secondaryDelivery : {},
      hasTertiaryPickup: hasTertiaryPickup === 'yes',
      tertiaryPickup: hasTertiaryPickup === 'yes' ? tertiaryPickup : {},
      hasTertiaryDelivery: hasTertiaryDelivery === 'yes',
      tertiaryDelivery: hasTertiaryDelivery === 'yes' ? tertiaryDelivery : {},
    };

    const pendingMtoShipment = formatMtoShipmentForAPI(preformattedMtoShipment);

    const reviewPath = generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId });

    if (isCreatePage) {
      createMTOShipment(pendingMtoShipment)
        .then((response) => {
          updateMTOShipment(response);
          navigate(reviewPath);
        })
        .catch((e) => {
          const { response } = e;
          const errorMessage = getResponseError(response, 'failed to create MTO shipment due to server error');

          this.setState({ errorMessage });
        });
    } else {
      patchMTOShipment(mtoShipment.id, pendingMtoShipment, mtoShipment.eTag)
        .then((response) => {
          updateMTOShipment(response);
          navigate(reviewPath);
        })
        .catch((e) => {
          const { response } = e;
          const errorMessage = getResponseError(response, 'failed to update MTO shipment due to server error');

          this.setState({ errorMessage });
        });
    }
  };

  // eslint-disable-next-line class-methods-use-this
  getShipmentNumber = () => {
    const {
      router: {
        location: { search },
      },
    } = this.props;

    const params = new URLSearchParams(search);
    const shipmentNumber = params.get('shipmentNumber');
    return shipmentNumber;
  };

  render() {
    const {
      newDutyLocationAddress,
      shipmentType,
      isCreatePage,
      mtoShipment,
      orders,
      currentResidence,
      router: { params, navigate },
      handleBack,
    } = this.props;

    const { moveId } = params;
    const { isTertiaryAddressEnabled } = this.state;
    const { errorMessage } = this.state;
    const { showDeliveryFields, showPickupFields, schema } = getShipmentOptions(shipmentType, roleTypes.CUSTOMER);
    const isNTS = shipmentType === SHIPMENT_OPTIONS.NTS;
    const isNTSR = shipmentType === SHIPMENT_OPTIONS.NTSR;
    const isBoat = shipmentType === SHIPMENT_TYPES.BOAT_HAUL_AWAY || shipmentType === SHIPMENT_TYPES.BOAT_TOW_AWAY;
    const isMobileHome = shipmentType === SHIPMENT_TYPES.MOBILE_HOME;
    const isUB = shipmentType === SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE;
    const shipmentNumber =
      shipmentType === SHIPMENT_OPTIONS.HHG || isBoat || isMobileHome || isUB ? this.getShipmentNumber() : null;
    const isRetireeSeparatee =
      orders.orders_type === ORDERS_TYPE.RETIREMENT || orders.orders_type === ORDERS_TYPE.SEPARATION;

    const initialValues = formatMtoShipmentForDisplay(
      isCreatePage && !mtoShipment?.requestedPickupDate ? {} : mtoShipment, // check if data carried over from boat shipment
    );

    return (
      <Formik
        initialValues={initialValues}
        validateOnMount
        validateOnBlur
        validationSchema={schema}
        onSubmit={this.submitMTOShipment}
      >
        {({ values, isValid, isSubmitting, setValues, handleSubmit }) => {
          const {
            hasDeliveryAddress,
            hasSecondaryPickup,
            hasSecondaryDelivery,
            hasTertiaryPickup,
            hasTertiaryDelivery,
            pickup,
            delivery,
          } = values;

          const handleUseCurrentResidenceChange = (e) => {
            const { checked } = e.target;
            if (checked) {
              // use current residence
              setValues({
                ...values,
                pickup: {
                  ...values.pickup,
                  address: currentResidence,
                },
              });
            } else if (moveId === mtoShipment?.moveTaskOrderId) {
              // TODO - what is the purpose of this check?
              // Revert address
              setValues({
                ...values,
                pickup: {
                  ...values.pickup,
                  address: mtoShipment.pickupAddress,
                },
              });
            } else {
              // Revert address
              setValues({
                ...values,
                pickup: {
                  ...values.pickup,
                  ...blankAddress,
                },
              });
            }
          };

          const [isPreferredPickupDateAlertVisible, setIsPreferredPickupDateAlertVisible] = useState(false);
          const [isPreferredDeliveryDateAlertVisible, setIsPreferredDeliveryDateAlertVisible] = useState(false);
          const [preferredPickupDateAlertMessage, setPreferredPickupDateAlertMessage] = useState('');
          const [preferredDeliveryDateAlertMessage, setPreferredDeliveryDateAlertMessage] = useState('');
          const DEFAULT_COUNTRY_CODE = 'US';

          const onDateSelectionErrorHandler = (e) => {
            const { response } = e;
            const msg = getResponseError(response, 'failed to retrieve date selection weekend/holiday info');
            this.setState({ errorMessage: msg });
          };

          useEffect(() => {
            if (pickup?.requestedDate !== '') {
              const preferredPickupDateSelectionHandler = (countryCode, date) => {
                dateSelectionWeekendHolidayCheck(
                  dateSelectionIsWeekendHoliday,
                  countryCode,
                  date,
                  'Preferred pickup date',
                  setPreferredPickupDateAlertMessage,
                  setIsPreferredPickupDateAlertVisible,
                  onDateSelectionErrorHandler,
                );
              };
              const dateSelection = new Date(pickup.requestedDate);
              preferredPickupDateSelectionHandler(DEFAULT_COUNTRY_CODE, dateSelection);
            }
          }, [pickup.requestedDate]);

          useEffect(() => {
            if (delivery?.requestedDate !== '') {
              const preferredDeliveryDateSelectionHandler = (countryCode, date) => {
                dateSelectionWeekendHolidayCheck(
                  dateSelectionIsWeekendHoliday,
                  countryCode,
                  date,
                  'Preferred delivery date',
                  setPreferredDeliveryDateAlertMessage,
                  setIsPreferredDeliveryDateAlertVisible,
                  onDateSelectionErrorHandler,
                );
              };
              const dateSelection = new Date(delivery.requestedDate);
              preferredDeliveryDateSelectionHandler(DEFAULT_COUNTRY_CODE, dateSelection);
            }
          }, [delivery.requestedDate]);

          return (
            <GridContainer>
              <Grid row>
                <Grid col desktop={{ col: 8, offset: 2 }}>
                  {errorMessage && (
                    <Alert type="error" headingLevel="h4" heading="An error occurred">
                      {errorMessage}
                    </Alert>
                  )}

                  <div className={styles.MTOShipmentForm}>
                    <ShipmentTag shipmentType={shipmentType} shipmentNumber={shipmentNumber} />
                    <h1>{shipmentForm.header[`${shipmentType}`]}</h1>
                    <Alert headingLevel="h4" type="info" noIcon>
                      Remember: You can move
                      {isUB
                        ? ` up to ${formatUBAllowanceWeight(
                            orders?.entitlement?.ub_allowance,
                          )} for this UB shipment. The weight of your UB is part of your authorized weight allowance`
                        : ` ${formatWeight(orders.authorizedWeight)} total`}
                      . You’ll be billed for any excess weight you move.
                    </Alert>
                    <Form className={formStyles.form}>
                      {showPickupFields && (
                        <SectionWrapper className={formStyles.formSection}>
                          {showDeliveryFields && <h2>Pickup info</h2>}
                          <Fieldset legend="Date">
                            <Hint id="pickupDateHint" data-testid="pickupDateHint">
                              This is the day movers would put this shipment on their truck. Packing starts earlier.
                              Dates will be finalized when you talk to your Customer Care Representative. Your requested
                              pickup/load date should be your latest preferred pickup/load date, or the date you need to
                              be out of your origin residence.
                            </Hint>
                            {isPreferredPickupDateAlertVisible && (
                              <Alert type="warning" aria-live="polite" headingLevel="h4">
                                {preferredPickupDateAlertMessage}
                              </Alert>
                            )}
                            <DatePickerInput
                              name="pickup.requestedDate"
                              label="Preferred pickup date"
                              id="requestedPickupDate"
                              hint="Required"
                              validate={validateDate}
                            />
                          </Fieldset>

                          <AddressFields
                            name="pickup.address"
                            legend="Pickup Address"
                            labelHint="Required"
                            render={(fields) => (
                              <>
                                <p>What address are the movers picking up from?</p>
                                <Checkbox
                                  data-testid="useCurrentResidence"
                                  label="Use my current address"
                                  name="useCurrentResidence"
                                  onChange={handleUseCurrentResidenceChange}
                                  id="useCurrentResidenceCheckbox"
                                />
                                {fields}
                                <h4>Second Pickup Address</h4>
                                <FormGroup>
                                  <p>
                                    Do you want movers to pick up any belongings from a second address? (Must be near
                                    your pickup address. Subject to approval.)
                                  </p>
                                  <div className={formStyles.radioGroup}>
                                    <Field
                                      as={Radio}
                                      id="has-secondary-pickup"
                                      data-testid="has-secondary-pickup"
                                      label="Yes"
                                      name="hasSecondaryPickup"
                                      value="yes"
                                      title="Yes, I have a second pickup address"
                                      checked={hasSecondaryPickup === 'yes'}
                                    />
                                    <Field
                                      as={Radio}
                                      id="no-secondary-pickup"
                                      data-testid="no-secondary-pickup"
                                      label="No"
                                      name="hasSecondaryPickup"
                                      value="no"
                                      title="No, I do not have a second pickup address"
                                      checked={hasSecondaryPickup !== 'yes'}
                                    />
                                  </div>
                                </FormGroup>
                                {hasSecondaryPickup === 'yes' && (
                                  <AddressFields name="secondaryPickup.address" labelHint="Required" />
                                )}
                                {isTertiaryAddressEnabled && hasSecondaryPickup === 'yes' && (
                                  <div>
                                    <FormGroup>
                                      <p>Do you want movers to pick up any belongings from a third address?</p>
                                      <div className={formStyles.radioGroup}>
                                        <Field
                                          as={Radio}
                                          id="has-tertiary-pickup"
                                          data-testid="has-tertiary-pickup"
                                          label="Yes"
                                          name="hasTertiaryPickup"
                                          value="yes"
                                          title="Yes, I have a third pickup address"
                                          checked={hasTertiaryPickup === 'yes'}
                                        />
                                        <Field
                                          as={Radio}
                                          id="no-tertiary-pickup"
                                          data-testid="no-tertiary-pickup"
                                          label="No"
                                          name="hasTertiaryPickup"
                                          value="no"
                                          title="No, I do not have a third pickup address"
                                          checked={hasTertiaryPickup !== 'yes'}
                                        />
                                      </div>
                                    </FormGroup>
                                  </div>
                                )}
                                {isTertiaryAddressEnabled &&
                                  hasTertiaryPickup === 'yes' &&
                                  hasSecondaryPickup === 'yes' && (
                                    <>
                                      <h3>Third Pickup Address</h3>
                                      <AddressFields name="tertiaryPickup.address" labelHint="Required" />
                                    </>
                                  )}
                              </>
                            )}
                          />

                          <ContactInfoFields
                            name="pickup.agent"
                            legend={<div className={formStyles.legendContent}>Releasing agent</div>}
                            render={(fields) => (
                              <>
                                <p>Who can let the movers pick up your personal property if you are not there?</p>
                                {fields}
                              </>
                            )}
                          />
                        </SectionWrapper>
                      )}

                      {showDeliveryFields && (
                        <SectionWrapper className={formStyles.formSection}>
                          {showPickupFields && <h2>Delivery Address info</h2>}
                          <Fieldset legend="Date">
                            <Hint>
                              You will finalize an actual delivery date later by talking with your Customer Care
                              Representative once the shipment is underway.
                            </Hint>
                            {isPreferredDeliveryDateAlertVisible && (
                              <Alert type="warning" aria-live="polite" headingLevel="h4">
                                {preferredDeliveryDateAlertMessage}
                              </Alert>
                            )}
                            <DatePickerInput
                              name="delivery.requestedDate"
                              label="Preferred delivery date"
                              id="requestedDeliveryDate"
                              validate={validateDate}
                              hint="Required"
                            />
                          </Fieldset>

                          <Fieldset legend="Delivery Address">
                            {!isNTSR && (
                              <FormGroup>
                                <Label hint="Required" htmlFor="hasDeliveryAddress">
                                  Do you know your delivery address yet?
                                </Label>
                                <div className={formStyles.radioGroup}>
                                  <Field
                                    as={Radio}
                                    id="has-delivery-address"
                                    label="Yes"
                                    name="hasDeliveryAddress"
                                    value="yes"
                                    title="Yes, I know my delivery address"
                                    checked={hasDeliveryAddress === 'yes'}
                                  />
                                  <Field
                                    as={Radio}
                                    id="no-delivery-address"
                                    label="No"
                                    name="hasDeliveryAddress"
                                    value="no"
                                    title="No, I do not know my delivery address"
                                    checked={hasDeliveryAddress === 'no'}
                                  />
                                </div>
                              </FormGroup>
                            )}
                            {(hasDeliveryAddress === 'yes' || isNTSR) && (
                              <AddressFields
                                name="delivery.address"
                                labelHint="Required"
                                render={(fields) => (
                                  <>
                                    {fields}
                                    <h4>Second Delivery Address</h4>
                                    <FormGroup>
                                      <p>
                                        Do you want the movers to deliver any belongings to a second address? (Must be
                                        near your delivery address. Subject to approval.)
                                      </p>
                                      <div className={formStyles.radioGroup}>
                                        <Field
                                          as={Radio}
                                          data-testid="has-secondary-delivery"
                                          id="has-secondary-delivery"
                                          label="Yes"
                                          name="hasSecondaryDelivery"
                                          value="yes"
                                          title="Yes, I have a second delivery address"
                                          checked={hasSecondaryDelivery === 'yes'}
                                        />
                                        <Field
                                          as={Radio}
                                          data-testid="no-secondary-delivery"
                                          id="no-secondary-delivery"
                                          label="No"
                                          name="hasSecondaryDelivery"
                                          value="no"
                                          title="No, I do not have a second delivery address"
                                          checked={hasSecondaryDelivery === 'no'}
                                        />
                                      </div>
                                    </FormGroup>
                                    {hasSecondaryDelivery === 'yes' && (
                                      <AddressFields name="secondaryDelivery.address" labelHint="Required" />
                                    )}
                                    {isTertiaryAddressEnabled && hasSecondaryDelivery === 'yes' && (
                                      <div>
                                        <FormGroup>
                                          <p>Do you want movers to deliver any belongings to a third address?</p>
                                          <div className={formStyles.radioGroup}>
                                            <Field
                                              as={Radio}
                                              id="has-tertiary-delivery"
                                              data-testid="has-tertiary-delivery"
                                              label="Yes"
                                              name="hasTertiaryDelivery"
                                              value="yes"
                                              title="Yes, I have a third delivery address"
                                              checked={hasTertiaryDelivery === 'yes'}
                                            />
                                            <Field
                                              as={Radio}
                                              id="no-tertiary-delivery"
                                              data-testid="no-tertiary-delivery"
                                              label="No"
                                              name="hasTertiaryDelivery"
                                              value="no"
                                              title="No, I do not have a third delivery address"
                                              checked={hasTertiaryDelivery === 'no'}
                                            />
                                          </div>
                                        </FormGroup>
                                      </div>
                                    )}
                                    {isTertiaryAddressEnabled &&
                                      hasTertiaryDelivery === 'yes' &&
                                      hasSecondaryDelivery === 'yes' && (
                                        <>
                                          <h4>Third Delivery Address</h4>
                                          <AddressFields name="tertiaryDelivery.address" labelHint="Required" />
                                        </>
                                      )}
                                  </>
                                )}
                              />
                            )}
                            {hasDeliveryAddress === 'no' && !isRetireeSeparatee && !isNTSR && (
                              <p>
                                We can use the zip of your new duty location.
                                <br />
                                <strong>
                                  {newDutyLocationAddress.city}, {newDutyLocationAddress.state}{' '}
                                  {newDutyLocationAddress.postalCode}{' '}
                                </strong>
                                <br />
                                You can add the specific delivery address later, once you know it.
                              </p>
                            )}
                            {hasDeliveryAddress === 'no' && isRetireeSeparatee && !isNTSR && (
                              <p>
                                We can use the zip of the HOR, PLEAD or HOS you entered with your orders.
                                <br />
                                <strong>
                                  {newDutyLocationAddress.city}, {newDutyLocationAddress.state}{' '}
                                  {newDutyLocationAddress.postalCode}{' '}
                                </strong>
                                <br />
                              </p>
                            )}
                          </Fieldset>

                          <ContactInfoFields
                            name="delivery.agent"
                            legend={<div className={formStyles.legendContent}>Receiving agent</div>}
                            render={(fields) => (
                              <>
                                <p>Who can take delivery for you if the movers arrive and you are not there?</p>
                                {fields}
                              </>
                            )}
                          />
                        </SectionWrapper>
                      )}

                      {isNTS && (
                        <SectionWrapper className={formStyles.formSection} data-testid="nts-what-to-expect">
                          <Fieldset legend="What you can expect">
                            <p>
                              The moving company will find a storage facility approved by the government, and will move
                              your belongings there.
                            </p>
                            <p>
                              You will need to schedule an NTS release shipment to get your items back, most likely as
                              part of a future move.
                            </p>
                          </Fieldset>
                        </SectionWrapper>
                      )}

                      {!isBoat && !isMobileHome && (
                        <SectionWrapper className={formStyles.formSection}>
                          <Fieldset legend={<div className={formStyles.legendContent}>Remarks</div>}>
                            <Label htmlFor="customerRemarks">
                              Are there things about this shipment that your counselor or movers should discuss with
                              you?
                            </Label>

                            <Callout>
                              Examples
                              <ul>
                                {isNTSR && (
                                  <li>
                                    Details about the facility where your things are now, including the name or address
                                    (if you know them)
                                  </li>
                                )}
                                <li>Large, bulky, or fragile items</li>
                                <li>Access info for your pickup or delivery address</li>
                                <li>You’re shipping weapons or alcohol</li>
                              </ul>
                            </Callout>

                            <Field
                              as={Textarea}
                              data-testid="remarks"
                              name="customerRemarks"
                              className={`${formStyles.remarks}`}
                              placeholder="Do not itemize your personal property here. Your movers will help do that when they talk to you."
                              id="customerRemarks"
                              maxLength={250}
                            />
                            <Hint>
                              <p>250 characters</p>
                            </Hint>
                          </Fieldset>
                        </SectionWrapper>
                      )}
                      <Hint darkText>
                        <p>You can change details about your move by talking with your counselor or your movers</p>
                      </Hint>

                      {isBoat || isMobileHome ? (
                        <div className={boatShipmentstyles.buttonContainer}>
                          <Button
                            className={boatShipmentstyles.backButton}
                            type="button"
                            onClick={handleBack}
                            secondary
                            outline
                          >
                            Back
                          </Button>
                          <Button
                            className={boatShipmentstyles.saveButton}
                            type="button"
                            onClick={handleSubmit}
                            disabled={!isValid || isSubmitting}
                          >
                            Save & Continue
                          </Button>
                        </div>
                      ) : (
                        <div className={formStyles.formActions}>
                          <WizardNavigation
                            disableNext={isSubmitting || !isValid}
                            editMode={!isCreatePage}
                            onNextClick={handleSubmit}
                            onBackClick={() => {
                              navigate(-1);
                            }}
                            onCancelClick={() => {
                              navigate(-1);
                            }}
                          />
                        </div>
                      )}
                    </Form>
                  </div>
                </Grid>
              </Grid>
            </GridContainer>
          );
        }}
      </Formik>
    );
  }
}

MtoShipmentForm.propTypes = {
  router: RouterShape.isRequired,
  updateMTOShipment: func.isRequired,
  isCreatePage: bool,
  currentResidence: AddressShape.isRequired,
  newDutyLocationAddress: SimpleAddressShape,
  shipmentType: string.isRequired,
  mtoShipment: ShipmentShape,
  orders: OrdersShape,
};

MtoShipmentForm.defaultProps = {
  isCreatePage: false,
  newDutyLocationAddress: {
    city: '',
    state: '',
    postalCode: '',
  },
  mtoShipment: {
    id: '',
    customerRemarks: '',
    requestedPickupDate: '',
    requestedDeliveryDate: '',
    destinationAddress: {
      city: '',
      postalCode: '',
      state: '',
      streetAddress1: '',
    },
  },
  orders: {},
};

export default withRouter(MtoShipmentForm);
