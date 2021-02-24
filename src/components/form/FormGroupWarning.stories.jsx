import React from 'react';

import ConnectedFormGroupWarning from './FormGroupWarning';

export default {
  title: 'Components/Form',
};

export const FormGroupWithWarning = () => (
  <ConnectedFormGroupWarning inputLabel="Warning Message" warningMessage="Helpful warning message text go here" />
);
