import React, { useRef } from 'react';
import { func, node, string } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Label, Fieldset, ErrorMessage } from '@trussworks/react-uswds';
import { useFormikContext } from 'formik';
import classnames from 'classnames';

import formStyles from 'styles/form.module.scss';
import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { CheckboxField } from 'components/form/fields';

export const CustomerAltContactInfoFields = ({ legend, className, render }) => {
  const CustomerAltContactInfoFieldsUUID = useRef(uuidv4());
  const { errors } = useFormikContext();

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <div className="grid-row grid-gap">
            <div className="grid-col-6">
              <TextField label="First name" name="firstName" id="firstName" required />
            </div>
            <div className="grid-col-6">
              <TextField label="Middle name" name="middleName" id="middleName" labelHint="Optional" />
            </div>
            <div className="grid-col-6">
              <TextField label="Last name" name="lastName" id="lastName" required />
            </div>
            <div className="grid-col-6">
              <TextField label="Suffix" name="suffix" id="suffix" labelHint="Optional" />
            </div>
          </div>
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-6">
              <MaskedTextField
                label="Phone"
                id={`customerTelephone_${CustomerAltContactInfoFieldsUUID.current}`}
                name="customerTelephone"
                type="tel"
                minimum="12"
                mask="000{-}000{-}0000"
                required
              />
            </div>
            <div className="mobile-lg:grid-col-6">
              <MaskedTextField
                label="Alternate Phone"
                id={`secondaryPhone_${CustomerAltContactInfoFieldsUUID.current}`}
                name="secondaryPhone"
                type="tel"
                minimum="12"
                mask="000{-}000{-}0000"
                labelHint="Optional"
              />
            </div>
          </div>

          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-6">
              <TextField
                label="Email"
                id={`customerEmail_${CustomerAltContactInfoFieldsUUID.current}`}
                name="customerEmail"
                required
              />
            </div>
            <div className="grid-row grid-gap">
              <Label>Preferred contact method</Label>
              {errors.preferredContactMethod ? <ErrorMessage>{errors.preferredContactMethod}</ErrorMessage> : null}
              <div className={classnames(formStyles.radioGroup, formStyles.customerPreferredContact)}>
                <CheckboxField
                  id={`phoneIsPreferred_${CustomerAltContactInfoFieldsUUID.current}`}
                  label="Phone"
                  name="phoneIsPreferred"
                />
                <CheckboxField
                  id={`emailIsPreferred_ ${CustomerAltContactInfoFieldsUUID.current}`}
                  label="Email"
                  name="emailIsPreferred"
                />
              </div>
            </div>
          </div>
        </>,
      )}
    </Fieldset>
  );
};

CustomerAltContactInfoFields.propTypes = {
  legend: node,
  className: string,
  render: func,
};

CustomerAltContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default CustomerAltContactInfoFields;
