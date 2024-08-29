import React from 'react';

import ConnectedDeleteDocumentFileConfirmationModal, {
  DeleteDocumentFileConfirmationModal,
} from 'components/ConfirmationModals/DeleteDocumentFileConfirmationModal';

export default {
  title: 'Components/DeleteDocumentFileConfirmationModal',
  component: DeleteDocumentFileConfirmationModal,
  argTypes: {
    closeModal: { action: 'close button clicked' },
    submitModal: { action: 'submit button clicked' },
  },
};

const Template = (args) => <DeleteDocumentFileConfirmationModal {...args} />;

export const Basic = Template.bind({});
Basic.args = {
  fileInfo: {
    filename: 'test-file',
    bytes: '1212',
    createdAt: '12/01/2024',
  },
};

const ConnectedTemplate = (args) => <ConnectedDeleteDocumentFileConfirmationModal {...args} />;
export const ConnectedModal = ConnectedTemplate.bind({});
ConnectedModal.args = {
  fileInfo: {
    filename: 'test-file',
    bytes: '1212',
    createdAt: '12/01/2024',
  },
  isOpen: true,
};
