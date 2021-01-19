/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import ConnectedUploadOrders from './UploadOrders';

import { MockProviders } from 'testUtils';

const defaultProps = {
  pages: ['1', '2', '3'],
  pageKey: '2',
  fetchLatestOrders: () => {},
};

const initialState = {
  entities: {
    orders: {},
    user: {
      testUser: {
        id: 'testUser',
        service_member: 'serviceMemberId',
      },
    },
    serviceMembers: {
      serviceMemberId: {
        id: 'serviceMemberId',
      },
    },
  },
};

const mountUploadOrders = (props = {}) =>
  mount(
    <MockProviders initialState={initialState}>
      <ConnectedUploadOrders {...defaultProps} {...props} />
    </MockProviders>,
  );

describe('UploadOrders component', () => {
  it('renders without errors', () => {
    const wrapper = mountUploadOrders();
    expect(wrapper.exists()).toBe(true);
  });
});
