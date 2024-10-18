import React from 'react';

import DeleteMoveConfirmationModal from 'components/ConfirmationModals/CancelMoveConfirmationModal';

export default {
  title: 'Components/DeleteMoveConfirmationModal',
  component: DeleteMoveConfirmationModal,
  args: {
    customerSupportRemarkID: '111',
  },
  argTypes: {
    onClose: { action: 'close button clicked' },
    onSubmit: { action: 'submit button clicked' },
  },
};

const Template = (args) => <DeleteMoveConfirmationModal {...args} />;

export const Basic = Template.bind({});

export const WithOverrides = Template.bind({});
WithOverrides.args = {
  title: 'This is a sample title',
  content: 'Some sample description',
  submitText: 'YES!',
  closeText: 'NO',
};

const ConnectedTemplate = (args) => <DeleteMoveConfirmationModal {...args} />;
export const ConnectedModal = ConnectedTemplate.bind({});
ConnectedModal.args = {
  isOpen: true,
};
