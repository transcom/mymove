/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';

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
      entitlement: {
        authorizedWeight: 8000,
        totalWeight: 8500,
      },
    },
  },
  move: {
    id: '2',
    status: MOVE_STATUSES.SUBMITTED,
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
      status: shipmentStatuses.SUBMITTED,
      eTag: '1234',
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
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
      status: shipmentStatuses.SUBMITTED,
      eTag: '1234',
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
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
      entitlement: {
        authorizedWeight: 8000,
        totalWeight: 8500,
      },
    },
  },
  move: {
    id: '2',
    status: MOVE_STATUSES.APPROVALS_REQUESTED,
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
        eTag: '1234',
      },
      destinationAddress: {
        street_address_1: '10 Park Place',
        city: 'Atlantic City',
        state: 'NJ',
        postal_code: '08401',
      },
      status: shipmentStatuses.APPROVED,
      eTag: '1234',
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
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
      status: shipmentStatuses.SUBMITTED,
      eTag: '1234',
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
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
      entitlement: {
        authorizedWeight: 8000,
        totalWeight: 8500,
      },
    },
  },
  move: {
    id: '2',
    status: MOVE_STATUSES.APPROVALS_REQUESTED,
    availableToPrimeAt: '2020-03-01T00:00:00.000Z',
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
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 100,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
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
      eTag: '1234',
      primeEstimatedWeight: null,
      primeActualWeight: null,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
    },
    {
      id: '5',
      moveTaskOrderID: '2',
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
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 100,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
    },
    {
      id: '6',
      moveTaskOrderID: '2',
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
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 50,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
    },
    {
      id: '7',
      moveTaskOrderID: '2',
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
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 100,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
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

// weights returned are all null
const missingWeightQuery = {
  ...allApprovedMTOQuery,
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
      eTag: '1234',
      primeEstimatedWeight: null,
      primeActualWeight: null,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
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
      eTag: '1234',
      primeEstimatedWeight: null,
      primeActualWeight: null,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
    },
    {
      id: '5',
      moveTaskOrderID: '2',
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
      eTag: '1234',
      primeEstimatedWeight: null,
      primeActualWeight: null,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
    },
  ],
};

// weight on some shipments doesn't exist
const missingSomeWeightQuery = {
  ...allApprovedMTOQuery,
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
      eTag: '1234',
      primeEstimatedWeight: null,
      primeActualWeight: null,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
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
      eTag: '1234',
      primeEstimatedWeight: 25,
      primeActualWeight: 25,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
    },
    {
      id: '5',
      moveTaskOrderID: '2',
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
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 100,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
    },
  ],
};

// weight is not returned in payload at all
const noWeightQuery = {
  ...allApprovedMTOQuery,
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
      eTag: '1234',
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
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
      eTag: '1234',
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
    },
    {
      id: '5',
      moveTaskOrderID: '2',
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
      eTag: '1234',
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
    },
  ],
};

// primeEstimatedWeight and estimatedWeightTotal is returned for some shipments but missing in others
const someWeightNotReturned = {
  ...allApprovedMTOQuery,
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
      eTag: '1234',
      primeActualWeight: 100,
      primeEstimatedWeight: 100,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
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
      eTag: '1234',
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
    },
    {
      id: '5',
      moveTaskOrderID: '2',
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
      eTag: '1234',
      primeActualWeight: 1,
      primeEstimatedWeight: 1,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
    },
  ],
};

const riskOfExcessWeightQuery = {
  ...allApprovedMTOQuery,
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
      entitlement: {
        authorizedWeight: 100,
        totalWeight: 100,
      },
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
      eTag: '1234',
      primeEstimatedWeight: 50,
      primeActualWeight: 50,
    },
    {
      id: '5',
      moveTaskOrderID: '2',
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
      eTag: '1234',
      primeEstimatedWeight: 40,
      primeActualWeight: 40,
    },
  ],
};

const approvedMTOWithCancelledShipmentQuery = {
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
      entitlement: {
        authorizedWeight: 8000,
        totalWeight: 8500,
      },
    },
  },
  move: {
    id: '2',
    status: MOVE_STATUSES.APPROVED,
    availableToPrimeAt: '2020-03-01T00:00:00.000Z',
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
      status: 'CANCELED',
      eTag: '1234',
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
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
  ],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const setUnapprovedShipmentCount = jest.fn();
const setUnapprovedServiceItemCount = jest.fn();
const setExcessWeightRiskCount = jest.fn();

