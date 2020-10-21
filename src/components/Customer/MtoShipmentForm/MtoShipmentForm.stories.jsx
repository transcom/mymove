/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import { MtoShipmentFormComponent as MtoShipmentForm } from './MtoShipmentForm';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { history, store } from 'shared/store';

const defaultProps = {
  pageList: ['page1', 'anotherPage/:foo/:bar'],
  pageKey: 'page1',
  match: { isExact: false, path: '', url: '', params: { moveId: 'move123' } },
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
  isCreatePage: true,
};

const mockMtoShipment = {
  id: 'mock id',
  moveTaskOrderId: 'mock move id',
  customerRemarks: 'mock remarks',
  requestedPickupDate: '1 Mar 2020',
  requestedDeliveryDate: '30 Mar 2020',
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
  title: 'Customer Components | MtoShipmentForm',
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
