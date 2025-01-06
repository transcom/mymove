import React from 'react';

import { MovingInfo } from './MovingInfo';

import { MockProviders } from 'testUtils';

export default {
  title: 'Customer Components / Pages / Move 101',
};

const props = {
  fetchLatestOrders: () => {},
  serviceMemberId: 1231231231,
  location: {},
  router: {
    params: {
      moveId: 'A1B2C3',
    },
    navigate: () => {},
  },
};

export const WithEntitlementWeight = () => (
  <MockProviders>
    <MovingInfo {...props} entitlementWeight={1234} />
  </MockProviders>
);

export const WithoutEntitlementWeight = () => (
  <MockProviders>
    <MovingInfo {...props} />
  </MockProviders>
);

export const WithUBAllowance = () => (
  <MockProviders>
    <MovingInfo {...props} entitlementWeight={1234} ubAllowance={2000} />
  </MockProviders>
);
