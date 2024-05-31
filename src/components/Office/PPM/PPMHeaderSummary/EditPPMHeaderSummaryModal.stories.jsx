import React from 'react';
import { action } from '@storybook/addon-actions';

import EditPPMHeaderSummaryModal from './EditPPMHeaderSummaryModal';

export default {
  title: 'Office Components/EditPPMHeaderSummaryModal',
  component: EditPPMHeaderSummaryModal,
};

// Mock data for the story
const sectionInfo = {
  actualMoveDate: '2022-01-01',
  advanceAmountReceived: 50000,
};

export const Basic = (args) => {
  return <EditPPMHeaderSummaryModal {...args} />;
};

export const Default = Basic.bind({});
Default.args = {
  sectionType: 'shipmentInfo',
  sectionInfo,
  onClose: action('onClose'),
  onSubmit: action('onSubmit'),
  editSectionName: 'actualMoveDate',
};

export const EditShipmentInfo = Basic.bind({});
EditShipmentInfo.args = {
  sectionType: 'shipmentInfo',
  sectionInfo,
  onClose: action('onClose'),
  onSubmit: action('onSubmit'),
  editSectionName: 'actualMoveDate',
};

export const EditIncentives = Basic.bind({});
EditIncentives.args = {
  sectionType: 'incentives',
  sectionInfo,
  onClose: action('onClose'),
  onSubmit: action('onSubmit'),
  editSectionName: 'advanceAmountReceived',
};
