/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import MtoShipmentForm from './MtoShipmentForm';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { history, store } from 'shared/store';

const defaultProps = {
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  match: {
    isExact: false,
    path: 'moves/:moveId/shipments/:mtoShipmentId/edit?shipmentNumber=2',
    url: '',
    params: { moveId: 'move123' },
  },
  history: { push: () => {}, goBack: () => {} },
  showLoggedInUser: () => {},
  newDutyLocationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
    streetAddress1: '123 Main',
    streetAddress2: '',
  },
  useCurrentResidence: false,
  mtoShipment: {
    destinationAddress: undefined,
  },
  serviceMember: {
    weight_allotment: {
      total_weight_self: 5000,
    },
  },
  isCreatePage: true,
};

const mockMtoShipment = {
  id: 'mock id',
  moveTaskOrderId: 'mock move id',
  customerRemarks: 'mock remarks',
  requestedPickupDate: '2020-03-01',
  requestedDeliveryDate: '2020-03-30',
  agents: [
    {
      firstName: 'mock receiving',
      lastName: 'agent',
      telephone: '2225551234',
      email: 'mock.delivery.agent@example.com',
      agentType: 'RECEIVING_AGENT',
    },
    {
      firstName: 'Mock Releasing',
      lastName: 'Agent Jr, PhD, MD, DDS',
      telephone: '3335551234',
      email: 'mock.pickup.agent@example.com',
      agentType: 'RELEASING_AGENT',
    },
  ],
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

export default {
  title: 'Customer Components / Forms / MtoShipmentForm',
};

function renderStory(props) {
  return (
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <MtoShipmentForm {...defaultProps} {...props} />
      </ConnectedRouter>
    </Provider>
  );
}

// create shipment stories (form should not prefill customer data)
export const HHGShipment = () => renderStory({ selectedMoveType: SHIPMENT_OPTIONS.HHG });
export const NTSReleaseShipment = () => renderStory({ selectedMoveType: SHIPMENT_OPTIONS.NTSR });
export const NTSShipment = () => renderStory({ selectedMoveType: SHIPMENT_OPTIONS.NTS });

// edit shipment stories (form should prefill)
export const EditHHGShipment = () =>
  renderStory({
    selectedMoveType: SHIPMENT_OPTIONS.HHG,
    isCreatePage: false,
    mtoShipment: mockMtoShipment,
  });
export const EditNTSReleaseShipment = () =>
  renderStory({
    selectedMoveType: SHIPMENT_OPTIONS.NTSR,
    isCreatePage: false,
    mtoShipment: mockMtoShipment,
  });
export const EditNTSShipment = () =>
  renderStory({
    selectedMoveType: SHIPMENT_OPTIONS.NTS,
    isCreatePage: false,
    mtoShipment: mockMtoShipment,
  });

export const EditHHGShipmentWithSecondaryAddresses = () => {
  const extendedShipment = {
    ...mockMtoShipment,
    secondaryPickupAddress: {
      streetAddress1: '142 E Barrel Hoop Circle',
      streetAddress2: '#4A',
      city: 'Corpus Christi',
      state: 'TX',
      postalCode: '78412',
    },
    secondaryDeliveryAddress: {
      streetAddress1: '3373 NW Martin Luther King Jr Blvd',
      streetAddress2: '',
      city: mockMtoShipment.destinationAddress.city,
      state: mockMtoShipment.destinationAddress.state,
      postalCode: mockMtoShipment.destinationAddress.postalCode,
    },
  };

  return renderStory({
    selectedMoveType: SHIPMENT_OPTIONS.HHG,
    isCreatePage: false,
    mtoShipment: extendedShipment,
  });
};
