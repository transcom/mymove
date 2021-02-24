import React from 'react';
import PropTypes from 'prop-types';
import { Label, TextInput } from '@trussworks/react-uswds';

import styles from './FormGroupWarning.module.scss';

export const FormGroupWarning = ({ inputLabel, warningMessage }) => {
  return (
    <div>
      <Label htmlFor="input-type-text">{inputLabel}</Label>
      <TextInput id="input-type-text" name="input-type-text" type="text" validationStatus="warning" />
      <em>
        <p className="usa-hint">{warningMessage}</p>
      </em>
    </div>
  );
};

FormGroupWarning.propTypes = {
  inputLabel: PropTypes.string.isRequired,
  warningMessage: PropTypes.string.isRequired,
};

export default FormGroupWarning;
