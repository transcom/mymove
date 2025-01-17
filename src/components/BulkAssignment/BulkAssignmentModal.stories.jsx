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
  bulkAssignmentData: {
    bulkAssignmentMoveIDs: ['1', '2', '3', '4', '5', '6', '7', '8'],
    availableOfficeUsers: [
      { lastName: 'Monk', firstName: 'Art', workload: 81 },
      { lastName: 'Green', firstName: 'Darrell', workload: 28 },
      { lastName: 'Riggins', firstName: 'John', workload: 44 },
    ],
  },
};
