/* eslint-disable import/prefer-default-export */
import MOVE_STATUSES from 'constants/moves';
import { shipmentStatuses } from 'constants/shipments';
import { SERVICE_ITEM_STATUSES } from 'constants/serviceItems';
import { ORDERS_TYPE } from 'constants/orders';

const move = {
  id: '1',
  contractor: {
    contractNumber: 'HTC-123-3456',
  },
  orders: {
    sac: '1234456',
    tac: '1213',
  },
  billableWeightsReviewedAt: '2021-06-01',
};

const order = {
  orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  has_dependents: false,
  issue_date: '2020-08-11',
  grade: 'E_1',
  moves: ['123'],
  origin_duty_location: {
    name: 'Test Duty Location',
    address: {
      postalCode: '123456',
    },
  },
  new_duty_location: {
    name: 'New Test Duty Location',
    address: {
      postalCode: '123456',
    },
  },
  report_by_date: '2020-08-31',
  service_member_id: '666',
  spouse_has_pro_gear: false,
  status: MOVE_STATUSES.SUBMITTED,
  uploaded_orders: {
    uploads: [],
  },
  entitlement: {
    authorizedWeight: 8000,
    dependentsAuthorized: true,
    eTag: 'MjAyMS0wOC0yNFQxODoyNDo0MC45NzIzMTha',
    id: '188842d1-cf88-49ec-bd2f-dfa98da44bb2',
    nonTemporaryStorage: true,
    organizationalClothingAndIndividualEquipment: true,
    privatelyOwnedVehicle: true,
    proGearWeight: 2000,
    proGearWeightSpouse: 500,
    requiredMedicalEquipmentWeight: 1000,
    storageInTransit: 2,
    totalDependents: 1,
    totalWeight: 8000,
  },
};

export const multiplePaymentRequests = {
  paymentRequests: [
    {
      id: '09474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      paymentRequestNumber: '1843-9061-1',
      status: 'REVIEWED',
      moveTaskOrderID: '1',
      moveTaskOrder: move,
      serviceItems: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 2000001,
          status: 'APPROVED',
        },
        {
          id: '39474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 4000001,
          status: 'DENIED',
          rejectionReason: 'Requested amount exceeds guideline',
        },
      ],
      reviewedAt: '2020-12-01T00:00:00.000Z',
    },
    {
      id: '29474c6a-69b6-4501-8e08-670a12512e5f',
      createdAt: '2020-12-01T00:00:00.000Z',
      paymentRequestNumber: '1843-9061-2',
      status: 'PENDING',
      moveTaskOrderID: '1',
      moveTaskOrder: move,
      serviceItems: [
        {
          id: '09474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'f8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 2000001,
          status: 'REQUESTED',
        },
        {
          id: '39474c6a-69b6-4501-8e08-670a12512a5f',
          createdAt: '2020-12-01T00:00:00.000Z',
          mtoServiceItemID: 'a8c2f97f-99e7-4fb1-9cc4-473debd24dbc',
          priceCents: 4000001,
          status: 'REQUESTED',
        },
      ],
    },
  ],
  mtoShipments: [
    {
      shipmentType: 'HHG',
      id: '2',
      moveTaskOrderID: '1',
      status: shipmentStatuses.APPROVED,
      scheduledPickupDate: '2020-01-09T00:00:00.000Z',
      destinationAddress: { city: 'Princeton', state: 'NJ', postalCode: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postalCode: '02101' },
      calculatedBillableWeight: 9000,
      primeActualWeight: 500,
      reweigh: {
        id: 'reweighID1',
        weight: 100,
      },
      mtoServiceItems: [
        {
          id: '5',
          mtoShipmentID: '2',
          status: SERVICE_ITEM_STATUSES.APPROVED,
        },
        {
          id: '6',
          status: SERVICE_ITEM_STATUSES.REJECTED,
          mtoShipmentID: '2',
        },
        {
          id: '7',
          status: SERVICE_ITEM_STATUSES.SUBMITTED,
          mtoShipmentID: '2',
        },
      ],
    },
    {
      shipmentType: 'HHG',
      id: '3',
      moveTaskOrderID: '1',
      status: shipmentStatuses.APPROVED,
      scheduledPickupDate: '2020-01-10T00:00:00.000Z',
      destinationAddress: { city: 'Princeton', state: 'NJ', postalCode: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postalCode: '02101' },
      calculatedBillableWeight: 1000,
      primeActualWeight: 5000,
      reweigh: {
        id: 'reweighID2',
        weight: 600,
      },
      mtoServiceItems: [
        {
          id: '9',
          mtoShipmentID: '3',
          status: SERVICE_ITEM_STATUSES.APPROVED,
        },
        {
          id: '10',
          status: SERVICE_ITEM_STATUSES.REJECTED,
          mtoShipmentID: '3',
        },
        {
          id: '11',
          status: SERVICE_ITEM_STATUSES.SUBMITTED,
          mtoShipmentID: '3',
        },
      ],
    },
    {
      shipmentType: 'HHG',
      id: '4',
      moveTaskOrderID: '1',
      status: shipmentStatuses.SUBMITTED,
      scheduledPickupDate: '2020-01-11T00:00:00.000Z',
      destinationAddress: { city: 'Princeton', state: 'NJ', postalCode: '08540' },
      pickupAddress: { city: 'Boston', state: 'MA', postalCode: '02101' },
      calculatedBillableWeight: 2000,
      primeActualWeight: 300,
      reweigh: {
        id: 'reweighID3',
        weight: 900,
      },
      mtoServiceItems: [
        {
          id: '12',
          mtoShipmentID: '4',
          status: SERVICE_ITEM_STATUSES.APPROVED,
        },
        {
          id: '13',
          status: SERVICE_ITEM_STATUSES.REJECTED,
          mtoShipmentID: '4',
        },
        {
          id: '14',
          status: SERVICE_ITEM_STATUSES.SUBMITTED,
          mtoShipmentID: '4',
        },
      ],
    },
  ],
  order,
};
