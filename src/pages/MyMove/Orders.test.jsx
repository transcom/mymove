/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import ConnectedOrders, { Orders } from './Orders';

import { MockProviders } from 'testUtils';

describe('Orders page', () => {
  const mockHistory = {
    push: jest.fn(),
    goBack: jest.fn(),
  };

  const ordersOptions = [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ];

  describe('with the allOrdersType feature flag set to true', () => {
    const wrapper = mount(
      <Orders
        serviceMemberId="123"
        match={{ params: { moveId: 'test' }, path: '/orders', url: '/orders', isExact: false }}
        history={mockHistory}
        context={{ flags: { allOrdersTypes: true } }}
      />,
    );
    it('passes all orders types into the form', () => {
      expect(wrapper.find('OrdersInfoForm').prop('ordersTypeOptions')).toEqual(ordersOptions);
    });
  });

  describe('with the allOrdersType feature flag set to false', () => {
    const wrapper = mount(
      <Orders
        serviceMemberId="123"
        match={{ params: { moveId: 'test' }, path: '/orders', url: '/orders', isExact: false }}
        history={mockHistory}
        context={{ flags: { allOrdersTypes: false } }}
      />,
    );
    it('passes only the PCS option into the form', () => {
      expect(wrapper.find('OrdersInfoForm').prop('ordersTypeOptions')).toEqual([ordersOptions[0]]);
    });
  });

  describe('with no existing orders', () => {
    const initialState = {
      user: {
        userInfo: {
          service_member: {
            id: 'testServiceMember123',
          },
        },
      },
      entities: {
        orders: {},
        serviceMembers: {
          testServiceMember123: {
            id: 'testServiceMember123',
          },
        },
        users: {
          testUserId: {
            service_member: 'testServiceMember123',
          },
        },
      },
    };

    const testProps = {
      updateOrders: jest.fn(),
      updateServiceMember: jest.fn(),
      history: mockHistory,
      pages: [],
      pageKey: '',
      match: { params: { moveId: 'test' }, path: '/orders', url: '/orders', isExact: false },
    };

    const wrapper = mount(
      <MockProviders initialState={initialState} initialEntries={['/orders']}>
        <ConnectedOrders {...testProps} />
      </MockProviders>,
    );

    it('renders without errors', () => {
      expect(wrapper.exists()).toBe(true);
    });

    it('does not fetch latest orders on mount', () => {
      expect(testProps.updateOrders).not.toHaveBeenCalled();
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
      entities: {
        orders: {
          orders123: {
            service_member_id: 'testServiceMember123',
          },
        },
        serviceMembers: {
          testServiceMember123: {
            id: 'testServiceMember123',
          },
        },
        users: {
          testUserId: {
            service_member: 'testServiceMember123',
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
