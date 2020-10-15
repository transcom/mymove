/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import { NTSrDetailsFormComponent as NTSrDetailsForm } from './NTSrDetailsForm';

const defaultProps = {
  wizardPage: {
    pageList: ['page1', 'anotherPage/:foo/:bar'],
    pageKey: 'page1',
    match: { isExact: false, path: '', url: '', params: { moveId: '123' } },
    history: { push: () => {}, goBack: () => {} },
  },
  showLoggedInUser: () => {},
  newDutyStationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
  },
  mtoShipment: {
    destinationAddress: undefined,
  },
};

export default {
  title: 'Customer Components | NTSrDetailsForm',
};

export const Basic = () => <NTSrDetailsForm {...defaultProps} />;
