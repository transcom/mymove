import React from 'react';

import ConnectedServiceItemUpdateModal, { ServiceItemUpdateModal } from './ServiceItemUpdateModal';
import EditSitAddressChangeForm from './EditSitAddressChangeForm';

const destinationSIT = {
  id: 'abc123',
  code: 'DDDSIT',
  submittedAt: '2020-11-20',
  serviceItem: 'Domestic destination SIT',
  details: {
    reason: "Customer's housing at base is not ready",
    firstCustomerContact: { timeMilitary: '1200Z', firstAvailableDeliveryDate: '2020-09-15' },
    secondCustomerContact: { timeMilitary: '2300Z', firstAvailableDeliveryDate: '2020-09-21' },
    serviceItem: 'Domestic Destination SIT',
  },
};

const initialAddress = {
  city: 'Fairfax',
  state: 'VA',
  postalCode: '12345',
  streetAddress1: '123 Fake Street',
  streetAddress2: '',
  streetAddress3: '',
  country: 'USA',
};

export default {
  title: 'Office Components/ServiceItemUpdateModal',
  component: ServiceItemUpdateModal,
  args: {
    closeModal: () => {},
    onSave: () => {},
    isOpen: true,
    serviceItem: destinationSIT,
  },
  argTypes: {
    onClose: { action: 'close button clicked' },
    onSubmit: { action: 'submit button clicked' },
  },
};

const Template = (args) => <ServiceItemUpdateModal {...args} />;
export const ServiceItemUpdateModalStory = Template.bind({});
ServiceItemUpdateModalStory.args = {
  title: 'Title',
};

// Creates template of the modal
const ConnectedTemplate = (args) => <ConnectedServiceItemUpdateModal {...args} />;
// Story for Editing service item address
export const EditServiceItemAddress = ConnectedTemplate.bind({});
EditServiceItemAddress.args = {
  title: 'Edit Service Item',
  content: <EditSitAddressChangeForm initialAddress={initialAddress} />,
};

// To-do: Setup story for Reviewing Service Item requests
