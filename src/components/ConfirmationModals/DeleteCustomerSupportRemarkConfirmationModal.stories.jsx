import React from 'react';

import ConnectedDeleteCustomerSupportRemarkConfirmationModal, {
  DeleteCustomerSupportRemarkConfirmationModal,
} from 'components/ConfirmationModals/DeleteCustomerSupportRemarkConfirmationModal';

export default {
  title: 'Components/DeleteCustomerSupportRemarkConfirmationModal',
  component: DeleteCustomerSupportRemarkConfirmationModal,
  args: {
    customerSupportRemarkID: '111',
  },
  argTypes: {
    onClose: { action: 'close button clicked' },
    onSubmit: { action: 'submit button clicked' },
  },
};

const Template = (args) => <DeleteCustomerSupportRemarkConfirmationModal {...args} />;

export const Basic = Template.bind({});

export const WithOverrides = Template.bind({});
WithOverrides.args = {
  title: 'This is a sample title',
  content: 'Some sample description',
  submitText: 'YES!',
  closeText: 'NO',
};

const ConnectedTemplate = (args) => <ConnectedDeleteCustomerSupportRemarkConfirmationModal {...args} />;
export const ConnectedModal = ConnectedTemplate.bind({});
ConnectedModal.args = {
  isOpen: true,
};
