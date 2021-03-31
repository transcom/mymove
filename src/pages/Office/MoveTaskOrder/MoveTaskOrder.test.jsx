/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import { MoveTaskOrder } from 'pages/Office/MoveTaskOrder/MoveTaskOrder';
import MOVE_STATUSES from 'constants/moves';
import { shipmentStatuses } from 'constants/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import SERVICE_ITEM_STATUS from 'constants/serviceItems';
import { useMoveTaskOrderQueries } from 'hooks/queries';
import { MockProviders } from 'testUtils';

jest.mock('hooks/queries', () => ({
  useMoveTaskOrderQueries: jest.fn(),
}));

const unapprovedMTOQuery = {
  orders: {
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
      status: MOVE_STATUSES.SUBMITTED,
    },
  },
  mtoShipments: [
    {
      id: '3',
      shipmentType: SHIPMENT_OPTIONS.HHG,
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
      status: shipmentStatuses.SUBMITTED,
    },
    {
      id: '4',
      shipmentType: SHIPMENT_OPTIONS.NTS,
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
      status: shipmentStatuses.SUBMITTED,
    },
  ],
  mtoServiceItems: undefined,
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const someShipmentsApprovedMTOQuery = {
  orders: {
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
      status: MOVE_STATUSES.APPROVALS_REQUESTED,
    },
  },
  mtoShipments: [
    {
      id: '3',
      moveTaskOrderID: '2',
      shipmentType: SHIPMENT_OPTIONS.HHG,
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
      status: shipmentStatuses.APPROVED,
    },
    {
      id: '4',
      shipmentType: SHIPMENT_OPTIONS.NTS,
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
      status: shipmentStatuses.SUBMITTED,
    },
  ],
  mtoServiceItems: [
    {
      id: '5',
      mtoShipmentID: '3',
      reServiceName: 'Domestic origin 1st day SIT',
      status: SERVICE_ITEM_STATUS.SUBMITTED,
      reServiceCode: 'DOFSIT',
    },
    {
      id: '6',
      mtoShipmentID: '3',
      reServiceName: 'Domestic Linehaul',
      status: SERVICE_ITEM_STATUS.APPROVED,
      reServiceCode: 'DLH',
    },
    {
      id: '7',
      mtoShipmentID: '3',
      reServiceName: 'Domestic Unpacking',
      status: SERVICE_ITEM_STATUS.REJECTED,
      reServiceCode: 'DUPK',
    },
    {
      id: '8',
      reServiceName: 'Move management',
      status: SERVICE_ITEM_STATUS.APPROVED,
      reServiceCode: 'MS',
    },
  ],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const allApprovedMTOQuery = {
  orders: {
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
      status: MOVE_STATUSES.APPROVALS_REQUESTED,
      availableToPrimeAt: '2020-03-01T00:00:00.000Z',
    },
  },
  mtoShipments: [
    {
      id: '3',
      moveTaskOrderID: '2',
      shipmentType: SHIPMENT_OPTIONS.HHG,
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
    {
      id: '4',
      moveTaskOrderID: '2',
      shipmentType: SHIPMENT_OPTIONS.NTS,
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
    {
      id: '5',
      mtoShipmentID: '2',
      shipmentType: SHIPMENT_OPTIONS.NTSR,
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
    {
      id: '6',
      shipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
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
    {
      id: '7',
      shipmentType: SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
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
  ],
  mtoServiceItems: [
    {
      id: '8',
      mtoShipmentID: '3',
      reServiceName: 'Domestic origin 1st day SIT',
      status: SERVICE_ITEM_STATUS.SUBMITTED,
      reServiceCode: 'DOFSIT',
    },
    {
      id: '9',
      mtoShipmentID: '4',
      reServiceName: "Domestic origin add'l SIT",
      status: SERVICE_ITEM_STATUS.SUBMITTED,
      reServiceCode: 'DOASIT',
    },
  ],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const setUnapprovedShipmentCount = jest.fn();
const setUnapprovedServiceItemCount = jest.fn();

describe('MoveTaskOrder', () => {
  const moveCode = 'WE31AZ';
  const requiredProps = {
    match: { params: { moveCode } },
    history: { push: jest.fn() },
    setMessage: jest.fn(),
  };

  describe('move is not available to prime', () => {
    useMoveTaskOrderQueries.mockImplementation(() => unapprovedMTOQuery);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
        />
      </MockProviders>,
    );

    it('should render the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
      expect(wrapper.find('h1').text()).toBe('Move task order');
    });

    it('should display empty state message', () => {
      expect(
        wrapper
          .find('[data-testid="too-shipment-container"] p')
          .contains('This move does not have any approved shipments yet.'),
      ).toBe(true);
    });

    it('updates the unapproved shipments tag state', () => {
      expect(setUnapprovedShipmentCount).toHaveBeenCalledWith(2);
    });

    it('updates the unapproved shipments tag state', () => {
      expect(setUnapprovedServiceItemCount).toHaveBeenCalledWith(0);
    });
  });

  describe('approved mto with both submitted and approved shipments', () => {
    useMoveTaskOrderQueries.mockImplementation(() => someShipmentsApprovedMTOQuery);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
        />
      </MockProviders>,
    );

    it('should render the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
      expect(wrapper.find('h1').text()).toBe('Move task order');
    });

    it('should render the left nav with shipments', () => {
      expect(wrapper.find('LeftNav').exists()).toBe(true);

      const navLinks = wrapper.find('LeftNav a');
      expect(navLinks.length).toBe(1);
      expect(navLinks.at(0).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(0).prop('href')).toBe('#shipment-3');
    });

    it('should render the ShipmentContainer', () => {
      expect(wrapper.find('ShipmentContainer').length).toBe(1);
    });

    it('should render the ShipmentHeading', () => {
      expect(wrapper.find('ShipmentHeading').exists()).toBe(true);
      expect(wrapper.find('h3').at(0).text()).toEqual('Household goods');
      expect(wrapper.find('[data-testid="button"]').exists()).toBe(true);
    });

    it('should render the ImportantShipmentDates', () => {
      expect(wrapper.find('ImportantShipmentDates').exists()).toBe(true);
    });

    it('should render the ShipmentAddresses', () => {
      expect(wrapper.find('ShipmentAddresses').exists()).toBe(true);
    });

    it('should render the ShipmentWeightDetails', () => {
      expect(wrapper.find('ShipmentWeightDetails').exists()).toBe(true);
    });

    it('should render the RequestedServiceItemsTable for requested, approved, and rejected service items', () => {
      const requestedServiceItemsTable = wrapper.find('RequestedServiceItemsTable');
      // There should be 1 of each status table requested, approved, rejected service items
      expect(requestedServiceItemsTable.length).toBe(3);
      expect(requestedServiceItemsTable.at(0).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.SUBMITTED);
      expect(requestedServiceItemsTable.at(1).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.APPROVED);
      expect(requestedServiceItemsTable.at(2).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.REJECTED);
    });

    it('updates the unapproved shipments tag state', () => {
      expect(setUnapprovedShipmentCount).toHaveBeenCalledWith(0);
    });

    it('updates the unapproved service items tag state', () => {
      expect(setUnapprovedServiceItemCount).toHaveBeenCalledWith(1);
    });
  });

  describe('approved mto with approved shipments', () => {
    useMoveTaskOrderQueries.mockImplementation(() => allApprovedMTOQuery);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
        />
      </MockProviders>,
    );

    it('should render the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
      expect(wrapper.find('h1').text()).toBe('Move task order');
    });

    it('should render the left nav with shipments', () => {
      expect(wrapper.find('LeftNav').exists()).toBe(true);

      const navLinks = wrapper.find('LeftNav a');
      expect(navLinks.at(0).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(0).contains('1'));
      expect(navLinks.at(0).prop('href')).toBe('#shipment-3');

      expect(navLinks.at(1).contains('NTS shipment')).toBe(true);
      expect(navLinks.at(1).contains('1'));
      expect(navLinks.at(1).prop('href')).toBe('#shipment-4');

      expect(navLinks.at(2).contains('NTS-R shipment')).toBe(true);
      expect(navLinks.at(2).prop('href')).toBe('#shipment-5');

      expect(navLinks.at(3).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(3).prop('href')).toBe('#shipment-6');

      expect(navLinks.at(4).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(4).prop('href')).toBe('#shipment-7');
    });

    it('should render the ShipmentContainer', () => {
      expect(wrapper.find('ShipmentContainer').length).toBe(5);
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

    it('should render the ShipmentWeightDetails', () => {
      expect(wrapper.find('ShipmentWeightDetails').exists()).toBe(true);
    });

    it('should render the RequestedServiceItemsTable for SUBMITTED service item', () => {
      const requestedServiceItemsTable = wrapper.find('RequestedServiceItemsTable');
      // There are no approved or rejected service item tables to display
      expect(requestedServiceItemsTable.length).toBe(2);
      expect(requestedServiceItemsTable.at(0).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.SUBMITTED);
      expect(requestedServiceItemsTable.at(1).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.SUBMITTED);
    });

    it('updates the unapproved shipments tag state', () => {
      expect(setUnapprovedShipmentCount).toHaveBeenCalledWith(0);
    });

    it('updates the unapproved service items tag state', () => {
      expect(setUnapprovedServiceItemCount).toHaveBeenCalledWith(2);
    });
  });
});
