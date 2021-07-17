import React from 'react';
import { FormGroup, Label, TextInput, ErrorMessage } from '@trussworks/react-uswds';

import Hint from 'components/Hint';

export default {
  title: 'Components/FormGroup',
};

export const FormGroupDefaultState = () => (
  <FormGroup>
    <Label htmlFor="input-type-text">Text input label</Label>
    <TextInput id="input-type-text" name="input-type-text" type="text" />
  </FormGroup>
);
export const FormGroupWithWarning = () => (
  <FormGroup className="warning">
    <Label htmlFor="input-type-text">Text input label</Label>
    <Hint>
      This TAC does not appear in TGET, so it might not be valid. Make sure it matches what&apos;s on the orders before
      you continue.
    </Hint>
    <TextInput id="input-type-text" name="input-type-text" type="text" />
  </FormGroup>
);

export const FormGroupWithError = () => (
  <FormGroup error>
    <Label htmlFor="input-type-text" error>
      Text input label
    </Label>
    <ErrorMessage>Helpful error message</ErrorMessage>
    <TextInput id="input-type-text" name="input-type-text" type="text" validationStatus="error" />
  </FormGroup>
);
