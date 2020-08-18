import React from 'react';
import { mount } from 'enzyme';

import { MoveTaskOrder } from 'pages/Office/MoveTaskOrder/MoveTaskOrder';

jest.mock('hooks/queries', () => ({
  useMoveTaskOrderQueries: () => {
    return {
      moveOrders: {
        1: {
          id: '1',
        },
      },
      moveTaskOrders: {
        2: {
          id: '2',
        },
      },
      mtoShipments: {
        3: {
          id: '3',
          shipmentType: 'HHG',
          scheduledPickupDate: '2020-03-16',
          requestedPickupDate: '2020-03-15',
        },
      },
      mtoServiceItems: {
        4: {
          id: '4',
          mtoShipmentID: 'a1',
          reServiceName: 'Test Service Item',
        },
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    };
  },
}));

describe('MoveTaskOrder', () => {
  const testId = 'test-id-123';
  const requiredProps = {
    match: { params: { moveTaskOrderId: testId } },
    history: { push: jest.fn() },
  };

  // eslint-disable-next-line react/jsx-props-no-spreading
  const wrapper = mount(<MoveTaskOrder {...requiredProps} />);

  it('should render the h1', () => {
    expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
  });
});
