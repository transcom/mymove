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

const ConnectedTemplate = (args) => <BulkAssignmentModal {...args} />;
export const ConnectedModal = ConnectedTemplate.bind({});
ConnectedModal.args = {
  isOpen: true,
};
