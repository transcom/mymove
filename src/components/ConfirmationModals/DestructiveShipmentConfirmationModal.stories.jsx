import React from 'react';

import ConnectedDestructiveShipmentConfirmationModal, {
  DestructiveShipmentConfirmationModal,
} from 'components/ConfirmationModals/DestructiveShipmentConfirmationModal';

export default {
  title: 'Components/DestructiveShipmentConfirmationModal',
  component: DestructiveShipmentConfirmationModal,
  args: {
    shipmentID: '111',
  },
  argTypes: {
    onClose: { action: 'close button clicked' },
    onSubmit: { action: 'submit button clicked' },
  },
};

const Template = (args) => <DestructiveShipmentConfirmationModal {...args} />;

export const Basic = Template.bind({});

export const WithOverrides = Template.bind({});
WithOverrides.args = {
  title: 'This is a sample title',
  content: 'Some sample description',
  submitText: 'YES!',
  closeText: 'NO',
  onClose: { action: 'close button clicked' },
};

const ConnectedTemplate = (args) => <ConnectedDestructiveShipmentConfirmationModal {...args} />;
export const ConnectedModal = ConnectedTemplate.bind({});
ConnectedModal.args = {
  isOpen: true,
};