const moveCode = 'WE31AZ';
const requiredProps = {
  match: { params: { moveCode } },
  history: { push: jest.fn() },
  setMessage: jest.fn(),
};

const loadingReturnValue = {
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  isLoading: false,
  isError: true,
  isSuccess: false,
};

describe('MoveTaskOrder', () => {
  describe('weight display', () => {
    it('displays the weight allowance', async () => {
      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[0]).toHaveTextContent('8,500 lbs');
    });

    it('displays the max billable weight', async () => {
      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[2]).toHaveTextContent('8,000 lbs');
    });

    it('displays the estimated total weight with all weights not set', async () => {
      useMoveTaskOrderQueries.mockReturnValue(missingWeightQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[1]).toHaveTextContent('—');
    });

    it('displays the move weight total with all weights not set', async () => {
      useMoveTaskOrderQueries.mockReturnValue(missingWeightQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[3]).toHaveTextContent('—');
    });

    it('displays the estimated total weight with some weights missing', async () => {
      useMoveTaskOrderQueries.mockReturnValue(missingSomeWeightQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[1]).toHaveTextContent('125 lbs');
    });

    it('displays the move weight total with some weights missing', async () => {
      useMoveTaskOrderQueries.mockReturnValue(missingSomeWeightQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[3]).toHaveTextContent('125 lbs');
    });

    it('displays the estimated total weight with all not sent', async () => {
      useMoveTaskOrderQueries.mockReturnValue(noWeightQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[1]).toHaveTextContent('—');
    });

    it('displays the move weight total with all not sent', async () => {
      useMoveTaskOrderQueries.mockReturnValue(noWeightQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[3]).toHaveTextContent('—');
    });

    it('displays the estimated total weight with some sent and some not sent', async () => {
      useMoveTaskOrderQueries.mockReturnValue(someWeightNotReturned);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[1]).toHaveTextContent('101');
    });

    it('displays the move weight total with some sent and some not sent', async () => {
      useMoveTaskOrderQueries.mockReturnValue(someWeightNotReturned);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const weightSummaries = await screen.findAllByTestId('weight-display');
      expect(weightSummaries[3]).toHaveTextContent('101');
    });

    it('displays risk of excess tag', async () => {
      useMoveTaskOrderQueries.mockReturnValue(riskOfExcessWeightQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const riskOfExcessTag = await screen.getByText(/Risk of excess/);
      expect(riskOfExcessTag).toBeInTheDocument();
    });

    it('displays the estimated total weight', async () => {
      useMoveTaskOrderQueries.mockReturnValue(allApprovedMTOQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const estimatedWeightTotal = await screen.getByText(/400 lbs/);
      expect(estimatedWeightTotal).toBeInTheDocument();
    });

    it('displays the move weight total', async () => {
      useMoveTaskOrderQueries.mockReturnValue(allApprovedMTOQuery);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const moveWeightTotal = await screen.getByText(/350 lbs/);
      expect(moveWeightTotal).toBeInTheDocument();
    });
  });
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useMoveTaskOrderQueries.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useMoveTaskOrderQueries.mockReturnValue(errorReturnValue);

      render(
        <MockProviders initialEntries={['moves/1000/allowances']}>
          <MoveTaskOrder
            {...requiredProps}
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
          />
        </MockProviders>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('move is not available to prime', () => {
    useMoveTaskOrderQueries.mockReturnValue(unapprovedMTOQuery);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
        />
      </MockProviders>,
    );

    it('renders the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
      expect(wrapper.find('h1').text()).toBe('Move task order');
    });

    it('displays empty state message', () => {
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
    useMoveTaskOrderQueries.mockReturnValue(someShipmentsApprovedMTOQuery);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
        />
      </MockProviders>,
    );

    it('renders the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
      expect(wrapper.find('h1').text()).toBe('Move task order');
    });

    it('renders the left nav with shipments', () => {
      expect(wrapper.find('LeftNav').exists()).toBe(true);

      const navLinks = wrapper.find('LeftNav a');
      expect(navLinks.length).toBe(1);
      expect(navLinks.at(0).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(0).prop('href')).toBe('#shipment-3');
    });

    it('renders the ShipmentContainer', () => {
      expect(wrapper.find('ShipmentContainer').length).toBe(1);
    });

    it('renders the ShipmentHeading', () => {
      expect(wrapper.find('ShipmentHeading').exists()).toBe(true);
      expect(wrapper.find('h2').at(0).text()).toEqual('Household goods');
      expect(wrapper.find('[data-testid="button"]').exists()).toBe(true);
    });

    it('renders the ImportantShipmentDates', () => {
      expect(wrapper.find('ImportantShipmentDates').exists()).toBe(true);
    });

    it('renders the ShipmentAddresses', () => {
      expect(wrapper.find('ShipmentAddresses').exists()).toBe(true);
    });

    it('renders the ShipmentWeightDetails', () => {
      expect(wrapper.find('ShipmentWeightDetails').exists()).toBe(true);
    });

    it('renders the RequestedServiceItemsTable for requested, approved, and rejected service items', () => {
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
    useMoveTaskOrderQueries.mockReturnValue(allApprovedMTOQuery);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
        />
      </MockProviders>,
    );

    it('renders the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
      expect(wrapper.find('h1').text()).toBe('Move task order');
    });

    it('renders the left nav with shipments', () => {
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

    it('renders the ShipmentContainer', () => {
      expect(wrapper.find('ShipmentContainer').length).toBe(5);
    });

    it('renders the ShipmentHeading', () => {
      expect(wrapper.find('ShipmentHeading').exists()).toBe(true);
      expect(wrapper.find('h2').at(0).text()).toEqual('Household goods');
      expect(wrapper.find('h2').at(1).text()).toEqual('Non-temp storage');
    });

    it('renders the ImportantShipmentDates', () => {
      expect(wrapper.find('ImportantShipmentDates').exists()).toBe(true);
    });

    it('renders the ShipmentAddresses', () => {
      expect(wrapper.find('ShipmentAddresses').exists()).toBe(true);
    });

    it('renders the ShipmentWeightDetails', () => {
      expect(wrapper.find('ShipmentWeightDetails').exists()).toBe(true);
    });

    it('renders the RequestedServiceItemsTable for SUBMITTED service item', () => {
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

  describe('approved mto with cancelled shipment', () => {
    useMoveTaskOrderQueries.mockReturnValue(approvedMTOWithCancelledShipmentQuery);
    const wrapper = mount(
      <MockProviders>
        <MoveTaskOrder
          {...requiredProps}
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
        />
      </MockProviders>,
    );

    it('renders the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-shipment-container' }).exists()).toBe(true);
      expect(wrapper.find('h1').text()).toBe('Move task order');
    });

    it('renders the left nav with shipments', () => {
      expect(wrapper.find('LeftNav').exists()).toBe(true);

      const navLinks = wrapper.find('LeftNav a');
      expect(navLinks.at(0).contains('HHG shipment')).toBe(true);
      expect(navLinks.at(0).contains('1'));
      expect(navLinks.at(0).prop('href')).toBe('#shipment-3');
    });

    it('renders the ShipmentContainer', () => {
      expect(wrapper.find('ShipmentContainer').length).toBe(1);
    });

    it('renders the ShipmentHeading', () => {
      expect(wrapper.find('ShipmentHeading').exists()).toBe(true);
      expect(wrapper.find('h2').at(0).text()).toEqual('Household goods');
      expect(wrapper.find('span[data-testid="tag"]').at(0).text()).toEqual('cancelled');
    });

    it('renders the ImportantShipmentDates', () => {
      expect(wrapper.find('ImportantShipmentDates').exists()).toBe(true);
    });

    it('renders the ShipmentAddresses', () => {
      expect(wrapper.find('ShipmentAddresses').exists()).toBe(true);
    });

    it('renders the ShipmentWeightDetails', () => {
      expect(wrapper.find('ShipmentWeightDetails').exists()).toBe(true);
      expect(wrapper.find('span[data-testid="tag"]').at(1).text()).toEqual('reweigh requested');
    });

    it('renders the RequestedServiceItemsTable for SUBMITTED service item', () => {
      const requestedServiceItemsTable = wrapper.find('RequestedServiceItemsTable');
      // There are no approved or rejected service item tables to display
      expect(requestedServiceItemsTable.length).toBe(1);
      expect(requestedServiceItemsTable.at(0).prop('statusForTableType')).toBe(SERVICE_ITEM_STATUS.SUBMITTED);
    });

    it('updates the unapproved shipments tag state', () => {
      expect(setUnapprovedShipmentCount).toHaveBeenCalledWith(0);
    });

    it('updates the unapproved service items tag state', () => {
      expect(setUnapprovedServiceItemCount).toHaveBeenCalledWith(2);
    });
  });
});
