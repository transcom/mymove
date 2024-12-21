import React from 'react';

import BulkAssignmentModal from './BulkAssignmentModal';

export default {
  title: 'Components/BulkAssignmentModal',
  component: BulkAssignmentModal,
  argTypes: {
    onClose: { action: 'close button clicked' },
    onSubmit: { action: 'submit button clicked' },
  },
};

const Template = (args) => <BulkAssignmentModal {...args} />;

export const Basic = Template.bind({});

export const WithOverrides = Template.bind({});
WithOverrides.args = {
  title: 'This is a sample title',
  content: 'Some sample description',
  submitText: 'YES!',
  closeText: 'NO',
};

const ConnectedTemplate = (args) => <BulkAssignmentModal {...args} />;
export const ConnectedModal = ConnectedTemplate.bind({});
ConnectedModal.args = {
  isOpen: true,
};
