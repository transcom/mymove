import React from 'react';

import FormGroupWarning from './FormGroupWarning';

export default {
  title: 'Components/Form',
};

export const FormGroupWithWarning = () => (
  <FormGroupWarning inputLabel="Warning Message" warningMessage="Helpful warning message text go here" />
);
