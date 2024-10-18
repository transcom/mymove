import React from 'react';

import CancelMoveConfirmationModal from './CancelMoveConfirmationModal';

export default {
  title: 'Components/CancelMoveConfirmationModal',
  component: CancelMoveConfirmationModal,
  args: {
    moveID: '111',
  },
  argTypes: {
    onClose: { action: 'close button clicked' },
    onSubmit: { action: 'submit button clicked' },
  },
};

const Template = (args) => <CancelMoveConfirmationModal {...args} />;

export const Basic = Template.bind({});

export const WithOverrides = Template.bind({});
WithOverrides.args = {
  title: 'This is a sample title',
  content: 'Some sample description',
  submitText: 'YES!',
  closeText: 'NO',
};

const ConnectedTemplate = (args) => <CancelMoveConfirmationModal {...args} />;
export const ConnectedModal = ConnectedTemplate.bind({});
ConnectedModal.args = {
  isOpen: true,
};
