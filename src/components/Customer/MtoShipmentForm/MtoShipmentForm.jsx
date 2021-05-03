import React, { Component } from 'react';
import { bool, string, func, shape, number } from 'prop-types';
import { Formik, Field } from 'formik';
import { generatePath } from 'react-router';
import {
  Fieldset,
  Radio,
  Checkbox,
  Alert,
  FormGroup,
  Label,
  Textarea,
  GridContainer,
  Grid,
} from '@trussworks/react-uswds';

import getShipmentOptions from './getShipmentOptions';
import styles from './MtoShipmentForm.module.scss';

import { customerRoutes } from 'constants/routes';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { HhgShipmentShape, MatchShape, HistoryShape } from 'types/customerShapes';
import { formatMtoShipmentForAPI, formatMtoShipmentForDisplay } from 'utils/formatMtoShipment';
import { createMTOShipment, patchMTOShipment, getResponseError } from 'services/internalApi';
import { shipmentForm } from 'content/shipments';
import { DatePickerInput } from 'components/form/fields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { Form } from 'components/form/Form';
import Hint from 'components/Hint/index';
import { validateDate } from 'utils/validation';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';

class MtoShipmentForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMessage: null,
    };
  }

  submitMTOShipment = ({ shipmentType, pickup, hasDeliveryAddress, delivery, customerRemarks }) => {
    const { history, match, selectedMoveType, isCreatePage, mtoShipment, updateMTOShipment } = this.props;
    const { moveId } = match.params;

    const deliveryDetails = delivery;
    if (hasDeliveryAddress === 'no') {
      delete deliveryDetails.address;
    }

    const pendingMtoShipment = formatMtoShipmentForAPI({
      shipmentType: shipmentType || selectedMoveType,
      moveId,
      customerRemarks,
      pickup,
      delivery: deliveryDetails,
    });

    const reviewPath = generatePath(customerRoutes.MOVE_REVIEW_PATH, { moveId });

    if (isCreatePage) {
      createMTOShipment(pendingMtoShipment)
        .then((response) => {
          updateMTOShipment(response);
          history.push(reviewPath);
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
          history.push(reviewPath);
        })
        .catch((e) => {
          const { response } = e;
          const errorMessage = getResponseError(response, 'failed to update MTO shipment due to server error');

          this.setState({ errorMessage });
        });
    }
  };

  getShipmentNumber = () => {
    // TODO - this is not supported by IE11, shipment number should be calculable from Redux anyways
    // we should fix this also b/c it doesn't display correctly in storybook
    const { search } = window.location;
    const params = new URLSearchParams(search);
    const shipmentNumber = params.get('shipmentNumber');
    return shipmentNumber;
  };

  render() {
    const {
      match,
      history,
      newDutyStationAddress,
      selectedMoveType,
      isCreatePage,
      mtoShipment,
      serviceMember,
      currentResidence,
    } = this.props;

    const { errorMessage } = this.state;

    const shipmentType = mtoShipment.shipmentType || selectedMoveType;
    const { showDeliveryFields, showPickupFields, schema } = getShipmentOptions(shipmentType);
    const isNTS = shipmentType === SHIPMENT_OPTIONS.NTS;
    const shipmentNumber = shipmentType === SHIPMENT_OPTIONS.HHG ? this.getShipmentNumber() : null;

    const initialValues = formatMtoShipmentForDisplay(isCreatePage ? {} : mtoShipment);

    const optionalLabel = <span className={styles.optional}>Optional</span>;

    return (
      <Formik
        initialValues={initialValues}
        validateOnMount
        validateOnBlur
        validationSchema={schema}
        onSubmit={this.submitMTOShipment}
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
            } else if (match.params.moveId === mtoShipment?.moveTaskOrderId) {
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
            <GridContainer>
              <Grid row>
                <Grid col desktop={{ col: 8, offset: 2 }}>
                  {errorMessage && (
                    <Alert type="error" heading="An error occurred">
                      {errorMessage}
                    </Alert>
                  )}

                  <div className={styles.MTOShipmentForm}>
                    <ShipmentTag shipmentType={shipmentType} shipmentNumber={shipmentNumber} />

                    <h1>{shipmentForm.header[`${shipmentType}`]}</h1>

                    <Alert type="info" noIcon>
                      Remember: You can move {serviceMember.weight_allotment.total_weight_self} lbs total. You’ll be
                      billed for any excess weight you move.
                    </Alert>

                    <Form className={styles.form}>
                      {showPickupFields && (
                        <>
                          <SectionWrapper className={styles.formSection}>
                            {showDeliveryFields && <h2>Pickup information</h2>}
                            <Fieldset legend="Pickup date">
                              <DatePickerInput
                                name="pickup.requestedDate"
                                label="Requested pickup date"
                                id="requestedPickupDate"
                                validate={validateDate}
                              />
                              <Hint id="pickupDateHint">
                                <p>
                                  Movers will contact you to schedule the actual pickup date. That date should fall
                                  within 7 days of your requested date. Tip: Avoid scheduling multiple shipments on the
                                  same day.
                                </p>
                              </Hint>
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
                                  <Hint>
                                    <p>
                                      If you have more things at another pickup location, you can schedule a shipment
                                      for them later.
                                    </p>
                                  </Hint>
                                </>
                              )}
                            />

                            <ContactInfoFields
                              name="pickup.agent"
                              legend={<div className={styles.legendContent}>Releasing agent {optionalLabel}</div>}
                              render={(fields) => (
                                <>
                                  <p>Who can let the movers pick up your things if you’re not there?</p>
                                  {fields}
                                </>
                              )}
                            />
                          </SectionWrapper>
                        </>
                      )}

                      {showDeliveryFields && (
                        <>
                          <SectionWrapper className={styles.formSection}>
                            {showPickupFields && <h2>Delivery information</h2>}
                            <Fieldset legend="Delivery date">
                              <DatePickerInput
                                name="delivery.requestedDate"
                                label="Requested delivery date"
                                id="requestedDeliveryDate"
                                validate={validateDate}
                              />
                              <Hint>
                                <p>
                                  Shipments can take several weeks to arrive, depending on how far they’re going. Your
                                  movers will contact you close to the date you select to coordinate delivery.
                                </p>
                              </Hint>
                            </Fieldset>

                            <Fieldset legend="Delivery location">
                              <FormGroup>
                                <p>Do you know your delivery address yet?</p>
                                <div className={styles.radioGroup}>
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
                                <AddressFields
                                  name="delivery.address"
                                  render={(fields) => (
                                    <>
                                      {fields}
                                      <Hint>
                                        <p>
                                          If you have more things to go to another destination, you can schedule a
                                          shipment for them later.
                                        </p>
                                      </Hint>
                                    </>
                                  )}
                                />
                              ) : (
                                <p>
                                  We can use the zip of your new duty station.
                                  <br />
                                  <strong>
                                    {newDutyStationAddress.city}, {newDutyStationAddress.state}{' '}
                                    {newDutyStationAddress.postal_code}{' '}
                                  </strong>
                                  <br />
                                  You can add the specific delivery address later, once you know it.
                                </p>
                              )}
                            </Fieldset>

                            <ContactInfoFields
                              name="delivery.agent"
                              legend={<div className={styles.legendContent}>Receiving agent {optionalLabel}</div>}
                              render={(fields) => (
                                <>
                                  <p>Who can take delivery for you if the movers arrive and you’re not there?</p>
                                  {fields}
                                </>
                              )}
                            />
                          </SectionWrapper>
                        </>
                      )}

                      {isNTS && (
                        <>
                          <SectionWrapper className={styles.formSection} data-testid="nts-what-to-expect">
                            <Fieldset legend="What you can expect">
                              <p>
                                The moving company will find a storage facility approved by the government, and will
                                move your belongings there.
                              </p>
                              <p>
                                You’ll need to schedule an NTS release shipment to get your items back, most likely as
                                part of a future move.
                              </p>
                            </Fieldset>
                          </SectionWrapper>
                        </>
                      )}

                      <SectionWrapper className={styles.formSection}>
                        <Fieldset legend={<div className={styles.legendContent}>Remarks {optionalLabel}</div>}>
                          <Label htmlFor="customerRemarks">
                            Is there anything special about this shipment that the movers should know?
                          </Label>

                          <div className={styles.remarksExamples}>
                            Examples
                            <ul>
                              <li>Things that might need special handling</li>
                              <li>Access info for a location</li>
                              <li>Weapons or alcohol</li>
                            </ul>
                          </div>

                          <Field
                            as={Textarea}
                            data-testid="remarks"
                            name="customerRemarks"
                            className={`${styles.remarks}`}
                            placeholder="You don’t need to list all your belongings here. Your mover will get those details later."
                            id="customerRemarks"
                            maxLength={250}
                          />
                          <Hint>
                            <p>250 characters</p>
                          </Hint>
                        </Fieldset>
                      </SectionWrapper>

                      <Hint>
                        <p>
                          You can change details for your shipment when you talk to your move counselor or the person
                          who’s your point of contact with the movers. You can also edit in MilMove up to 24 hours
                          before your final pickup date.
                        </p>
                      </Hint>

                      <div className={styles.formActions}>
                        <WizardNavigation
                          disableNext={isSubmitting || !isValid}
                          editMode={!isCreatePage}
                          onNextClick={handleSubmit}
                          onBackClick={history.goBack}
                          onCancelClick={history.goBack}
                        />
                      </div>
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

MtoShipmentForm.defaultProps = {
  isCreatePage: false,
  match: { isExact: false, params: { moveID: '' } },
  history: { goBack: () => {}, push: () => {} },
  newDutyStationAddress: {
    city: '',
    state: '',
    postal_code: '',
  },
  mtoShipment: {
    id: '',
    customerRemarks: '',
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

export default MtoShipmentForm;
