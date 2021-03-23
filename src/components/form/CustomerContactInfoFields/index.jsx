import React, { useRef } from 'react';
import { func, node, string } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Label, Fieldset } from '@trussworks/react-uswds';

import formStyles from 'styles/form.module.scss';
import TextField from 'components/form/fields/TextField';
import CheckboxField from 'components/form/fields/CheckboxField';

export const CustomerContactInfoFields = ({ legend, className, render }) => {
  const CustomerContactInfoFieldsUUID = useRef(uuidv4());

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-7">
              <TextField
                label="Best contact phone"
                id={`telephone_${CustomerContactInfoFieldsUUID}`}
                name="telephone"
                type="tel"
                maxLength="10"
                required
              />
            </div>
          </div>
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-7">
              <TextField
                label="Alt. phone"
                labelHint="Optional"
                id={`secondaryTelephone_${CustomerContactInfoFieldsUUID}`}
                name="secondary_telephone"
                type="tel"
                maxLength="10"
              />
            </div>
          </div>
          <TextField
            label="Personal email"
            id={`personalEmail_${CustomerContactInfoFieldsUUID}`}
            name="personal_email"
            required
          />
          <Label>Preferred contact method</Label>
          <div className={formStyles.radioGroup}>
            <CheckboxField
              id={`phoneIsPreferred_${CustomerContactInfoFieldsUUID}`}
              label="Phone"
              name="phone_is_preferred"
            />
            <CheckboxField
              id={`emailIsPreferred_ ${CustomerContactInfoFieldsUUID}`}
              label="Email"
              name="email_is_preferred"
            />
          </div>
        </>,
      )}
    </Fieldset>
  );
};

CustomerContactInfoFields.propTypes = {
  legend: node,
  className: string,
  render: func,
};

CustomerContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default CustomerContactInfoFields;
