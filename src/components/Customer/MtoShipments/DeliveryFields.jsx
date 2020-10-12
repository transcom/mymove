import React from 'react';
import { Fieldset, Radio, Label } from '@trussworks/react-uswds';
import { string, bool, shape, func } from 'prop-types';

import { simpleAddressShape, fullAddressShape, agentShape } from './propShapes';

import { DatePickerInput } from 'components/form/fields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { validateDate } from 'utils/formikValidators';

export const DeliveryFields = ({
  fieldsetClasses,
  values,
  hasDeliveryAddress,
  onHasAddressChange,
  newDutyStationAddress,
}) => {
  return (
    <div>
      <Fieldset legend="Delivery date" className={fieldsetClasses}>
        <DatePickerInput
          name="requestedDeliveryDate"
          label="Requested delivery date"
          id="requestedDeliveryDate"
          value={values.requestedDate}
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
            onChange={onHasAddressChange}
            checked={hasDeliveryAddress}
          />
          <Radio
            id="no-delivery-address"
            label="No"
            name="hasDeliveryAddress"
            checked={!hasDeliveryAddress}
            onChange={onHasAddressChange}
          />
        </div>
        {hasDeliveryAddress ? (
          <AddressFields name="deliveryAddress" values={values.address} />
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
        name="receivingAgent"
        legend="Receiving agent"
        className={fieldsetClasses}
        subtitle="Who can take delivery for you if the movers arrive and you're not there?"
        subtitleClassName="margin-y-2"
        values={values.agent}
      />
    </div>
  );
};

DeliveryFields.propTypes = {
  fieldsetClasses: string,
  hasDeliveryAddress: bool,
  onHasAddressChange: func,
  newDutyStationAddress: simpleAddressShape.isRequired,
  values: shape({
    address: fullAddressShape,
    agent: agentShape,
    requestedDate: string,
  }),
};

DeliveryFields.defaultProps = {
  fieldsetClasses: '',
  hasDeliveryAddress: false,
  onHasAddressChange: () => {},
  values: {},
};

export default DeliveryFields;
