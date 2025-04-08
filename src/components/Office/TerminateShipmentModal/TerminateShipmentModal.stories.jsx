import React from 'react';
import { action } from '@storybook/addon-actions';

import { TerminateShipmentModal } from './TerminateShipmentModal';

export default {
  title: 'Office Components/Terminate Shipment Modal',
  component: TerminateShipmentModal,
};

const props = {
  onClose: action('clicked'),
  onSubmit: action('clicked'),
  shipmentID: 'shipmentID',
  shipmentLocator: 'SHIPMT-01',
};

export const Modal = () => <TerminateShipmentModal {...props} />;
