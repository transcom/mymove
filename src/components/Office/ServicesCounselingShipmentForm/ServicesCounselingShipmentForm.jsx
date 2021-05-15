import React from 'react';
import { bool, string, func, shape, number } from 'prop-types';
import { Formik, Field } from 'formik';
import { generatePath } from 'react-router';
import { Fieldset, Radio, Checkbox, Alert, FormGroup, Label, Textarea } from '@trussworks/react-uswds';

import getShipmentOptions from '../../Customer/MtoShipmentForm/getShipmentOptions';

import styles from './ServicesCounselingShipmentForm.module.scss';

import formStyles from 'styles/form.module.scss';
import { servicesCounselingRoutes } from 'constants/routes';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { HhgShipmentShape, HistoryShape } from 'types/customerShapes';
import { MatchShape } from 'types/officeShapes';
import { formatMtoShipmentForAPI, formatMtoShipmentForDisplay } from 'utils/formatMtoShipment';
import { createMTOShipment, patchMTOShipment, getResponseError } from 'services/internalApi';
import { DatePickerInput } from 'components/form/fields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { Form } from 'components/form/Form';
import Hint from 'components/Hint/index';
import { validateDate } from 'utils/validation';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';

const ServicesCounselingShipmentForm = ({
  match,
  history,
  newDutyStationAddress,
  selectedMoveType,
  isCreatePage,
  mtoShipment,
  serviceMember,
  currentResidence,
  updateMTOShipment,
}) => {
  const [errorMessage, setErrorMessage] = React.useState(null);

  const getShipmentNumber = () => {
    // TODO - this is not supported by IE11, shipment number should be calculable from Redux anyways
    // we should fix this also b/c it doesn't display correctly in storybook
    const { search } = window.location;
    const params = new URLSearchParams(search);
    const shipmentNumber = params.get('shipmentNumber');
    return shipmentNumber;
  };

  const shipmentType = mtoShipment.shipmentType || selectedMoveType;
  const { showDeliveryFields, showPickupFields, schema } = getShipmentOptions(shipmentType);
  const isNTS = shipmentType === SHIPMENT_OPTIONS.NTS;
  const shipmentNumber = shipmentType === SHIPMENT_OPTIONS.HHG ? getShipmentNumber() : null;

  const initialValues = formatMtoShipmentForDisplay(isCreatePage ? {} : mtoShipment);

  const optionalLabel = <span className={formStyles.optional}>Optional</span>;

  const submitMTOShipment = ({ shipmentOption, pickup, hasDeliveryAddress, delivery, customerRemarks }) => {
    const { moveCode } = match.params;

    const deliveryDetails = delivery;
    if (hasDeliveryAddress === 'no') {
      delete deliveryDetails.address;
    }

    const pendingMtoShipment = formatMtoShipmentForAPI({
      shipmentType: shipmentOption || selectedMoveType,
      moveCode,
      customerRemarks,
      pickup,
      delivery: deliveryDetails,
    });

    const moveDetailsPath = generatePath(servicesCounselingRoutes.MOVE_DETAILS_INFO_PATH, { moveCode });

    if (isCreatePage) {
      createMTOShipment(pendingMtoShipment)
        .then((response) => {
          updateMTOShipment(response);
          history.push(moveDetailsPath);
        })
        .catch((e) => {
          const { response } = e;
          const error = getResponseError(response, 'failed to create MTO shipment due to server error');

          setErrorMessage(error);
        });
    } else {
      patchMTOShipment(mtoShipment.id, pendingMtoShipment, mtoShipment.eTag)
        .then((response) => {
          updateMTOShipment(response);
          history.push(moveDetailsPath);
        })
        .catch((e) => {
          const { response } = e;
          const error = getResponseError(response, 'failed to update MTO shipment due to server error');

          setErrorMessage(error);
        });
    }
  };

  return (
    <Formik
      initialValues={initialValues}
      validateOnMount
      validateOnBlur
      validationSchema={schema}
      onSubmit={submitMTOShipment}
    >
      {({ values, isValid, isSubmitting, setValues, handleSubmit }) => {
        const { hasDeliveryAddress } = values;

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
          } else if (match.params.moveCode === mtoShipment?.moveTaskOrderId) {
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
                address: {
                  street_address_1: '',
                  street_address_2: '',
                  city: '',
                  state: '',
                  postal_code: '',
                },
              },
            });
          }
        };

        return (
          <>
            {errorMessage && (
              <Alert type="error" heading="An error occurred">
                {errorMessage}
              </Alert>
            )}

            <div className={styles.ServicesCounselingShipmentForm}>
              <ShipmentTag shipmentType={shipmentType} shipmentNumber={shipmentNumber} />

              <h1>Edit shipment details</h1>

              <SectionWrapper className={styles.weightAllowance}>
                <p>
                  <strong>Weight Allowance: </strong>
                  {serviceMember.weight_allotment.total_weight_self} lbs
                </p>
              </SectionWrapper>

              <Form className={formStyles.form}>
                {showPickupFields && (
                  <>
                    <SectionWrapper className={formStyles.formSection}>
                      {showDeliveryFields && <h2>Pickup information</h2>}
                      <Fieldset>
                        <DatePickerInput
                          name="pickup.requestedDate"
                          label="Requested pickup date"
                          id="requestedPickupDate"
                          validate={validateDate}
                        />
                      </Fieldset>

                      <AddressFields
                        name="pickup.address"
                        legend="Pickup location"
                        render={(fields) => (
                          <>
                            <Checkbox
                              data-testid="useCurrentResidence"
                              label="Use my current address"
                              name="useCurrentResidence"
                              onChange={handleUseCurrentResidenceChange}
                              id="useCurrentResidenceCheckbox"
                            />
                            {fields}
                          </>
                        )}
                      />

                      <ContactInfoFields
                        name="pickup.agent"
                        legend={<div className={formStyles.legendContent}>Releasing agent {optionalLabel}</div>}
                        render={(fields) => <>{fields}</>}
                      />
                    </SectionWrapper>
                  </>
                )}

                {showDeliveryFields && (
                  <>
                    <SectionWrapper className={formStyles.formSection}>
                      {showPickupFields && <h2>Delivery information</h2>}
                      <Fieldset>
                        <DatePickerInput
                          name="delivery.requestedDate"
                          label="Requested delivery date"
                          id="requestedDeliveryDate"
                          validate={validateDate}
                        />
                      </Fieldset>

                      <Fieldset legend="Delivery location">
                        <FormGroup>
                          <p>Does the customer know their delivery address yet?</p>
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
                        {hasDeliveryAddress === 'yes' ? (
                          <AddressFields name="delivery.address" render={(fields) => <>{fields}</>} />
                        ) : (
                          <p>
                            We can use the zip of their new duty station:
                            <br />
                            <strong>
                              {newDutyStationAddress.city}, {newDutyStationAddress.state}{' '}
                              {newDutyStationAddress.postal_code}{' '}
                            </strong>
                          </p>
                        )}
                      </Fieldset>

                      <ContactInfoFields
                        name="delivery.agent"
                        legend={<div className={formStyles.legendContent}>Receiving agent {optionalLabel}</div>}
                        render={(fields) => <>{fields}</>}
                      />
                    </SectionWrapper>
                  </>
                )}

                {isNTS && (
                  <>
                    <SectionWrapper className={formStyles.formSection} data-testid="nts-what-to-expect">
                      <Fieldset legend="What you can expect">
                        <p>
                          The moving company will find a storage facility approved by the government, and will move your
                          belongings there.
                        </p>
                        <p>
                          You’ll need to schedule an NTS release shipment to get your items back, most likely as part of
                          a future move.
                        </p>
                      </Fieldset>
                    </SectionWrapper>
                  </>
                )}

                <SectionWrapper className={formStyles.formSection}>
                  <Fieldset legend={<h2 className={formStyles.legendContent}>Remarks {optionalLabel}</h2>}>
                    <Label htmlFor="customerRemarks">Customer remarks</Label>
                    <Field
                      as={Textarea}
                      data-testid="remarks"
                      name="customerRemarks"
                      className={`${formStyles.remarks}`}
                      placeholder=""
                      id="customerRemarks"
                      maxLength={500}
                    />
                    <Hint>
                      <p>500 characters</p>
                    </Hint>

                    <Label htmlFor="counselorRemarks">Counselor remarks</Label>
                    <Field
                      as={Textarea}
                      data-testid="counselor-remarks"
                      name="counselorRemarks"
                      className={`${formStyles.remarks}`}
                      placeholder=""
                      id="counselorRemarks"
                      maxLength={500}
                    />
                    <Hint>
                      <p>500 characters</p>
                    </Hint>
                  </Fieldset>
                </SectionWrapper>

                <div className={formStyles.formActions}>
                  <WizardNavigation
                    disableNext={isSubmitting || !isValid}
                    editMode
                    onNextClick={handleSubmit}
                    onBackClick={history.goBack}
                    onCancelClick={history.goBack}
                  />
                </div>
              </Form>
            </div>
          </>
        );
      }}
    </Formik>
  );
};

ServicesCounselingShipmentForm.propTypes = {
  match: MatchShape,
  history: HistoryShape,
  updateMTOShipment: func.isRequired,
  isCreatePage: bool,
  currentResidence: AddressShape.isRequired,
  newDutyStationAddress: SimpleAddressShape,
  selectedMoveType: string.isRequired,
  mtoShipment: HhgShipmentShape,
  serviceMember: shape({
    weight_allotment: shape({
      total_weight_self: number,
    }),
  }).isRequired,
};

ServicesCounselingShipmentForm.defaultProps = {
  isCreatePage: false,
  match: { isExact: false, params: { moveCode: '', shipmentId: '' } },
  history: { goBack: () => {}, push: () => {} },
  newDutyStationAddress: {
    city: '',
    state: '',
    postal_code: '',
  },
  mtoShipment: {
    id: '',
    customerRemarks: '',
    counselorRemarks: '',
    requestedPickupDate: '',
    requestedDeliveryDate: '',
    destinationAddress: {
      city: '',
      postal_code: '',
      state: '',
      street_address_1: '',
    },
  },
};

export default ServicesCounselingShipmentForm;
