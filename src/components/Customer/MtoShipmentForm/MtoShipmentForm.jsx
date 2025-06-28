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
import moment from 'moment';

import boatShipmentstyles from '../BoatShipment/BoatShipmentForm/BoatShipmentForm.module.scss';

import getShipmentOptions from './getShipmentOptions';
import styles from './MtoShipmentForm.module.scss';

import { RouterShape } from 'types';
import Callout from 'components/Callout';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { DatePickerInput } from 'components/form/fields';
import { Form } from 'components/form';
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
import { MOVE_LOCKED_WARNING, SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
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
import { isPreceedingAddressComplete } from 'shared/utils';
import { datePickerFormat, formatDate, formatDateWithUTC } from 'shared/dates';
import { handleAddressToggleChange, blankAddress } from 'utils/shipments';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

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

    const preformattedMtoShipment = {
      shipmentType,
      moveId,
      customerRemarks,
      pickup,
      delivery,
      hasSecondaryPickup: hasSecondaryPickup === 'true',
      secondaryPickup: hasSecondaryPickup === 'true' ? secondaryPickup : {},
      hasSecondaryDelivery: hasSecondaryDelivery === 'true',
      secondaryDelivery: hasSecondaryDelivery === 'true' ? secondaryDelivery : {},
      hasTertiaryPickup: hasTertiaryPickup === 'true',
      tertiaryPickup: hasTertiaryPickup === 'true' ? tertiaryPickup : {},
      hasTertiaryDelivery: hasTertiaryDelivery === 'true',
      tertiaryDelivery: hasTertiaryDelivery === 'true' ? tertiaryDelivery : {},
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
      isMoveLocked,
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
      <>
        {isMoveLocked && (
          <Alert headingLevel="h4" type="warning">
            {MOVE_LOCKED_WARNING}
          </Alert>
        )}
        <Formik
          initialValues={initialValues}
          validateOnMount
          validateOnBlur
          validationSchema={schema}
          onSubmit={this.submitMTOShipment}
        >
          {({ values, isValid, isSubmitting, setValues, handleSubmit, ...formikProps }) => {
            const {
              hasDeliveryAddress,
              hasSecondaryPickup,
              hasSecondaryDelivery,
              hasTertiaryPickup,
              hasTertiaryDelivery,
              delivery,
            } = values;

            const handleUseCurrentResidenceChange = (e) => {
              const { checked } = e.target;
              if (checked) {
                // use current residence
                setValues(
                  {
                    ...values,
                    pickup: {
                      ...values.pickup,
                      address: currentResidence,
                    },
                  },
                  { shouldValidate: true },
                );
              } else if (moveId === mtoShipment?.moveTaskOrderId) {
                // TODO - what is the purpose of this check?
                // Revert address
                setValues(
                  {
                    ...values,
                    pickup: {
                      ...values.pickup,
                      address: mtoShipment.pickupAddress,
                    },
                  },
                  { shouldValidate: true },
                );
              } else {
                // Revert address
                setValues(
                  {
                    ...values,
                    pickup: {
                      ...values.pickup,
                      address: blankAddress.address,
                    },
                  },
                  { shouldValidate: true },
                );
              }
            };

            const [isPreferredPickupDateAlertVisible, setIsPreferredPickupDateAlertVisible] = useState(false);
            const [isPreferredDeliveryDateAlertVisible, setIsPreferredDeliveryDateAlertVisible] = useState(false);
            const [preferredPickupDateAlertMessage, setPreferredPickupDateAlertMessage] = useState('');
            const [isPreferredPickupDateInvalid, setIsPreferredPickupDateInvalid] = useState(false);
            const [isPreferredPickupDateChanged, setIsPreferredPickupDateChanged] = useState(false);
            const [preferredDeliveryDateAlertMessage, setPreferredDeliveryDateAlertMessage] = useState('');
            const DEFAULT_COUNTRY_CODE = 'US';

            const onDateSelectionErrorHandler = (e) => {
              const { response } = e;
              const msg = getResponseError(response, 'failed to retrieve date selection weekend/holiday info');
              this.setState({ errorMessage: msg });
            };

            const validatePickupDate = (e) => {
              let error = validateDate(e);

              // preferredPickupDate must be in the future for non-PPM shipments
              const pickupDate = moment(formatDateWithUTC(e)).startOf('day');
              const today = moment().startOf('day');

              if (!error && isPreferredPickupDateChanged && !pickupDate.isAfter(today)) {
                setIsPreferredPickupDateInvalid(true);
                error = 'Preferred pickup date must be in the future.';
              } else {
                setIsPreferredPickupDateInvalid(false);
              }

              return error;
            };

            const handlePickupDateChange = (e) => {
              setValues({
                ...values,
                pickup: {
                  ...values.pickup,
                  requestedDate: formatDate(e, datePickerFormat),
                },
              });

              setIsPreferredPickupDateChanged(true);

              if (!validatePickupDate(e)) {
                dateSelectionWeekendHolidayCheck(
                  dateSelectionIsWeekendHoliday,
                  DEFAULT_COUNTRY_CODE,
                  new Date(e),
                  'Preferred pickup date',
                  setPreferredPickupDateAlertMessage,
                  setIsPreferredPickupDateAlertVisible,
                  onDateSelectionErrorHandler,
                );
              }
            };

            useEffect(() => {
              if (mtoShipment.requestedPickupDate !== '') {
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
                const dateSelection = new Date(mtoShipment.requestedPickupDate);
                preferredPickupDateSelectionHandler(DEFAULT_COUNTRY_CODE, dateSelection);
              }
            }, []);

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
                <NotificationScrollToTop dependency={errorMessage} />
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
                        Remember:
                        {isUB
                          ? ` You can move up to ${formatUBAllowanceWeight(
                              orders?.entitlement?.ub_allowance,
                            )} for this UB shipment. The weight of your UB is part of your authorized weight allowance`
                          : ` Your standard weight allowance is ${formatWeight(
                              orders.authorizedWeight,
                            )} total. If you are moving to an administratively restricted HHG weight location this amount may be less`}
                        . You’ll be billed for any excess weight you move.
                      </Alert>
                      <Form className={formStyles.form}>
                        {showPickupFields && (
                          <SectionWrapper className={formStyles.formSection}>
                            {showDeliveryFields && <h2>Pickup info</h2>}
                            <Fieldset legend="Date" data-testid="preferredPickupDateFieldSet">
                              <Hint id="pickupDateHint" data-testid="pickupDateHint">
                                This is the day movers would put this shipment on their truck. Packing starts earlier.
                                Dates will be finalized when you talk to your Customer Care Representative. Your
                                requested pickup/load date should be your latest preferred pickup/load date, or the date
                                you need to be out of your origin residence.
                              </Hint>
                              {requiredAsteriskMessage}
                              {isPreferredPickupDateAlertVisible && !isPreferredPickupDateInvalid && (
                                <Alert
                                  type="warning"
                                  aria-live="polite"
                                  headingLevel="h4"
                                  data-testid="preferredPickupDateAlert"
                                >
                                  {preferredPickupDateAlertMessage}
                                </Alert>
                              )}
                              <DatePickerInput
                                name="pickup.requestedDate"
                                label="Preferred pickup date"
                                showRequiredAsterisk
                                required
                                id="requestedPickupDate"
                                validate={validatePickupDate}
                                onChange={handlePickupDateChange}
                              />
                            </Fieldset>
                            <AddressFields
                              name="pickup.address"
                              legend="Pickup Address"
                              formikProps={formikProps}
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
                                  <FormGroup>
                                    <p aria-label="Do you want movers to pick up any belongings from a second address? (Must be near your pickup address. Subject to approval.">
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
                                        value="true"
                                        title="Yes, I have a second pickup address"
                                        checked={hasSecondaryPickup === 'true'}
                                        disabled={!isPreceedingAddressComplete('true', values.pickup.address)}
                                        onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                      />
                                      <Field
                                        as={Radio}
                                        id="no-secondary-pickup"
                                        data-testid="no-secondary-pickup"
                                        label="No"
                                        name="hasSecondaryPickup"
                                        value="false"
                                        title="No, I do not have a second pickup address"
                                        checked={hasSecondaryPickup !== 'true'}
                                        disabled={!isPreceedingAddressComplete('true', values.pickup.address)}
                                        onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                      />
                                    </div>
                                  </FormGroup>
                                  {hasSecondaryPickup === 'true' && (
                                    <>
                                      <h4>Second Pickup Address</h4>
                                      <AddressFields name="secondaryPickup.address" formikProps={formikProps} />
                                    </>
                                  )}
                                  {isTertiaryAddressEnabled && hasSecondaryPickup === 'true' && (
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
                                            value="true"
                                            title="Yes, I have a third pickup address"
                                            checked={hasTertiaryPickup === 'true'}
                                            disabled={
                                              !isPreceedingAddressComplete(
                                                hasSecondaryPickup,
                                                values.secondaryPickup.address,
                                              )
                                            }
                                            onChange={(e) =>
                                              handleAddressToggleChange(e, values, setValues, blankAddress)
                                            }
                                          />
                                          <Field
                                            as={Radio}
                                            id="no-tertiary-pickup"
                                            data-testid="no-tertiary-pickup"
                                            label="No"
                                            name="hasTertiaryPickup"
                                            value="false"
                                            title="No, I do not have a third pickup address"
                                            checked={hasTertiaryPickup !== 'true'}
                                            disabled={
                                              !isPreceedingAddressComplete(
                                                hasSecondaryPickup,
                                                values.secondaryPickup.address,
                                              )
                                            }
                                            onChange={(e) =>
                                              handleAddressToggleChange(e, values, setValues, blankAddress)
                                            }
                                          />
                                        </div>
                                      </FormGroup>
                                    </div>
                                  )}
                                  {isTertiaryAddressEnabled &&
                                    hasTertiaryPickup === 'true' &&
                                    hasSecondaryPickup === 'true' && (
                                      <>
                                        <h4>Third Pickup Address</h4>
                                        <AddressFields name="tertiaryPickup.address" formikProps={formikProps} />
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
                              {requiredAsteriskMessage}
                              {isPreferredDeliveryDateAlertVisible && (
                                <Alert type="warning" aria-live="polite" headingLevel="h4" required>
                                  {preferredDeliveryDateAlertMessage}
                                </Alert>
                              )}
                              <DatePickerInput
                                name="delivery.requestedDate"
                                label="Preferred delivery date"
                                showRequiredAsterisk
                                required
                                id="requestedDeliveryDate"
                                validate={validateDate}
                              />
                            </Fieldset>

                            <Fieldset legend="Delivery Address">
                              {!isNTSR && (
                                <FormGroup>
                                  <legend className="usa-label" htmlFor="hasDeliveryAddress">
                                    Do you know your delivery address yet?
                                  </legend>
                                  <div className={formStyles.radioGroup} required>
                                    <Field
                                      as={Radio}
                                      id="has-delivery-address"
                                      label="Yes"
                                      name="hasDeliveryAddress"
                                      value="true"
                                      title="Yes, I know my delivery address"
                                      checked={hasDeliveryAddress === 'true'}
                                      onChange={(e) => handleAddressToggleChange(e, values, setValues, blankAddress)}
                                    />
                                    <Field
                                      as={Radio}
                                      id="no-delivery-address"
                                      label="No"
                                      name="hasDeliveryAddress"
                                      value="false"
                                      title="No, I do not know my delivery address"
                                      checked={hasDeliveryAddress === 'false'}
                                      onChange={(e) =>
                                        handleAddressToggleChange(e, values, setValues, newDutyLocationAddress)
                                      }
                                    />
                                  </div>
                                </FormGroup>
                              )}
                              {(hasDeliveryAddress === 'true' || isNTSR) && (
                                <AddressFields
                                  name="delivery.address"
                                  formikProps={formikProps}
                                  render={(fields) => (
                                    <>
                                      {fields}
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
                                            value="true"
                                            title="Yes, I have a second delivery address"
                                            checked={hasSecondaryDelivery === 'true'}
                                            disabled={!isPreceedingAddressComplete('true', values.delivery.address)}
                                            onChange={(e) =>
                                              handleAddressToggleChange(e, values, setValues, blankAddress)
                                            }
                                          />
                                          <Field
                                            as={Radio}
                                            data-testid="no-secondary-delivery"
                                            id="no-secondary-delivery"
                                            label="No"
                                            name="hasSecondaryDelivery"
                                            value="false"
                                            title="No, I do not have a second delivery address"
                                            checked={hasSecondaryDelivery === 'false'}
                                            disabled={!isPreceedingAddressComplete('true', values.delivery.address)}
                                            onChange={(e) =>
                                              handleAddressToggleChange(e, values, setValues, blankAddress)
                                            }
                                          />
                                        </div>
                                      </FormGroup>
                                      {hasSecondaryDelivery === 'true' && (
                                        <>
                                          <h4>Second Delivery Address</h4>
                                          <AddressFields name="secondaryDelivery.address" formikProps={formikProps} />
                                        </>
                                      )}
                                      {isTertiaryAddressEnabled && hasSecondaryDelivery === 'true' && (
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
                                                value="true"
                                                title="Yes, I have a third delivery address"
                                                checked={hasTertiaryDelivery === 'true'}
                                                disabled={
                                                  !isPreceedingAddressComplete(
                                                    hasSecondaryDelivery,
                                                    values.secondaryDelivery.address,
                                                  )
                                                }
                                                onChange={(e) =>
                                                  handleAddressToggleChange(e, values, setValues, blankAddress)
                                                }
                                              />
                                              <Field
                                                as={Radio}
                                                id="no-tertiary-delivery"
                                                data-testid="no-tertiary-delivery"
                                                label="No"
                                                name="hasTertiaryDelivery"
                                                value="false"
                                                title="No, I do not have a third delivery address"
                                                checked={hasTertiaryDelivery === 'false'}
                                                disabled={
                                                  !isPreceedingAddressComplete(
                                                    hasSecondaryDelivery,
                                                    values.secondaryDelivery.address,
                                                  )
                                                }
                                                onChange={(e) =>
                                                  handleAddressToggleChange(e, values, setValues, blankAddress)
                                                }
                                              />
                                            </div>
                                          </FormGroup>
                                        </div>
                                      )}
                                      {isTertiaryAddressEnabled &&
                                        hasTertiaryDelivery === 'true' &&
                                        hasSecondaryDelivery === 'true' && (
                                          <>
                                            <h4>Third Delivery Address</h4>
                                            <AddressFields name="tertiaryDelivery.address" formikProps={formikProps} />
                                          </>
                                        )}
                                    </>
                                  )}
                                />
                              )}
                              {hasDeliveryAddress === 'false' && !isRetireeSeparatee && !isNTSR && (
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
                              {hasDeliveryAddress === 'false' && isRetireeSeparatee && !isNTSR && (
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
                                The moving company will find a storage facility approved by the government, and will
                                move your belongings there.
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
                                      Details about the facility where your things are now, including the name or
                                      address (if you know them)
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
                              disabled={!isValid || isSubmitting || isMoveLocked}
                            >
                              Save & Continue
                            </Button>
                          </div>
                        ) : (
                          <div className={formStyles.formActions}>
                            <WizardNavigation
                              disableNext={isSubmitting || !isValid || isMoveLocked}
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
      </>
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
