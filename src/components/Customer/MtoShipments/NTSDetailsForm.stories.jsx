/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import { NTSDetailsFormComponent as NTSDetailsForm } from './NTSDetailsForm';

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
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postal_code: '31905',
    street_address_1: '123 Main',
  },
};

export default {
  title: 'Customer Components | NTSDetailsForm',
};

export const Basic = () => <NTSDetailsForm {...defaultProps} />;
