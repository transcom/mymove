import React from 'react';
import { bool, shape, string, func, number } from 'prop-types';
import { Field } from 'formik';
import { Button, Fieldset, Label, Radio, Checkbox, Alert } from '@trussworks/react-uswds';

import styles from './MtoShipmentForm.module.scss';

import { shipmentForm } from 'content/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { DatePickerInput, TextInput } from 'components/form/fields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { Form } from 'components/form/Form';
import Divider from 'shared/Divider';
// import Fieldset from 'shared/Fieldset';
import Hint from 'shared/Hint';
import { SimpleAddressShape } from 'types/address';
import { MtoDisplayOptionsShape, MtoShipmentFormValuesShape } from 'types/customerShapes';
import { validateDate } from 'utils/formikValidators';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';

const MtoShipmentFormFields = ({
  // formik data
  values,
  history,
  dirty,
  isValid,
  isSubmitting,
  // shipment-related data
  shipmentNumber,
  displayOptions,
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
        Remember: You can move {serviceMember.weight_allotment.total_weight_self} lbs total. You&rsquo;ll be billed for
        any excess weight you move.
      </Alert>
      <Form className={styles.HHGDetailsForm}>
        {displayOptions.showPickupFields && (
          <>
            <Fieldset legend="Pickup date">
              <Field
                as={DatePickerInput}
                name="pickup.requestedDate"
                label="Requested pickup date"
                id="requestedPickupDate"
                validate={validateDate}
              />
              <Hint id="pickupDateHint">
                Movers will contact you to schedule the actual pickup date. That date should fall within 7 days of your
                requested date. Tip: Avoid scheduling multiple shipments on the same day.
              </Hint>
            </Fieldset>

            <Divider />

            <AddressFields
              name="pickup.address"
              legend="Pickup location"
              render={(fields) => (
                <>
                  <Checkbox
                    data-testid="useCurrentResidence"
                    label="Use my current residence address"
                    name="useCurrentResidence"
                    onChange={onUseCurrentResidenceChange}
                    id="useCurrentResidenceCheckbox"
                  />
                  {fields}
                  <Hint>If you have more things at another pickup location, you can schedule for them later.</Hint>
                </>
              )}
              values={values.pickup.address}
            />

            <Divider />

            <ContactInfoFields
              name="pickup.agent"
              legend={<>Releasing agent {optionalLabel}</>}
              subtitle="Who can allow the movers to take your stuff if you're not there?"
              values={values.pickup.agent}
            />
          </>
        )}

        {displayOptions.showDeliveryFields && (
          <>
            <Fieldset legend="Delivery date">
              <Field
                as={DatePickerInput}
                name="delivery.requestedDate"
                label="Requested delivery date"
                id="requestedDeliveryDate"
                validate={validateDate}
              />
              <Hint>
                Shipments can take several weeks to arrive, depending on how far they&rsquo;re going. Your movers will
                contact you close to the date you select to coordinate delivery.
              </Hint>
            </Fieldset>

            <Divider />

            <Fieldset legend="Delivery location">
              <Label>Do you know your delivery address?</Label>
              <div>
                <Field
                  as={Radio}
                  id="has-delivery-address"
                  label="Yes"
                  name="hasDeliveryAddress"
                  value="yes"
                  checked={hasDeliveryAddress === 'yes'}
                />
                <Field
                  as={Radio}
                  id="no-delivery-address"
                  label="No"
                  name="hasDeliveryAddress"
                  value="no"
                  checked={hasDeliveryAddress === 'no'}
                />
              </div>
              {hasDeliveryAddress === 'yes' ? (
                <AddressFields name="delivery.address" values={values.delivery.address} />
              ) : (
                <>
                  <p>
                    We can use the zip of your new duty station.
                    <br />
                    <strong>
                      {newDutyStationAddress.city}, {newDutyStationAddress.state} {newDutyStationAddress.postal_code}{' '}
                    </strong>
                  </p>
                </>
              )}
            </Fieldset>

            <Divider />

            <ContactInfoFields
              name="delivery.agent"
              legend={<>Receiving agent {optionalLabel}</>}
              subtitle="Who can take delivery for you if the movers arrive and you're not there?"
              values={values.delivery.agent}
            />
          </>
        )}

        {isNTS && (
          <>
            <Divider />

            <Fieldset legend="What you can expect" data-testid="nts-what-to-expect">
              <p>
                The moving company will find a storage facility approved by the government, and will move your
                belongings there.
              </p>
              <p>
                You’ll need to schedule an NTS release shipment to get your items back, most likely as part of a future
                move.
              </p>
            </Fieldset>
          </>
        )}

        <Divider />

        <Fieldset legend={<>Remarks {optionalLabel}</>}>
          <div>Is there anything special about this shipment that the movers should know?</div>
          <div className={`${styles['hhg-examples-container']}`}>
            <strong>Examples</strong>
            <ul>
              <li>Things that might need special handling</li>
              <li>Access info for a location</li>
              <li>Weapons or alcohol</li>
            </ul>
          </div>

          <TextInput
            label="Anything else you would like us to know?"
            labelHint="(optional)"
            data-testid="remarks"
            name="customerRemarks"
            className={`${styles.remarks}`}
            placeholder="You don&rsquo;t need to list all belongings here. Your mover will get those details later."
            id="customerRemarks"
            maxLength={250}
            value={values.customerRemarks}
          />
          <Hint>250 characters</Hint>
        </Fieldset>

        <Divider />

        <Hint>
          You can change details for your shipment when you talk to your move counselor or the person who’s your point
          of contact with the movers. You can also edit in MilMove up to 24 hours before your final pickup date.
        </Hint>

        {!isCreatePage && (
          <div style={{ display: 'flex', flexDirection: 'column' }}>
            <Button
              disabled={isSubmitting || (!isValid && !dirty) || (isValid && !dirty)}
              onClick={() => submitHandler(values)}
            >
              <span>Save</span>
            </Button>
            <Button className={`${styles['cancel-button']}`} onClick={history.goBack}>
              <span>Cancel</span>
            </Button>
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
  displayOptions: MtoDisplayOptionsShape.isRequired,
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
};

export default MtoShipmentFormFields;
