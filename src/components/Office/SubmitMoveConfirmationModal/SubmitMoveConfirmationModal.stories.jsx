import React from 'react';

import { SubmitMoveConfirmationModal } from './SubmitMoveConfirmationModal';

export default {
  title: 'Office Components/SubmitMoveConfirmationModal',
  component: SubmitMoveConfirmationModal,
};

export const Basic = () => <SubmitMoveConfirmationModal onClose={() => {}} onSubmit={() => {}} />;
export const Shipment = () => <SubmitMoveConfirmationModal onClose={() => {}} onSubmit={() => {}} isShipment />;
