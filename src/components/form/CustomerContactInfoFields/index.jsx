import React, { useRef } from 'react';
import { func, node, string } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Label, Fieldset } from '@trussworks/react-uswds';

import formStyles from 'styles/form.module.scss';
import TextField from 'components/form/fields/TextField';
import CheckboxField from 'components/form/fields/CheckboxField';

export const CustomerContactInfoFields = ({ legend, className, name, render }) => {
  const CustomerContactInfoFieldsUUID = useRef(uuidv4());

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-7">
              <TextField
                label="Best contact phone"
                id={`phone_${CustomerContactInfoFieldsUUID}`}
                name={`${name}.phone`}
                type="tel"
                maxLength="10"
              />
            </div>
          </div>
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-7">
              <TextField
                label="Alt. phone"
                labelHint="Optional"
                id={`alternatePhone_${CustomerContactInfoFieldsUUID}`}
                name={`${name}.alternatePhone`}
                type="tel"
                maxLength="10"
              />
            </div>
          </div>
          <TextField label="Personal email" id={`email_${CustomerContactInfoFieldsUUID}`} name={`${name}.email`} />
          <Label>Preferred contact method</Label>
          <div className={formStyles.radioGroup}>
            <CheckboxField
              id={`preferPhone_${CustomerContactInfoFieldsUUID}`}
              label="Phone"
              name={`${name}.preferPhone`}
            />
            <CheckboxField
              id={`preferEmail_ ${CustomerContactInfoFieldsUUID}`}
              label="Email"
              name={`${name}.preferEmail`}
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
  name: string.isRequired,
  render: func,
};

CustomerContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default CustomerContactInfoFields;
