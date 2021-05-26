/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';
import { action } from '@storybook/addon-actions';

import ServicesCounselingShipmentForm from './ServicesCounselingShipmentForm';

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
  newDutyStationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
    street_address_1: '123 Main',
  },
  useCurrentResidence: false,
  mtoShipment: {
    destinationAddress: undefined,
  },
  serviceMember: {
    weightAllotment: {
      totalWeightSelf: 5000,
    },
  },
  isCreatePage: true,
  updateMTOShipment: action('update MTO shipment'),
};

const mockMtoShipment = {
  id: 'mock id',
  moveTaskOrderId: 'mock move id',
  customerRemarks: 'mock customer remarks',
  counselorRemarks: 'mock counselor remarks',
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
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    street_address_1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postal_code: '98421',
  },
};

export default {
  title: 'Office Components / Forms / ServicesCounselingShipmentForm',
};

function renderStory(props) {
  return (
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <ServicesCounselingShipmentForm {...defaultProps} {...props} />
      </ConnectedRouter>
    </Provider>
  );
}

// create shipment stories (form should not prefill customer data)
export const HHGShipment = () => renderStory({ selectedMoveType: SHIPMENT_OPTIONS.HHG });

// edit shipment stories (form should prefill)
export const EditHHGShipment = () =>
  renderStory({
    selectedMoveType: SHIPMENT_OPTIONS.HHG,
    isCreatePage: false,
    mtoShipment: mockMtoShipment,
  });
