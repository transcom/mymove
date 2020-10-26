/*  react/jsx-props-no-spreading */
import React from 'react';
import { shallow } from 'enzyme';

import { Summary } from './index';

import { MOVE_STATUSES } from 'shared/constants';

const defaultProps = {
  serviceMember: {
    id: '666',
    current_station: {},
    residential_address: {},
  },
  currentOrders: {
    orders_type: 'PERMANENT_CHANGE_OF_STATION',
    has_dependents: false,
    issue_date: '2020-08-11',
    moves: ['123'],
    new_duty_station: {
      address: {
        postal_code: '123456',
      },
    },
    report_by_date: '2020-08-31',
    service_member_id: '666',
    spouse_has_pro_gear: false,
    status: MOVE_STATUSES.DRAFT,
    uploaded_orders: {
      uploads: [],
    },
  },
  match: { path: '', url: '/moves/123/review', params: { moveId: '123' } },
  currentMove: {
    id: '123',
    locator: 'CXVV3F',
    selected_move_type: 'HHG',
    service_member_id: '666',
    status: MOVE_STATUSES.DRAFT,
  },
  selectedMoveType: 'HHG',
  currentPPM: {},
  mtoShipment: {
    agents: [],
    customerRemarks: 'please be carefule',
    moveTaskOrderID: '123',
    pickupAddress: {
      city: 'Beverly Hills',
    },
    requestedDeliveryDate: '2020-08-31',
    requestedPickupDate: '2020-08-31',
    shipmentType: 'HHG',
    status: MOVE_STATUSES.SUBMITTED,
    updatedAt: '2020-09-02T21:08:38.392Z',
  },
  mtoShipments: [
    {
      agents: [],
      customerRemarks: 'please be carefule',
      moveTaskOrderID: '123',
      pickupAddress: {
        city: 'Beverly Hills',
      },
      requestedDeliveryDate: '2020-08-31',
      requestedPickupDate: '2020-08-31',
      shipmentType: 'HHG',
      status: MOVE_STATUSES.SUBMITTED,
      updatedAt: '2020-09-02T21:08:38.392Z',
    },
  ],
  reviewState: {
    editSuccess: null,
    entitlementChange: null,
    error: null,
  },
  onDidMount: jest.fn(),
  onCheckEntitlement: jest.fn(),
  showLoggedInUser: jest.fn(),
};

describe('Summary page', () => {
  const wrapper = shallow(<Summary {...defaultProps} />);
  it('renders heading with existing mtoShipment', () => {
    expect(wrapper.containsMatchingElement(<h3>Add another shipment</h3>)).toBe(true);
  });
  it('add shipment button exists', () => {
    const btn = wrapper.find('.usa-button--secondary');
    expect(btn.props().children).toBe('Add another shipment');
  });
});
