import React from 'react';
import { bool, shape, string, func } from 'prop-types';
import { Field } from 'formik';
import { Button, /* Fieldset, */ Label, Radio, Checkbox } from '@trussworks/react-uswds';

import styles from './MtoShipmentForm.module.scss';

import { DatePickerInput, TextInput } from 'components/form/fields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { Form } from 'components/form/Form';
import Divider from 'shared/Divider';
import Fieldset from 'shared/Fieldset';
import Hint from 'shared/Hint';
import { SimpleAddressShape } from 'types/address';
import { MtoDisplayOptionsShape, MtoShipmentFormValuesShape } from 'types/customerShapes';
import { validateDate } from 'utils/formikValidators';

const hhgFormHeader = 'When and where can the movers pick up and deliver this shipment?';
const ntsFormHeader = 'Where and when should the movers pick up your things going into storage?';
const ntsrFormHeader = 'Where and when should the movers release your things from storage?';

const MtoShipmentFormFields = ({
  // formik data
  values,
  history,
  dirty,
  isValid,
  isSubmitting,
  // shipment-related data
  displayOptions,
  shipmentNumber,
  onUseCurrentResidenceChange,
  submitHandler,
  newDutyStationAddress,
  isCreatePage,
  serviceMember,
}) => {
  const isHHG = displayOptions.displayName === 'HHG';
  const isNTS = displayOptions.displayName === 'NTS';
  const isNTSR = displayOptions.displayName === 'NTS-R';
  const { hasDeliveryAddress } = values;

  return (
    <>
      <div className={`margin-top-2 ${styles['hhg-label']}`}>{`${displayOptions.displayName} ${shipmentNumber}`}</div>
      <h1 className="margin-top-1">
        {isHHG && hhgFormHeader}
        {isNTS && ntsFormHeader}
        {isNTSR && ntsrFormHeader}
      </h1>
      <p>
        Remember: You can move {serviceMember.weight_allotment.total_weight_self} lbs total. You&rsquo;ll be billed for
        any excess weight you move.
      </p>
      <Form className={styles.HHGDetailsForm}>
        {displayOptions.showPickupFields && (
          <div>
            <Fieldset legend="Pickup date" className="margin-top-4">
              <Field
                as={DatePickerInput}
                name="pickup.requestedDate"
                label="Requested pickup date"
                labelClassName={`margin-top-2 ${styles['small-bold']}`}
                id="requestedPickupDate"
                validate={validateDate}
              />
            </Fieldset>
            <Hint className="margin-top-1" id="pickupDateHint">
              Movers will contact you to schedule the actual pickup date. That date should fall within 7 days of your
              requested date. Tip: Avoid scheduling multiple shipments on the same day.
            </Hint>
            <Divider className="margin-top-4 margin-bottom-4" />
            <AddressFields
              className="margin-bottom-3"
              name="pickup.address"
              legend="Pickup location"
              renderExistingAddressCheckbox={() => (
                <div className="margin-y-2">
                  <Checkbox
                    data-testid="useCurrentResidence"
                    label="Use my current residence address"
                    name="useCurrentResidence"
                    onChange={onUseCurrentResidenceChange}
                    id="useCurrentResidenceCheckbox"
                  />
                </div>
              )}
              values={values.pickup.address}
            />
            <Hint>If you have more things at another pickup location, you can schedule for them later.</Hint>
            <hr className="margin-top-4 margin-bottom-4" />
            <ContactInfoFields
              className="margin-bottom-5"
              name="pickup.agent"
              legend="Releasing agent"
              hintText="Optional"
              subtitle="Who can allow the movers to take your stuff if you're not there?"
              subtitleClassName="margin-top-3"
              values={values.pickup.agent}
            />
          </div>
        )}
        {displayOptions.showDeliveryFields && (
          <>
            <Divider className="margin-bottom-6" />
            <Fieldset legend="Delivery date">
              <Field
                as={DatePickerInput}
                name="delivery.requestedDate"
                label="Requested delivery date"
                labelClassName={`${styles['small-bold']}`}
                id="requestedDeliveryDate"
                validate={validateDate}
              />
              <Hint className="margin-top-1">
                Shipments can take several weeks to arrive, depending on how far they&rsquo;re going. Your movers will
                contact you close to the date you select to coordinate delivery.
              </Hint>
            </Fieldset>
            <Divider className="margin-top-4 margin-bottom-4" />
            <Fieldset legend="Delivery location">
              <Label className="margin-top-3 margin-bottom-1">Do you know your delivery address?</Label>
              <div className="display-flex margin-top-1">
                <Field
                  as={Radio}
                  className="margin-right-3"
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
                  <p className="margin-top-2">
                    We can use the zip of your new duty station.
                    <br />
                    <strong>
                      {newDutyStationAddress.city}, {newDutyStationAddress.state} {newDutyStationAddress.postal_code}{' '}
                    </strong>
                  </p>
                </>
              )}
            </Fieldset>
            <Divider className="margin-top-4 margin-bottom-4" />
            <ContactInfoFields
              name="delivery.agent"
              legend="Receiving agent"
              hintText="Optional"
              subtitle="Who can take delivery for you if the movers arrive and you're not there?"
              subtitleClassName="margin-top-3"
              values={values.delivery.agent}
            />
          </>
        )}
        <Divider className="margin-top-4 margin-bottom-4" />
        <Fieldset legend="What you can expect">
          {isNTS && (
            <>
              <p>
                The moving company will find a storage facility approved by the government, and will move your
                belongings there.
              </p>
              <p>
                You&rsquo;ll need to schedule an NTS release shipment to get your items back, most likely as part of a
                future move.
              </p>
            </>
          )}
        </Fieldset>
        <Divider className="margin-top-4 margin-bottom-4" />
        <Fieldset hintText="Optional" legend="Remarks">
          <div className={`${styles['small-bold']} margin-top-3 margin-bottom-1`}>
            Is there anything special about this shipment that the movers should know?
          </div>
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
          <Hint className="margin-bottom-2">250 characters</Hint>
        </Fieldset>
        <Divider className="margin-top-6 margin-bottom-3" />
        <Hint className="margin-bottom-2">
          You can change details for your shipment when you talk to your move counselor or the person who&rsquo;s your
          point of contact with the movers. You can also edit in MilMove up to 24 hours before your final pickup date.
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
  isCreatePage: bool,
  newDutyStationAddress: SimpleAddressShape.isRequired,
  onUseCurrentResidenceChange: func.isRequired,
  submitHandler: func.isRequired,
  serviceMember: shape({
    weight_allotment: shape({
      total_weight_self: string,
    }),
  }).isRequired,

  // shipment-related data
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
