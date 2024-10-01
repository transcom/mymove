import React, { useRef } from 'react';
import { func, node, string } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

export const BackupContactInfoFields = ({ name, legend, className, render, labelHint: labelHintProp }) => {
  const backupContactInfoFieldsUUID = useRef(uuidv4());

  let nameFieldName = 'name';
  let emailFieldName = 'email';
  let phoneFieldName = 'telephone';

  if (name !== '') {
    nameFieldName = `${name}.name`;
    emailFieldName = `${name}.email`;
    phoneFieldName = `${name}.telephone`;
  }

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField
            label="Name"
            id={`name_${backupContactInfoFieldsUUID.current}`}
            name={nameFieldName}
            required
            labelHint={labelHintProp}
          />
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-7">
              <TextField
                label="Email"
                id={`email_${backupContactInfoFieldsUUID.current}`}
                name={emailFieldName}
                required
                labelHint={labelHintProp}
              />
            </div>
          </div>
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-4">
              <MaskedTextField
                label="Phone"
                id={`phone_${backupContactInfoFieldsUUID.current}`}
                name={phoneFieldName}
                type="tel"
                minimum="12"
                mask="000{-}000{-}0000"
                required
                labelHint={labelHintProp}
              />
            </div>
          </div>
        </>,
      )}
    </Fieldset>
  );
};

BackupContactInfoFields.propTypes = {
  name: string,
  legend: node,
  className: string,
  render: func,
};

BackupContactInfoFields.defaultProps = {
  name: '',
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default BackupContactInfoFields;
