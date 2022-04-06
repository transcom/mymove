import React from 'react';

import { DestructiveShipmentConfirmationModal } from 'components/ConfirmationModals/DestructiveShipmentConfirmationModal';

export default {
  title: 'Components/DestructiveShipmentConfirmationModal',
  component: DestructiveShipmentConfirmationModal,
};

export const Basic = () => <DestructiveShipmentConfirmationModal onClose={() => {}} onSubmit={() => {}} />;

export const withOverrides = () => (
  <DestructiveShipmentConfirmationModal
    onClose={() => {}}
    onSubmit={() => {}}
    title="This is a sample title"
    content="Some sample description"
    submitText="YES!"
    backText="NO"
  />
);
