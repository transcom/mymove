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

export const HHGShipment = () => renderStory({ selectedMoveType: SHIPMENT_OPTIONS.HHG });
export const NTSReleaseShipment = () => renderStory({ selectedMoveType: SHIPMENT_OPTIONS.NTSR });
export const NTSShipment = () => renderStory({ selectedMoveType: SHIPMENT_OPTIONS.NTS });
