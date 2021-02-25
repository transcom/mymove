import React from 'react';

import {FormGroup, Label, TextInput } from '@trussworks/react-uswds';

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
  </FormGroup>
);

export const FormGroupWithError= () => (
  <FormGroup error>
    <Label htmlFor="input-type-text">Text input label</Label>
    <TextInput id="input-type-text" name="input-type-text" type="text" />
  </FormGroup>
);

