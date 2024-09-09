import React, { useRef } from 'react';
import { func, node, string } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Label, Fieldset } from '@trussworks/react-uswds';
import classnames from 'classnames';

import formStyles from 'styles/form.module.scss';
import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { CheckboxField } from 'components/form/fields';

export const CustomerContactInfoFields = ({ legend, className, render, labelHint: labelHintProp }) => {
  const CustomerContactInfoFieldsUUID = useRef(uuidv4());

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-7">
              <MaskedTextField
                label="Best contact phone"
                id={`telephone_${CustomerContactInfoFieldsUUID.current}`}
                name="telephone"
                labelHint={labelHintProp}
                type="tel"
                minimum="12"
                mask="000{-}000{-}0000"
                required
              />
            </div>
          </div>
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-7">
              <MaskedTextField
                label="Alt. phone"
                labelHint={labelHintProp ? null : 'Optional'}
                id={`secondaryTelephone_${CustomerContactInfoFieldsUUID.current}`}
                name="secondary_telephone"
                type="tel"
                minimum="12"
                mask="000{-}000{-}0000"
              />
            </div>
          </div>
          <TextField
            label="Personal email"
            id={`personalEmail_${CustomerContactInfoFieldsUUID.current}`}
            name="personal_email"
            labelHint={labelHintProp}
            required
          />
          <Label hint={labelHintProp}>Preferred contact method</Label>
          <div className={classnames(formStyles.radioGroup, formStyles.customerPreferredContact)}>
            <CheckboxField
              id={`phoneIsPreferred_${CustomerContactInfoFieldsUUID.current}`}
              label="Phone"
              name="phone_is_preferred"
            />
            <CheckboxField
              id={`emailIsPreferred_ ${CustomerContactInfoFieldsUUID.current}`}
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
