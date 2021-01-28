import React from 'react';
import { bool, shape, string, func, number } from 'prop-types';
import { Field } from 'formik';
import { Fieldset, Radio, Checkbox, Alert, FormGroup, Label, Textarea } from '@trussworks/react-uswds';

import styles from './MtoShipmentForm.module.scss';

import { shipmentForm } from 'content/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { DatePickerInput } from 'components/form/fields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { Form } from 'components/form/Form';
import Hint from 'components/Hint/index';
import { SimpleAddressShape } from 'types/address';
import { MtoShipmentFormValuesShape } from 'types/customerShapes';
import { validateDate } from 'utils/formikValidators';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import SectionWrapper from 'components/Customer/SectionWrapper';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';

const MtoShipmentFormFields = ({
  // formik data
  values,
  history,
  dirty,
  isValid,
  isSubmitting,
  // shipment-related data
  shipmentNumber,
  showPickupFields,
  showDeliveryFields,
  shipmentType,
  onUseCurrentResidenceChange,
  submitHandler,
  newDutyStationAddress,
  isCreatePage,
  serviceMember,
}) => {
  const isNTS = shipmentType === SHIPMENT_OPTIONS.NTS;
  const { hasDeliveryAddress } = values;

  const optionalLabel = <span className={styles.optional}>Optional</span>;
  return (
    <div className={styles.MTOShipmentForm}>
      <ShipmentTag shipmentType={shipmentType} shipmentNumber={shipmentNumber} />
      <h1>{shipmentForm.header[`${shipmentType}`]}</h1>
      <Alert type="info" noIcon>
        Remember: You can move {serviceMember.weight_allotment.total_weight_self} lbs total. You’ll be billed for any
        excess weight you move.
      </Alert>
      <Form className={styles.form}>
        {showPickupFields && (
          <>
            <SectionWrapper className={styles.formSection}>
              {showDeliveryFields && <h2>Pickup information</h2>}
              <Fieldset legend="Pickup date">
                <Field
                  as={DatePickerInput}
                  name="pickup.requestedDate"
                  label="Requested pickup date"
                  id="requestedPickupDate"
                  validate={validateDate}
                />
                <Hint id="pickupDateHint">
                  <p>
                    Movers will contact you to schedule the actual pickup date. That date should fall within 7 days of
                    your requested date. Tip: Avoid scheduling multiple shipments on the same day.
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
                      onChange={onUseCurrentResidenceChange}
                      id="useCurrentResidenceCheckbox"
                    />
                    {fields}
                    <Hint>
                      <p>
                        If you have more things at another pickup location, you can schedule a shipment for them later.
                      </p>
                    </Hint>
                  </>
                )}
                values={values.pickup.address}
              />

              <ContactInfoFields
                name="pickup.agent"
                legend={<div className={styles.legendContent}>Releasing agent {optionalLabel}</div>}
                values={values.pickup.agent}
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
                <Field
                  as={DatePickerInput}
                  name="delivery.requestedDate"
                  label="Requested delivery date"
                  id="requestedDeliveryDate"
                  validate={validateDate}
                />
                <Hint>
                  <p>
                    Shipments can take several weeks to arrive, depending on how far they’re going. Your movers will
                    contact you close to the date you select to coordinate delivery.
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
                    values={values.delivery.address}
                    render={(fields) => (
                      <>
                        {fields}
                        <Hint>
                          <p>
                            If you have more things to go to another destination, you can schedule a shipment for them
                            later.
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
                      {newDutyStationAddress.city}, {newDutyStationAddress.state} {newDutyStationAddress.postal_code}{' '}
                    </strong>
                    <br />
                    You can add the specific delivery address later, once you know it.
                  </p>
                )}
              </Fieldset>

              <ContactInfoFields
                name="delivery.agent"
                legend={<div className={styles.legendContent}>Receiving agent {optionalLabel}</div>}
                values={values.delivery.agent}
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
                  The moving company will find a storage facility approved by the government, and will move your
                  belongings there.
                </p>
                <p>
                  You’ll need to schedule an NTS release shipment to get your items back, most likely as part of a
                  future move.
                </p>
              </Fieldset>
            </SectionWrapper>
          </>
        )}

        <SectionWrapper className={styles.formSection}>
          <Fieldset legend={<div className={styles.legendContent}>Remarks {optionalLabel}</div>}>
            <Label for="customerRemarks">
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
              value={values.customerRemarks}
            />
            <Hint>
              <p>250 characters</p>
            </Hint>
          </Fieldset>
        </SectionWrapper>

        <Hint>
          <p>
            You can change details for your shipment when you talk to your move counselor or the person who’s your point
            of contact with the movers. You can also edit in MilMove up to 24 hours before your final pickup date.
          </p>
        </Hint>

        {!isCreatePage && (
          <div className={styles.formActions}>
            <WizardNavigation
              disableNext={isSubmitting || (!isValid && !dirty) || (isValid && !dirty)}
              editMode
              onNextClick={() => submitHandler(values)}
              onBackClick={history.goBack}
              onCancelClick={history.goBack}
            />
          </div>
        )}
      </Form>
    </div>
  );
};

MtoShipmentFormFields.propTypes = {
  // history params
  history: shape({
    goBack: func.isRequired,
  }).isRequired,

  // formik data
  values: MtoShipmentFormValuesShape,
  isSubmitting: bool,
  isValid: bool,
  dirty: bool,

  // customer data for pre-fill (& submit)
  isCreatePage: bool,
  newDutyStationAddress: SimpleAddressShape.isRequired,
  onUseCurrentResidenceChange: func.isRequired,
  submitHandler: func.isRequired,
  serviceMember: shape({
    weight_allotment: shape({
      total_weight_self: number,
    }),
  }).isRequired,

  // shipment-related data
  shipmentType: string.isRequired,
  showDeliveryFields: bool,
  showPickupFields: bool,
  shipmentNumber: string,
};

MtoShipmentFormFields.defaultProps = {
  // formik data
  values: {},
  isValid: false,
  isSubmitting: false,
  dirty: true,

  // shipment-related data
  shipmentNumber: '',
  isCreatePage: false,
  showDeliveryFields: true,
  showPickupFields: true,
};

export default MtoShipmentFormFields;
