import React from 'react';
import PropTypes from 'prop-types';
import { Label, TextInput } from '@trussworks/react-uswds';

export const FormGroupWarning = ({ inputLabel, warningMessage }) => {
  <div>
    <Label htmlFor="input-type-text">{inputLabel}</Label>
    <TextInput id="input-type-text" name="input-type-text" type="text" validationStatus="warning" />
    <em>
      <p className="usa-hint">{warningMessage}</p>
    </em>
  </div>;
};

FormGroupWarning.propTypes = {
  inputLabel: PropTypes.string,
  warningMessage: PropTypes.string,
};

export default FormGroupWarning;
