import React from 'react';
import { mount } from 'enzyme';

import { MoveTaskOrder } from 'pages/Office/MoveTaskOrder/MoveTaskOrder';

jest.mock('hooks/queries', () => ({
  useMoveTaskOrderQueries: () => {
    return {
      moveOrders: {
        1: {
          id: '1',
          originDutyStation: {
            address: {
              street_address_1: '',
              city: 'Fort Knox',
              state: 'KY',
              postal_code: '40121',
            },
          },
          destinationDutyStation: {
            address: {
              street_address_1: '',
              city: 'Fort Irwin',
              state: 'CA',
              postal_code: '92310',
            },
          },
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
          pickupAddress: {
            street_address_1: '932 Baltic Avenue',
            city: 'Chicago',
            state: 'IL',
            postal_code: '60601',
          },
          destinationAddress: {
            street_address_1: '10 Park Place',
            city: 'Atlantic City',
            state: 'NJ',
            postal_code: '08401',
          },
          status: 'APPROVED',
        },
        4: {
          id: '4',
          shipmentType: 'HHG_INTO_NTS_DOMESTIC',
          scheduledPickupDate: '2020-03-16',
          requestedPickupDate: '2020-03-15',
          pickupAddress: {
            street_address_1: '932 Baltic Avenue',
            city: 'Chicago',
            state: 'IL',
            postal_code: '60601',
          },
          destinationAddress: {
            street_address_1: '10 Park Place',
            city: 'Atlantic City',
            state: 'NJ',
            postal_code: '08401',
          },
          status: 'APPROVED',
        },
      },
      mtoServiceItems: {
        5: {
          id: '5',
          mtoShipmentID: '3',
          reServiceName: 'Test Service Item',
          status: 'SUBMITTED',
          reServiceCode: 'DOFSIT',
        },
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    };
  },
}));

const setUnapprovedShipmentCount = jest.fn();

describe('MoveTaskOrder', () => {
  const moveCode = 'WE31AZ';
  const requiredProps = {
    match: { params: { moveCode } },
    history: { push: jest.fn() },
  };

  // eslint-disable-next-line react/jsx-props-no-spreading
  const wrapper = mount(<MoveTaskOrder {...requiredProps} setUnapprovedShipmentCount={setUnapprovedShipmentCount} />);

  it('should render the h1', () => {
    expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
    expect(wrapper.find('h1').text()).toBe('Move task order');
  });

  it('should render the ShipmentContainer', () => {
    expect(wrapper.find('ShipmentContainer').exists()).toBe(true);
  });

  it('should render the ShipmentHeading', () => {
    expect(wrapper.find('ShipmentHeading').exists()).toBe(true);
    expect(wrapper.find('h3').at(0).text()).toEqual('Household goods');
    expect(wrapper.find('h3').at(1).text()).toEqual('Non-temp storage');
  });

  it('should render the ImportantShipmentDates', () => {
    expect(wrapper.find('ImportantShipmentDates').exists()).toBe(true);
  });

  it('should render the ShipmentAddresses', () => {
    expect(wrapper.find('ShipmentAddresses').exists()).toBe(true);
  });

  it('should render the RequestedServiceItemsTable', () => {
    expect(wrapper.find('RequestedServiceItemsTable').exists()).toBe(true);
  });

  it('updates the unapproved shipments tag state', () => {
    expect(setUnapprovedShipmentCount).toHaveBeenCalledWith(0);
  });
});
