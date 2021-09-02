import React from 'react';

import EditMaxBillableWeightModal from './EditMaxBillableWeightModal';

export default {
  title: 'Office Components/EditMaxBillableWeightModal',
  component: EditMaxBillableWeightModal,
  argTypes: {
    defaultWeight: { type: 'number', defaultValue: 10000 },
    maxBillableWeight: { type: 'number', defaultValue: 10999 },
    onSubmit: { action: 'submit form' },
    onClose: { action: 'close modal' },
  },
};

export const Basic = (args) => <EditMaxBillableWeightModal {...args} />;
