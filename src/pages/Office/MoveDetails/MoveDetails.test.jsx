/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';

import { ORDERS_TYPE, ORDERS_TYPE_DETAILS } from '../../../constants/orders';

import MoveDetails from './MoveDetails';

import { MockProviders } from 'testUtils';
import { useMoveDetailsQueries } from 'hooks/queries';

const mockRequestedMoveCode = 'LR4T8V';

jest.mock('hooks/queries', () => ({
  useMoveDetailsQueries: jest.fn(),
}));

const setUnapprovedShipmentCount = jest.fn();
const setUnapprovedServiceItemCount = jest.fn();
const setExcessWeightRiskCount = jest.fn();
const setUnapprovedSITExtensionCount = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: () => ({ moveCode: 'TE5TC0DE' }),
}));

const requestedMoveDetailsQuery = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
  },
  order: {
    id: '1',
    originDutyStation: {
      address: {
        streetAddress1: '',
        city: 'Fort Knox',
        state: 'KY',
        postalCode: '40121',
      },
    },
    destinationDutyStation: {
      address: {
        streetAddress1: '',
        city: 'Fort Irwin',
        state: 'CA',
        postalCode: '92310',
      },
    },
    customer: {
      agency: 'ARMY',
      backup_contact: {
        email: 'email@example.com',
        name: 'name',
        phone: '555-555-5555',
      },
      current_address: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41Mzg0Njha',
        id: '3a5f7cf2-6193-4eb3-a244-14d21ca05d7b',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      dodID: '6833908165',
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NjAzNTJa',
      email: 'combo@ppm.hhg',
      first_name: 'Submitted',
      id: 'f6bd793f-7042-4523-aa30-34946e7339c9',
      last_name: 'Ppmhhg',
      phone: '555-555-5555',
    },
    entitlement: {
      authorizedWeight: 8000,
      dependentsAuthorized: true,
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NzgwMzda',
      id: 'e0fefe58-0710-40db-917b-5b96567bc2a8',
      nonTemporaryStorage: true,
      privatelyOwnedVehicle: true,
      proGearWeight: 2000,
      proGearWeightSpouse: 500,
      storageInTransit: 2,
      totalDependents: 1,
      totalWeight: 8000,
    },
    order_number: 'ORDER3',
    order_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    order_type_detail: ORDERS_TYPE_DETAILS.HHG_PERMITTED,
    tac: '9999',
    ntsTac: '1111',
    ntsSac: '2222',
  },
  mtoShipments: [
    {
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
    {
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
      deletedAt: '2018-03-16',
    },
  ],
  mtoServiceItems: [],
  mtoAgents: [],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const requestedMoveDetailsQueryRetiree = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
  },
  order: {
    id: '1',
    originDutyStation: {
      address: {
        streetAddress1: '',
        city: 'Fort Knox',
        state: 'KY',
        postalCode: '40121',
      },
    },
    destinationDutyStation: {
      address: {
        streetAddress1: '',
        city: 'Fort Irwin',
        state: 'CA',
        postalCode: '92310',
      },
    },
    customer: {
      agency: 'ARMY',
      backup_contact: {
        email: 'email@example.com',
        name: 'name',
        phone: '555-555-5555',
      },
      current_address: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41Mzg0Njha',
        id: '3a5f7cf2-6193-4eb3-a244-14d21ca05d7b',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      dodID: '6833908165',
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NjAzNTJa',
      email: 'combo@ppm.hhg',
      first_name: 'Submitted',
      id: 'f6bd793f-7042-4523-aa30-34946e7339c9',
      last_name: 'Ppmhhg',
      phone: '555-555-5555',
    },
    entitlement: {
      authorizedWeight: 8000,
      dependentsAuthorized: true,
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NzgwMzda',
      id: 'e0fefe58-0710-40db-917b-5b96567bc2a8',
      nonTemporaryStorage: true,
      privatelyOwnedVehicle: true,
      proGearWeight: 2000,
      proGearWeightSpouse: 500,
      storageInTransit: 2,
      totalDependents: 1,
      totalWeight: 8000,
    },
    order_number: 'ORDER3',
    order_type: ORDERS_TYPE.RETIREMENT,
    order_type_detail: ORDERS_TYPE_DETAILS.HHG_PERMITTED,
    tac: '9999',
    ntsTac: '1111',
    ntsSac: '2222',
  },
  mtoShipments: [
    {
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
    {
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
      deletedAt: '2018-03-16',
    },
  ],
  mtoServiceItems: [],
  mtoAgents: [],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const requestedMoveDetailsAmendedOrdersQuery = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
  },
  order: {
    id: '1',
    originDutyStation: {
      address: {
        streetAddress1: '',
        city: 'Fort Knox',
        state: 'KY',
        postalCode: '40121',
      },
    },
    destinationDutyStation: {
      address: {
        streetAddress1: '',
        city: 'Fort Irwin',
        state: 'CA',
        postalCode: '92310',
      },
    },
    customer: {
      agency: 'ARMY',
      backup_contact: {
        email: 'email@example.com',
        name: 'name',
        phone: '555-555-5555',
      },
      current_address: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41Mzg0Njha',
        id: '3a5f7cf2-6193-4eb3-a244-14d21ca05d7b',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      dodID: '6833908165',
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NjAzNTJa',
      email: 'combo@ppm.hhg',
      first_name: 'Submitted',
      id: 'f6bd793f-7042-4523-aa30-34946e7339c9',
      last_name: 'Ppmhhg',
      phone: '555-555-5555',
    },
    entitlement: {
      authorizedWeight: 8000,
      dependentsAuthorized: true,
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NzgwMzda',
      id: 'e0fefe58-0710-40db-917b-5b96567bc2a8',
      nonTemporaryStorage: true,
      privatelyOwnedVehicle: true,
      proGearWeight: 2000,
      proGearWeightSpouse: 500,
      storageInTransit: 2,
      totalDependents: 1,
      totalWeight: 8000,
    },
    order_number: 'ORDER3',
    order_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
    order_type_detail: ORDERS_TYPE_DETAILS.HHG_PERMITTED,
    uploadedAmendedOrderID: '3',
    tac: '9999',
  },
  mtoShipments: [
    {
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
    {
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
      deletedAt: '2018-03-16',
    },
  ],
  mtoServiceItems: [],
  mtoAgents: [],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const requestedMoveDetailsMissingInfoQuery = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
  },
  order: {
    id: '1',
    originDutyStation: {
      address: {
        streetAddress1: '',
        city: 'Fort Knox',
        state: 'KY',
        postalCode: '40121',
      },
    },
    destinationDutyStation: {
      address: {
        streetAddress1: '',
        city: 'Fort Irwin',
        state: 'CA',
        postalCode: '92310',
      },
    },
    customer: {
      agency: 'ARMY',
      backup_contact: {
        email: 'email@example.com',
        name: 'name',
        phone: '555-555-5555',
      },
      current_address: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41Mzg0Njha',
        id: '3a5f7cf2-6193-4eb3-a244-14d21ca05d7b',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      dodID: '6833908165',
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NjAzNTJa',
      email: 'combo@ppm.hhg',
      first_name: 'Submitted',
      id: 'f6bd793f-7042-4523-aa30-34946e7339c9',
      last_name: 'Ppmhhg',
      phone: '555-555-5555',
    },
    entitlement: {
      authorizedWeight: 8000,
      dependentsAuthorized: true,
      eTag: 'MjAyMS0wMS0yMVQxNTo0MTozNS41NzgwMzda',
      id: 'e0fefe58-0710-40db-917b-5b96567bc2a8',
      nonTemporaryStorage: true,
      privatelyOwnedVehicle: true,
      proGearWeight: 2000,
      proGearWeightSpouse: 500,
      storageInTransit: 2,
      totalDependents: 1,
      totalWeight: 8000,
    },
  },
  mtoShipments: [
    {
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
    {
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4abf',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b7d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-ca1467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
  ],
  mtoServiceItems: [],
  mtoAgents: [],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const requestedAndApprovedMoveDetailsQuery = {
  ...requestedMoveDetailsQuery,
  mtoShipments: [
    {
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
    {
      approvedDate: '2020-01-01',
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'APPROVED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
  ],
};

const approvedMoveDetailsQuery = {
  ...requestedMoveDetailsQuery,
  mtoShipments: [
    {
      approvedDate: '2020-01-01',
      customerRemarks: 'please treat gently',
      destinationAddress: {
        city: 'Fairfield',
        country: 'US',
        id: '672ff379-f6e3-48b4-a87d-796713f8f997',
        postalCode: '94535',
        state: 'CA',
        streetAddress1: '987 Any Avenue',
        streetAddress2: 'P.O. Box 9876',
        streetAddress3: 'c/o Some Person',
      },
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
      id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
      moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
      pickupAddress: {
        city: 'Beverly Hills',
        country: 'US',
        eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
        id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
        postalCode: '90210',
        state: 'CA',
        streetAddress1: '123 Any Street',
        streetAddress2: 'P.O. Box 12345',
        streetAddress3: 'c/o Some Person',
      },
      requestedPickupDate: '2018-03-15',
      scheduledPickupDate: '2018-03-16',
      shipmentType: 'HHG',
      status: 'APPROVED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
  ],
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

describe('MoveDetails page', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useMoveDetailsQueries.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders initialEntries={[`/moves/${mockRequestedMoveCode}/details`]}>
          <MoveDetails
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useMoveDetailsQueries.mockReturnValue(errorReturnValue);

      render(
        <MockProviders initialEntries={[`/moves/${mockRequestedMoveCode}/details`]}>
          <MoveDetails
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });
  describe('requested shipment', () => {
    useMoveDetailsQueries.mockReturnValue(requestedMoveDetailsQuery);

    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${mockRequestedMoveCode}/details`]}>
        <MoveDetails
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
          setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
        />
      </MockProviders>,
    );

    it('renders the h1', () => {
      expect(wrapper.find({ 'data-testid': 'too-move-details' }).exists()).toBe(true);
      expect(wrapper.containsMatchingElement(<h1>Move details</h1>)).toBe(true);
    });

    it('renders side navigation for each section', () => {
      expect(wrapper.find('LeftNav').exists()).toBe(true);

      const navLinks = wrapper.find('LeftNav a');

      expect(navLinks.at(0).contains('Requested shipments')).toBe(true);
      expect(navLinks.at(0).contains(1)).toBe(true);
      expect(navLinks.at(0).prop('href')).toBe('#requested-shipments');

      expect(navLinks.at(1).contains('Orders')).toBe(true);
      expect(navLinks.at(1).prop('href')).toBe('#orders');

      expect(navLinks.at(2).contains('Allowances')).toBe(true);
      expect(navLinks.at(2).prop('href')).toBe('#allowances');

      expect(navLinks.at(3).contains('Customer info')).toBe(true);
      expect(navLinks.at(3).prop('href')).toBe('#customer-info');
    });

    it('renders the Requested Shipments component', () => {
      expect(wrapper.find('RequestedShipments')).toHaveLength(1);
    });

    it('renders the Orders Table', () => {
      expect(wrapper.find('#orders h2').text()).toEqual('Orders');
      expect(wrapper.find('dd[data-testid="NTStac"]').text()).toEqual('1111');
      expect(wrapper.find('dd[data-testid="NTSsac"]').text()).toEqual('2222');
    });

    it('renders the Allowances Table', () => {
      expect(wrapper.find('#allowances h2').text()).toEqual('Allowances');
    });

    it('renders the Customer Info Table', () => {
      expect(wrapper.find('#customer-info h2').text()).toEqual('Customer info');
    });

    it('renders the requested shipments tag', () => {
      expect(wrapper.find('span[data-testid="requestedShipmentsTag"]').text()).toEqual('1');
    });

    it('updates the unapproved shipments tag state', () => {
      expect(setUnapprovedShipmentCount).toHaveBeenCalledWith(1);
      expect(setUnapprovedShipmentCount.mock.calls[0][0]).toBe(1);
    });
  });

  describe('retiree move with shipment', () => {
    useMoveDetailsQueries.mockReturnValue(requestedMoveDetailsQueryRetiree);

    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${mockRequestedMoveCode}/details`]}>
        <MoveDetails
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
          setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
        />
      </MockProviders>,
    );
    it('renders the Orders Table', () => {
      expect(wrapper.find('#orders h2').text()).toEqual('Orders');
      expect(wrapper.find('[data-testid="newDutyLocationLabel"]').text()).toEqual('HOR, HOS, or PLEAD');
      expect(wrapper.find('[data-testid="reportByDateLabel"]').text()).toEqual('Date of retirement');
    });
  });

  describe('requested shipment with amended orders', () => {
    useMoveDetailsQueries.mockReturnValue(requestedMoveDetailsAmendedOrdersQuery);

    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${mockRequestedMoveCode}/details`]}>
        <MoveDetails
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
          setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
        />
      </MockProviders>,
    );

    it('renders a NEW tag in the orders navigation section', () => {
      expect(wrapper.find('[data-testid="newOrdersNavTag"]').exists()).toBe(true);
    });

    it('renders the Orders Table with NEW tag', () => {
      expect(wrapper.find('[data-testid="detailsPanelTag"]').exists()).toBe(true);
    });
  });

  describe('requested and approved shipment', () => {
    useMoveDetailsQueries.mockReturnValue(requestedAndApprovedMoveDetailsQuery);

    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${mockRequestedMoveCode}/details`]}>
        <MoveDetails
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
          setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
        />
      </MockProviders>,
    );

    it('renders side navigation for each section', () => {
      expect(wrapper.find('LeftNav').exists()).toBe(true);

      const navLinks = wrapper.find('LeftNav a');

      expect(navLinks.at(0).contains('Requested shipments')).toBe(true);
      expect(navLinks.at(0).contains(1)).toBe(true);
      expect(navLinks.at(0).prop('href')).toBe('#requested-shipments');

      expect(navLinks.at(1).contains('Approved shipments')).toBe(true);
      expect(navLinks.at(1).prop('href')).toBe('#approved-shipments');

      expect(navLinks.at(2).contains('Orders')).toBe(true);
      expect(navLinks.at(2).prop('href')).toBe('#orders');

      expect(navLinks.at(3).contains('Allowances')).toBe(true);
      expect(navLinks.at(3).prop('href')).toBe('#allowances');

      expect(navLinks.at(4).contains('Customer info')).toBe(true);
      expect(navLinks.at(4).prop('href')).toBe('#customer-info');
    });
  });

  describe('approved shipment', () => {
    it.each([['Approved shipments'], ['Orders'], ['Allowances'], ['Customer info']])(
      'renders side navigation for section %s',
      async (sectionName) => {
        useMoveDetailsQueries.mockReturnValue(approvedMoveDetailsQuery);

        render(
          <MockProviders initialEntries={[`/moves/${mockRequestedMoveCode}/details`]}>
            <MoveDetails
              setUnapprovedShipmentCount={setUnapprovedShipmentCount}
              setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
              setExcessWeightRiskCount={setExcessWeightRiskCount}
              setUnapprovedSITExtensionCount={setUnapprovedServiceItemCount}
            />
          </MockProviders>,
        );

        expect(await screen.findByRole('link', { name: sectionName })).toBeInTheDocument();
      },
    );
  });

  describe('When required Orders information (like TAC) is missing', () => {
    useMoveDetailsQueries.mockReturnValue(requestedMoveDetailsMissingInfoQuery);

    const wrapper = mount(
      <MockProviders initialEntries={[`/moves/${mockRequestedMoveCode}/details`]}>
        <MoveDetails
          setUnapprovedShipmentCount={setUnapprovedShipmentCount}
          setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
          setExcessWeightRiskCount={setExcessWeightRiskCount}
          setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
        />
      </MockProviders>,
    );

    it('renders an error indicator in the sidebar', () => {
      expect(wrapper.find('a[href="#orders"] span[data-testid="tag"]').exists()).toBe(true);
    });
  });

  describe('When required shipment information (like TAC) is missing', () => {
    it('renders an error indicator in the sidebar', async () => {
      useMoveDetailsQueries.mockReturnValue(requestedMoveDetailsMissingInfoQuery);

      render(
        <MockProviders initialEntries={[`/moves/${mockRequestedMoveCode}/details`]}>
          <MoveDetails
            setUnapprovedShipmentCount={setUnapprovedShipmentCount}
            setUnapprovedServiceItemCount={setUnapprovedServiceItemCount}
            setExcessWeightRiskCount={setExcessWeightRiskCount}
            setUnapprovedSITExtensionCount={setUnapprovedSITExtensionCount}
          />
        </MockProviders>,
      );

      expect(await screen.findByTestId('shipment-missing-info-alert')).toBeInTheDocument();
    });
  });
});
