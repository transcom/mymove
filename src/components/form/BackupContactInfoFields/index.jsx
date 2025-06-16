import React, { useRef } from 'react';
import { func, node, string } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import { requiredAsteriskMessage } from '../RequiredAsterisk';

import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

export const BackupContactInfoFields = ({ name, legend, className, render, labelHint: labelHintProp }) => {
  const backupContactInfoFieldsUUID = useRef(uuidv4());

  let firstNameFieldName = 'firstName';
  let lastNameFieldName = 'lastName';
  let emailFieldName = 'email';
  let phoneFieldName = 'telephone';

  if (name !== '') {
    firstNameFieldName = `${name}.firstName`;
    lastNameFieldName = `${name}.lastName`;
    emailFieldName = `${name}.email`;
    phoneFieldName = `${name}.telephone`;
  }

  const showRequiredAsterisk = labelHintProp !== 'Optional';

  return (
    <Fieldset legend={legend} className={className}>
      {requiredAsteriskMessage}
      {render(
        <>
          <TextField
            label="First Name"
            id={`firstName_${backupContactInfoFieldsUUID.current}`}
            name={firstNameFieldName}
            required
            labelHint={labelHintProp}
          />
          <TextField
            label="Last Name"
            id={`lastName_${backupContactInfoFieldsUUID.current}`}
            name={lastNameFieldName}
            required
            showRequiredAsterisk={showRequiredAsterisk}
          />
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-7">
              <TextField
                label="Email"
                id={`email_${backupContactInfoFieldsUUID.current}`}
                name={emailFieldName}
                required
                showRequiredAsterisk={showRequiredAsterisk}
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
                showRequiredAsterisk={showRequiredAsterisk}
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
