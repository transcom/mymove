import React from 'react';
import FormGroupWarning from './FormGroupWarning';

export default {
  title: 'Components/Form',
  component: FormGroupWarning,
};

export const FormFieldWithWarning = () => (
  <FormGroupWarning inputLabel="Warning Message" warningMessage="Helpful warning message text go here" />
);
