/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { shallow } from 'enzyme';

import { MoveDetails } from './MoveDetails';

describe('MoveDetails page', () => {
  const testMoveOrder = {
    id: '123',
    customerID: 'abc',
  };

  const testMoveTaskOrders = [{ id: '1a' }, { id: '2b' }];

  const moveDetailsProps = {
    match: {
      params: { locator: '123' },
      isExact: true,
      path: '',
      url: '',
    },
    getMoveOrder: jest.fn(() => new Promise((res) => res({ response: { body: testMoveOrder } }))),
    getCustomer: jest.fn(),
    getAllMoveTaskOrders: jest.fn(() => new Promise((res) => res({ response: { body: testMoveTaskOrders } }))),
    getMTOShipments: jest.fn(),
  };

  const wrapper = shallow(<MoveDetails {...moveDetailsProps} />);

  it('loads data from the API', () => {
    expect(moveDetailsProps.getMoveOrder).toHaveBeenCalledWith(moveDetailsProps.match.params.locator);
    expect(moveDetailsProps.getCustomer).toHaveBeenCalledWith(testMoveOrder.customerID);
    expect(moveDetailsProps.getAllMoveTaskOrders).toHaveBeenCalledWith(testMoveOrder.id);
    expect(moveDetailsProps.getMTOShipments).toHaveBeenCalledWith(testMoveTaskOrders[0].id);
    expect(moveDetailsProps.getMTOShipments).toHaveBeenCalledWith(testMoveTaskOrders[1].id);
  });

  it('renders the h1', () => {
    expect(wrapper.find({ 'data-cy': 'too-move-details' }).exists()).toBe(true);
    expect(wrapper.containsMatchingElement(<h1>Move details</h1>)).toBe(true);
  });

  it('renders side navigation for each section', () => {
    expect(wrapper.containsMatchingElement(<a href="#requested-shipments">Requested shipments</a>)).toBe(true);
    expect(wrapper.containsMatchingElement(<a href="#orders">Orders</a>)).toBe(true);
    expect(wrapper.containsMatchingElement(<a href="#allowances">Allowances</a>)).toBe(true);
    expect(wrapper.containsMatchingElement(<a href="#customer-info">Customer info</a>)).toBe(true);
  });

  it('renders the Requested Shipments component', () => {
    expect(wrapper.find('RequestedShipments')).toHaveLength(1);
  });

  it('renders the Orders Table', () => {
    expect(wrapper.find('OrdersTable')).toHaveLength(1);
  });

  it('renders the Allowances Table', () => {
    expect(wrapper.find('AllowancesTable')).toHaveLength(1);
  });

  it('renders the Customer Info Table', () => {
    expect(wrapper.find('CustomerInfoTable')).toHaveLength(1);
  });
});
