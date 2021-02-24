import React from 'react';

import ConnectedFormGroupWarning from './FormGroupWarning';

export default {
  title: 'Components/Form',
};

export const FormGroupWithWarning = () => (
  <ConnectedFormGroupWarning
    inputLabel="TAC/MDC"
    warningMessage="This TAC does not appear in TGET, so might not be valid. Make sure it matches what's on the orders before you continue."
  />
);
