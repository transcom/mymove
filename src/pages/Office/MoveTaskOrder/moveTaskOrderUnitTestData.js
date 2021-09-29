/* eslint-disable import/prefer-default-export */
import MOVE_STATUSES from 'constants/moves';
import { shipmentStatuses } from 'constants/shipments';
import SERVICE_ITEM_STATUS from 'constants/serviceItems';
import { SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import { SHIPMENT_OPTIONS } from 'shared/constants';

export const unapprovedMTOQuery = {
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
      sitExtensions: [],
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
      sitExtensions: [],
    },
  ],
  mtoServiceItems: undefined,
  isLoading: false,
  isError: false,
  isSuccess: true,
};

export const someShipmentsApprovedMTOQuery = {
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
      sitExtensions: [],
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
      sitExtensions: [],
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

export const allApprovedMTOQuery = {
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
      sitExtensions: [],
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
      sitExtensions: [],
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
      sitExtensions: [],
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
      sitExtensions: [],
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
      sitExtensions: [],
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
export const missingWeightQuery = {
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
      sitExtensions: [],
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
      sitExtensions: [],
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
      sitExtensions: [],
    },
  ],
};

// weight on some shipments doesn't exist
export const missingSomeWeightQuery = {
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
      sitExtensions: [],
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
      sitExtensions: [],
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
      sitExtensions: [],
    },
  ],
};

// weight is not returned in payload at all
export const noWeightQuery = {
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
      sitExtensions: [],
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
      sitExtensions: [],
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
      sitExtensions: [],
    },
  ],
};

// primeEstimatedWeight and estimatedWeightTotal is returned for some shipments but missing in others
export const someWeightNotReturned = {
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
      sitExtensions: [],
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
      sitExtensions: [],
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
      sitExtensions: [],
    },
  ],
};

// SIT extension present on the shipment
export const sitExtensionPresent = {
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
      sitExtensions: [
        {
          mtoShipmentID: '3',
          requestReason: 'reason',
          status: SIT_EXTENSION_STATUS.PENDING,
          requestedDays: 42,
          id: '1',
        },
      ],
    },
  ],
};
// SIT extension present on the shipment that's been approved
export const sitExtensionApproved = {
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
      status: 'APPROVED',
      eTag: '1234',
      primeActualWeight: 100,
      primeEstimatedWeight: 100,
      sitExtensions: [
        {
          mtoShipmentID: '3',
          requestReason: 'reason',
          status: SIT_EXTENSION_STATUS.APPROVED,
          requestedDays: 42,
          id: '2',
        },
      ],
    },
  ],
};

export const riskOfExcessWeightQuery = {
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
      sitExtensions: [],
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
      sitExtensions: [],
    },
  ],
};

export const approvedMTOWithCancelledShipmentQuery = {
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
      sitExtensions: [],
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

export const lowerReweighsMTOQuery = {
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
      status: shipmentStatuses.APPROVED,
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 100,
      reweigh: {
        weight: 99,
      },
      sitExtensions: [],
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
      status: shipmentStatuses.APPROVED,
      eTag: '1234',
      primeEstimatedWeight: null,
      primeActualWeight: null,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
      sitExtensions: [],
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
      status: shipmentStatuses.DIVERSION_REQUESTED,
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 100,
      reweigh: {
        weight: 99,
      },
      sitExtensions: [],
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
      status: shipmentStatuses.CANCELLATION_REQUESTED,
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 50,
      reweigh: {
        weight: 49,
      },
      sitExtensions: [],
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
      status: shipmentStatuses.SUBMITTED,
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 100,
      diversion: true,
      reweigh: {
        weight: 99,
      },
      sitExtensions: [],
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
      status: shipmentStatuses.CANCELED,
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 100,
      diversion: true,
      reweigh: {
        weight: 99,
      },
      sitExtensions: [],
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

export const lowerActualsMTOQuery = {
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
        weight: 101,
      },
      sitExtensions: [],
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
      status: shipmentStatuses.APPROVED,
      eTag: '1234',
      primeEstimatedWeight: null,
      primeActualWeight: null,
      reweigh: {
        id: '00000000-0000-0000-0000-000000000000',
      },
      sitExtensions: [],
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
      status: shipmentStatuses.DIVERSION_REQUESTED,
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 100,
      reweigh: {
        weight: 101,
      },
      sitExtensions: [],
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
      status: shipmentStatuses.CANCELLATION_REQUESTED,
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 50,
      reweigh: {
        weight: 51,
      },
      sitExtensions: [],
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
      status: shipmentStatuses.SUBMITTED,
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 100,
      diversion: true,
      reweigh: {
        weight: 101,
      },
      sitExtensions: [],
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
      status: shipmentStatuses.CANCELED,
      eTag: '1234',
      primeEstimatedWeight: 100,
      primeActualWeight: 100,
      diversion: true,
      reweigh: {
        weight: 101,
      },
      sitExtensions: [],
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
