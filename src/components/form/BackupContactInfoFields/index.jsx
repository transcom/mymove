import React, { useRef } from 'react';
import { func, node, string } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField';

export const BackupContactInfoFields = ({ legend, className, render }) => {
  const backupContactInfoFieldsUUID = useRef(uuidv4());

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField label="Name" id={`name_${backupContactInfoFieldsUUID}`} name="name" required />
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-7">
              <TextField label="Email" id={`email_${backupContactInfoFieldsUUID}`} name="email" required />
            </div>
          </div>
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-4">
              <MaskedTextField
                label="Phone"
                id={`phone_${backupContactInfoFieldsUUID}`}
                name="telephone"
                type="tel"
                minimum="12"
                mask="000{-}000{-}0000"
                required
              />
            </div>
          </div>
        </>,
      )}
    </Fieldset>
  );
};

BackupContactInfoFields.propTypes = {
  legend: node,
  className: string,
  render: func,
};

BackupContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default BackupContactInfoFields;
