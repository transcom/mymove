import React from 'react';
import { func, node, string } from 'prop-types';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField/TextField';

export const OktaInfoFields = ({ legend, className, render }) => {
  const usernameFieldName = 'oktaUsername';
  const emailFieldName = 'oktaEmail';
  const firstNameFieldName = 'oktaFirstName';
  const lastNameFieldName = 'oktaLastName';
  const edipiFieldName = 'oktaEdipi';

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField isDisabled label="Okta Username" name={usernameFieldName} id="oktaUsername" required />
          <TextField label="Okta Email" name={emailFieldName} id="oktaEmail" required />
          <TextField label="First Name" name={firstNameFieldName} id="oktaFirstName" required />
          <TextField label="Last Name" name={lastNameFieldName} id="oktaLastName" required />
          <TextField
            label="DoD ID number | EDIPI"
            name={edipiFieldName}
            id="oktaEdipi"
            maxLength="10"
            inputMode="numeric"
          />
        </>,
      )}
    </Fieldset>
  );
};

OktaInfoFields.propTypes = {
  legend: node,
  className: string,
  render: func,
};

OktaInfoFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default OktaInfoFields;
