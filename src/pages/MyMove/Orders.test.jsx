/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import ConnectedOrders from './Orders';

import { MockProviders } from 'testUtils';

describe('Orders page', () => {
  describe('with no existing orders', () => {
    const initialState = {
      serviceMember: {
        currentServiceMember: {
          id: 'testServiceMember123',
        },
      },
    };

    const mockHistory = {
      push: jest.fn(),
      goBack: jest.fn(),
    };

    const testProps = {
      fetchLatestOrders: jest.fn(),
      updateOrders: jest.fn(),
      createOrders: jest.fn(),
      history: mockHistory,
      pages: [],
      pageKey: '',
      match: { params: {} },
    };

    const wrapper = mount(
      <MockProviders initialState={initialState} initialEntries={['/']}>
        <ConnectedOrders {...testProps} />
      </MockProviders>,
    );

    it('renders without errors', () => {
      expect(wrapper.exists()).toBe(true);
    });

    it('does not fetch latest orders on mount', () => {
      expect(testProps.fetchLatestOrders).not.toHaveBeenCalled();
    });
  });

  // TODO - add tests after updating API data flow for this component
  /*
  describe('with existing orders', () => {
    const initialState = {
      user: {
        userInfo: {
          service_member: {
            id: 'testServiceMember123',
          },
        },
      },
      serviceMember: {
        currentServiceMember: {
          id: 'testServiceMember123',
        },
      },
      entities: {
        orders: {
          orders123: {
            service_member_id: 'testServiceMember123',
          },
        },
      },
    };

    const mockHistory = {
      push: jest.fn(),
      goBack: jest.fn(),
    };

    const testProps = {
      fetchLatestOrders: jest.fn(),
      updateOrders: jest.fn(),
      createOrders: jest.fn(),
      history: mockHistory,
      pages: [],
      pageKey: '',
      match: { params: {} },
    };

    const wrapper = mount(
      <MockProviders initialState={initialState} initialEntries={['/']}>
        <ConnectedOrders {...testProps} />
      </MockProviders>,
    );

    it('renders without errors', () => {
      expect(wrapper.exists()).toBe(true);
      expect(testProps.fetchLatestOrders).toHaveBeenCalled();
    });
  });
  */
});
