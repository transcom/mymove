import React from 'react';
import { Field } from 'formik';
import { Fieldset } from '@trussworks/react-uswds';

import { DatePickerInput } from '../../form/fields';
import { ContactInfoFields } from '../../form/ContactInfoFields/ContactInfoFields';
import { AddressFields } from '../../form/AddressFields/AddressFields';
import Checkbox from 'shared/Checkbox';
import { validateDate } from 'utils/formikValidators';

export const PickupDetails = ({ fieldsetClasses, values, useCurrentResidence, onCurrentResidenceChange }) => {
  return (
    <div>
      <Fieldset legend="Pickup date" className={fieldsetClasses}>
        <Field
          as={DatePickerInput}
          name="requestedPickupDate"
          label="Requested pickup date"
          id="requestedPickupDate"
          value={values.requestedDate}
          validate={validateDate}
        />
        <span className="usa-hint" id="pickupDateHint">
          Your movers will confirm this date or one shortly before or after.
        </span>
      </Fieldset>
      <AddressFields
        name="pickupAddress"
        legend="Pickup location"
        className={fieldsetClasses}
        renderExistingAddressCheckbox={() => (
          <div className="margin-y-2">
            <Checkbox
              data-testid="useCurrentResidence"
              label="Use my current residence address"
              name="useCurrentResidence"
              checked={useCurrentResidence}
              onChange={() => onCurrentResidenceChange(values)}
            />
          </div>
        )}
        values={values.address}
      />
      <ContactInfoFields
        name="releasingAgent"
        legend="Releasing agent"
        className={fieldsetClasses}
        subtitle="Who can allow the movers to take your stuff if you're not there?"
        subtitleClassName="margin-y-2"
        values={values.agent}
      />
    </div>
  );
}

export default PickupDetails;
