import React from 'react';
import { bool, shape, string, func } from 'prop-types';
import { Field } from 'formik';
import { Button, /* Fieldset, */ Label, Radio } from '@trussworks/react-uswds';

import styles from './MtoShipmentForm.module.scss';

import { DatePickerInput, TextInput } from 'components/form/fields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { Form } from 'components/form/Form';
import Checkbox from 'shared/Checkbox';
import Fieldset from 'shared/Fieldset';
import { SimpleAddressShape } from 'types/address';
import { MtoDisplayOptionsShape, MtoShipmentFormValuesShape } from 'types/customerShapes';
import { validateDate } from 'utils/formikValidators';

const MtoShipmentFormFields = (
  // formik data
  values,
  history,
  dirty,
  isValid,
  isSubmitting,
  // shipment-related data
  displayOptions,
  shipmentNumber,
  hasDeliveryAddress,
  onHasDeliveryAddressChange,
  useCurrentResidence,
  onUseCurrentResidenceChange,
  submitHandler,
  newDutyStationAddress,
  isCreateForm,
) => {
  const fieldsetClasses = 'margin-top-2';

  return (
    <>
      <div className={`margin-top-2 ${styles['hhg-label']}`}>{`${displayOptions.displayName} ${shipmentNumber}`}</div>
      <h1 className="margin-top-1">When and where can the movers pick up and deliver this shipment?</h1>
      <Form className={styles.HHGDetailsForm}>
        {displayOptions.showPickupFields && (
          <div>
            <Fieldset legend="Pickup date" className={fieldsetClasses}>
              <Field
                as={DatePickerInput}
                name="pickup.requestedDate"
                label="Requested pickup date"
                id="requestedPickupDate"
                value={values.pickup.requestedDate}
                validate={validateDate}
              />
              <span className="usa-hint" id="pickupDateHint">
                Your movers will confirm this date or one shortly before or after.
              </span>
            </Fieldset>

            <AddressFields
              name="pickup.address"
              legend="Pickup location"
              className={fieldsetClasses}
              renderExistingAddressCheckbox={() => (
                <div className="margin-y-2">
                  <Checkbox
                    data-testid="useCurrentResidence"
                    label="Use my current residence address"
                    name="useCurrentResidence"
                    checked={useCurrentResidence}
                    onChange={() => onUseCurrentResidenceChange(values)}
                  />
                </div>
              )}
              values={values.pickup.address}
            />
            <ContactInfoFields
              name="pickup.agent"
              legend="Releasing agent"
              className={fieldsetClasses}
              subtitle="Who can allow the movers to take your stuff if you're not there?"
              subtitleClassName="margin-y-2"
              values={values.pickup.agent}
            />
          </div>
        )}
        {displayOptions.showDeliveryFields && (
          <div>
            <Fieldset legend="Delivery date" className={fieldsetClasses}>
              <DatePickerInput
                name="delivery.requestedDate"
                label="Requested delivery date"
                id="requestedDeliveryDate"
                value={values.delivery.requestedDate}
                validate={validateDate}
              />
              <small className="usa-hint" id="deliveryDateHint">
                Your movers will confirm this date or one shortly before or after.
              </small>
            </Fieldset>
            <Fieldset legend="Delivery location" className={fieldsetClasses}>
              <Label>Do you know your delivery address?</Label>
              <div className="display-flex margin-top-1">
                <Radio
                  id="has-delivery-address"
                  label="Yes"
                  name="hasDeliveryAddress"
                  onChange={onHasDeliveryAddressChange}
                  checked={hasDeliveryAddress}
                />
                <Radio
                  id="no-delivery-address"
                  label="No"
                  name="hasDeliveryAddress"
                  checked={!hasDeliveryAddress}
                  onChange={onHasDeliveryAddressChange}
                />
              </div>
              {hasDeliveryAddress ? (
                <AddressFields name="destinationAddress" values={values.delivery.address} />
              ) : (
                <>
                  <div>
                    <p className={fieldsetClasses}>
                      We can use the zip of your new duty station.
                      <br />
                      <strong>
                        {newDutyStationAddress.city}, {newDutyStationAddress.state} {newDutyStationAddress.postal_code}{' '}
                      </strong>
                    </p>
                  </div>
                </>
              )}
            </Fieldset>
            <ContactInfoFields
              name="delivery.agent"
              legend="Receiving agent"
              className={fieldsetClasses}
              subtitle="Who can take delivery for you if the movers arrive and you're not there?"
              subtitleClassName="margin-y-2"
              values={values.delivery.agent}
            />
          </div>
        )}
        <Fieldset legend="Remarks" className={fieldsetClasses}>
          <TextInput
            label="Anything else you would like us to know?"
            labelHint="(optional)"
            data-testid="remarks"
            name="customerRemarks"
            id="customerRemarks"
            maxLength={1500}
            value={values.customerRemarks}
          />
        </Fieldset>
        {!isCreateForm && (
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
    </>
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
  isCreateForm: bool.isRequired,
  newDutyStationAddress: SimpleAddressShape.isRequired,
  onHasDeliveryAddressChange: func.isRequired,
  onUseCurrentResidenceChange: func.isRequired,
  submitHandler: func.isRequired,

  // shipment-related data
  displayOptions: MtoDisplayOptionsShape.isRequired,
  hasDeliveryAddress: bool.isRequired,
  useCurrentResidence: bool.isRequired,
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
};

export default { MtoShipmentFormFields };
