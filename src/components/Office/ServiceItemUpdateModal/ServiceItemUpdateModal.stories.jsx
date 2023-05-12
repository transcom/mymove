import React from 'react';

import ConnectedServiceItemUpdateModal, { ServiceItemUpdateModal } from './ServiceItemUpdateModal';

export default {
  title: 'Office Components/ServiceItemUpdateModal',
};

const defaultProps = {
  closeModal: () => {},
  onSave: () => {},
  title: 'Edit Service Item',
};

const Template = () => <ServiceItemUpdateModal {...defaultProps} />;
export const ServiceItemUpdateModalStory = () => Template.bind();

const ConnectedTemplate = (args) => <ConnectedServiceItemUpdateModal {...args} />;
export const ConnectedModal = ConnectedTemplate.bind({});
ConnectedModal.args = {
  isOpen: true,
};
