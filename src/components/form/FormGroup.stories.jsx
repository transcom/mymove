import React from 'react';
import { FormGroup, Label, TextInput, ErrorMessage } from '@trussworks/react-uswds';

export default {
  title: 'Components/Form',
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
    <TextInput id="input-type-text" name="input-type-text" type="text" />
    <p className="usa-hint">
      This TAC does not appear in TGET, so it might not be valid. Make sure it matches what&apos;s on the orders before
      you continue.
    </p>
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
