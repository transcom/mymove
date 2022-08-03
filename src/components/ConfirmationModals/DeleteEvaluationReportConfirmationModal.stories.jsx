import React from 'react';

import ConnectedDeleteEvaluationReportConfirmationModal, {
  DeleteEvaluationReportConfirmationModal,
} from 'components/ConfirmationModals/DeleteEvaluationReportConfirmationModal';

export default {
  title: 'Components/DeleteEvaluationReportConfirmationModal',
  component: DeleteEvaluationReportConfirmationModal,
  argTypes: {
    closeModal: { action: 'close button clicked' },
    submitModal: { action: 'submit button clicked' },
  },
};

const Template = (args) => <DeleteEvaluationReportConfirmationModal {...args} />;

export const Basic = Template.bind({});

const ConnectedTemplate = (args) => <ConnectedDeleteEvaluationReportConfirmationModal {...args} />;
export const ConnectedModal = ConnectedTemplate.bind({});
ConnectedModal.args = {
  isOpen: true,
};
