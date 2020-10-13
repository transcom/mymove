/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'connected-react-router';

import { HHGDetailsFormComponent as HHGDetailsForm } from './HHGDetailsForm';

import { history, store } from 'shared/store';

const defaultProps = {
  wizardPage: {
    pageList: ['page1', 'anotherPage/:foo/:bar'],
    pageKey: 'page1',
    match: { isExact: false, path: '', url: '', params: { moveId: 'move123' } },
    history: { push: () => {}, goBack: () => {} },
  },
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
  title: 'Customer Components | HHGDetailsForm',
};

function renderStory(props) {
  return (
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <HHGDetailsForm {...defaultProps} {...props} />
      </ConnectedRouter>
    </Provider>
  );
}

export const DefaultInitialState = () => renderStory();
