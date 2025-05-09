import React from 'react';
import { action } from '@storybook/addon-actions';
import { Provider } from 'react-redux';

import EditPPMHeaderSummaryModal from './EditPPMHeaderSummaryModal';

import { configureStore } from 'shared/store';

export default {
  title: 'Office Components/EditPPMHeaderSummaryModal',
  component: EditPPMHeaderSummaryModal,
};

const mockStore = configureStore({});

// Mock data for the story
const sectionInfo = {
  actualMoveDate: '2022-01-01',
  advanceAmountReceived: 50000,
  destinationAddressObj: {
    city: 'Fairfield',
    country: 'US',
    id: '672ff379-f6e3-48b4-a87d-796713f8f997',
    postalCode: '94535',
    county: 'Solano',
    state: 'CA',
    streetAddress1: '987 Any Avenue',
    streetAddress2: 'P.O. Box 9876',
    streetAddress3: 'c/o Some Person',
  },
  pickupAddressObj: {
    city: 'Beverly Hills',
    country: 'US',
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
    id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
    postalCode: '90210',
    county: 'Los Angeles',
    state: 'CA',
    streetAddress1: '123 Any Street',
    streetAddress2: 'P.O. Box 12345',
    streetAddress3: 'c/o Some Person',
  },
};

export const Basic = (args) => {
  return (
    <Provider store={mockStore.store}>
      <EditPPMHeaderSummaryModal sectionInfo={sectionInfo} {...args} />
    </Provider>
  );
};

export const EditShipmentInfo = Basic.bind({});
EditShipmentInfo.args = {
  sectionType: 'shipmentInfo',
  sectionInfo,
  onClose: action('onClose'),
  onSubmit: action('onSubmit'),
  editItemName: 'actualMoveDate',
};

export const EditPickupAddress = Basic.bind({});
EditPickupAddress.args = {
  sectionType: 'shipmentInfo',
  sectionInfo,
  onClose: action('onClose'),
  onSubmit: action('onSubmit'),
  editItemName: 'pickupAddress',
};

export const EditDestinationAddress = Basic.bind({});
EditDestinationAddress.args = {
  sectionType: 'shipmentInfo',
  sectionInfo,
  onClose: action('onClose'),
  onSubmit: action('onSubmit'),
  editItemName: 'destinationAddress',
};

export const EditIncentives = Basic.bind({});
EditIncentives.args = {
  sectionType: 'incentives',
  sectionInfo,
  onClose: action('onClose'),
  onSubmit: action('onSubmit'),
  editItemName: 'advanceAmountReceived',
};

export const EditExpenseType = Basic.bind({});
EditExpenseType.args = {
  sectionType: 'shipmentInfo',
  sectionInfo,
  onClose: action('onClose'),
  onSubmit: action('onSubmit'),
  editItemName: 'expenseType',
};
