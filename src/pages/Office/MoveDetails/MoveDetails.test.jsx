/*  react/jsx-props-no-spreading */
import React from 'react';
import { shallow } from 'enzyme';

import { MoveDetails } from './MoveDetails';

describe('MoveDetails page', () => {
  const testMoveOrder = {
    id: '123',
    customerID: 'abc',
  };

  const shipments = [
    {
      approvedDate: '2020-01-01',
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postal_code: '94535',
        state: 'CA',
        street_address_1: '987 Any Avenue',
        street_address_2: 'P.O. Box 9876',
        street_address_3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postal_code: '90210',
        state: 'CA',
        street_address_1: '123 Any Street',
        street_address_2: 'P.O. Box 12345',
        street_address_3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
  ];

  const testMoveTaskOrders = [{ id: '1a' }, { id: '2b' }];

  const testMTOShipments = [
    { id: '1a', status: 'SUBMITTED', currentAddress: { street_address_1: 'test' } },
    { id: '2b', status: 'SUBMITTED', currentAddress: { street_address_1: 'test' } },
  ];

  const testMtoServiceItems = [{ id: '1a', status: 'APPROVED' }];

  const moveDetailsProps = {
    match: {
      params: { moveOrderId: '123' },
      isExact: true,
      path: '',
      url: '',
    },
    getMoveOrder: jest.fn(() => new Promise((res) => res({ response: { body: testMoveOrder } }))),
    getCustomer: jest.fn(),
    getAllMoveTaskOrders: jest.fn(() => new Promise((res) => res({ response: { body: testMoveTaskOrders } }))),
    getMTOShipments: jest.fn(() => new Promise((res) => res({ response: { body: testMTOShipments } }))),
    patchMTOShipmentStatus: jest.fn(() => new Promise((res) => res({ response: { body: testMTOShipments[0] } }))),
    getMTOServiceItems: jest.fn(() => new Promise((res) => res({ response: { body: testMtoServiceItems } }))),
    mtoShipments: shipments,
    updateMoveTaskOrderStatus: jest
      .fn()
      .mockResolvedValue({ response: { status: 200, body: { id: '1a', eTag: '1a2b3c4d' } } }),
  };

  const wrapper = shallow(<MoveDetails {...moveDetailsProps} />);

  it('loads data from the API', () => {
    expect(moveDetailsProps.getMoveOrder).toHaveBeenCalledWith(moveDetailsProps.match.params.moveOrderId);
    expect(moveDetailsProps.getCustomer).toHaveBeenCalledWith(testMoveOrder.customerID);
    expect(moveDetailsProps.getAllMoveTaskOrders).toHaveBeenCalledWith(testMoveOrder.id);
    expect(moveDetailsProps.getMTOShipments).toHaveBeenCalledWith(testMoveTaskOrders[0].id);
    expect(moveDetailsProps.getMTOShipments).toHaveBeenCalledWith(testMoveTaskOrders[1].id);
  });

  it('renders the h1', () => {
    expect(wrapper.find({ 'data-testid': 'too-move-details' }).exists()).toBe(true);
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
