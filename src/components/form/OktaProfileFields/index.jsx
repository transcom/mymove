import React from 'react';
import { func, node, string } from 'prop-types';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField/TextField';

export const OktaProfileFields = ({ legend, className, render }) => {
  const usernameFieldName = 'username';
  const emailFieldName = 'email';
  const fNameFieldName = 'fName';
  const lNameFieldName = 'lName';
  const edipiFieldName = 'edipi';

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField label="Okta Username" name={usernameFieldName} required />
          <TextField label="Okta Email" name={emailFieldName} required />
          <TextField label="First Name" name={fNameFieldName} required />
          <TextField label="Last Name" name={lNameFieldName} required />
          <TextField label="DoD ID Number" name={edipiFieldName} required />
        </>,
      )}
    </Fieldset>
  );
};

OktaProfileFields.propTypes = {
  legend: node,
  className: string,
  render: func,
};

OktaProfileFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default OktaProfileFields;
