/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import { Home } from './index';

import { MockProviders } from 'testUtils';

export default {
  title: 'Customer Components / Pages / Home',
};

const defaultProps = {
  serviceMember: {
    id: 'testServiceMemberId',
    first_name: 'John',
    last_name: 'Lee',
    current_station: {
      name: 'Fort Knox',
      transportation_office: {
        name: 'Test Transportation Office Name',
        phone_lines: ['555-555-5555'],
      },
      weight_allotment: {},
    },
  },
  showLoggedInUser() {},
  loadMTOShipments() {},
  history: { push: () => {}, goBack: () => {} },
  getSignedCertification() {},
  mtoShipments: [],
  mtoShipment: {},
  isLoggedIn: true,
  loggedInUserIsLoading: false,
  loggedInUserSuccess: true,
  isProfileComplete: true,
  currentPpm: {},
  orders: {},
  location: {},
  move: {},
  uploadedOrderDocuments: [],
};

export const Basic = () => (
  <MockProviders>
    <div className="grid-container usa-prose">
      <Home {...defaultProps} />
    </div>
  </MockProviders>
);
