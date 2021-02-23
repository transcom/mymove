import React from 'react';
import PropTypes from 'prop-types';
import { Label, TextInput } from '@trussworks/react-uswds';

import { Hint } from '../Hint';

export const FormGroupWarning = ({ inputLabel, warningMessage }) => {
  <div>
    <Label htmlFor="input-type-text">{inputLabel}</Label>
    <TextInput id="input-type-text" name="input-type-text" type="text" validationStatus="warning" />
    <em>
      <Hint>{warningMessage}</Hint>
    </em>
  </div>;
};

FormGroupWarning.propTypes = {
  inputLabel: PropTypes.string,
  warningMessage: PropTypes.string,
};

export default FormGroupWarning;
