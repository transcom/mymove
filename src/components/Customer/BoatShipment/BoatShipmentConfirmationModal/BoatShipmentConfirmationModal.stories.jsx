import React from 'react';

import ConnectedBoatShipmentConfirmationModal, {
  BoatShipmentConfirmationModal,
} from 'components/Customer/BoatShipment/BoatShipmentConfirmationModal/BoatShipmentConfirmationModal';

export default {
  title: 'Components/BoatShipmentConfirmationModal',
  component: BoatShipmentConfirmationModal,
  args: {
    isDimensionsMeetReq: true,
    boatShipmentType: 'TOW_AWAY',
    isEditPage: false,
    isSubmitting: false,
  },
  argTypes: {
    closeModal: { action: 'close button clicked' },
    handleConfirmationContinue: { action: 'continue button clicked' },
    handleConfirmationRedirect: { action: 'redirect button clicked' },
    handleConfirmationDeleteAndRedirect: { action: 'delete and continue button clicked' },
  },
};

const Template = (args) => <BoatShipmentConfirmationModal {...args} />;

export const Default = Template.bind({});

export const EditPageWithNonMeetingDimensions = Template.bind({});
EditPageWithNonMeetingDimensions.args = {
  isDimensionsMeetReq: false,
  isEditPage: true,
};

export const TowAway = Template.bind({});
TowAway.args = {
  isDimensionsMeetReq: true,
  boatShipmentType: 'TOW_AWAY',
};

export const HaulAway = Template.bind({});
HaulAway.args = {
  isDimensionsMeetReq: true,
  boatShipmentType: 'HAUL_AWAY',
};

export const NonMeetingDimensions = Template.bind({});
NonMeetingDimensions.args = {
  isDimensionsMeetReq: false,
  boatShipmentType: '',
};

const ConnectedTemplate = (args) => <ConnectedBoatShipmentConfirmationModal {...args} />;
export const ConnectedModal = ConnectedTemplate.bind({});
ConnectedModal.args = {
  isOpen: true,
};
