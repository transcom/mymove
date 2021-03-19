import React, { useRef } from 'react';
import { func, node, string } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField';

export const BackupContactInfoFields = ({ legend, className, name, render }) => {
  const backupContactInfoFieldsUUID = useRef(uuidv4());

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField label="Name" id={`name_${backupContactInfoFieldsUUID}`} name={`${name}.name`} />
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-7">
              <TextField label="Email" id={`email_${backupContactInfoFieldsUUID}`} name={`${name}.email`} />
            </div>
          </div>
          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-4">
              <TextField
                label="Phone"
                id={`phone_${backupContactInfoFieldsUUID}`}
                name={`${name}.phone`}
                type="tel"
                maxLength="10"
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
  name: string.isRequired,
  render: func,
};

BackupContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default BackupContactInfoFields;
