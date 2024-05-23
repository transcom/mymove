import React from 'react';

import { IncompleteShipmentModal } from 'components/Customer/Review/IncompleteShipmentModal/IncompleteShipmentModal';

const noop = () => {};

export default {
  title: 'Components/IncompleteShipmentModal',
  component: IncompleteShipmentModal,
  args: {
    shipmentLabel: 'PPM 1',
    shipmentMoveCode: '20FDBF58',
    shipmentType: 'PPM',
    closeModal: noop,
  },
};

const Template = (args) => <IncompleteShipmentModal {...args} />;

export const Basic = Template.bind({});
