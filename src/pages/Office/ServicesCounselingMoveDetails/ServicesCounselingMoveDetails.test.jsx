/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import ServicesCounselingMoveDetails from './ServicesCounselingMoveDetails';

import { ORDERS_TYPE, ORDERS_TYPE_DETAILS } from 'constants/orders';
import { MockProviders } from 'testUtils';
import { useMoveDetailsQueries } from 'hooks/queries';
import MOVE_STATUSES from 'constants/moves';

const mockRequestedMoveCode = 'LR4T8V';

jest.mock('hooks/queries', () => ({
  useMoveDetailsQueries: jest.fn(),
}));

const newMoveDetailsQuery = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
    status: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING,
  },
  order: {
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
        postal_code: '90210',
        state: 'CA',
        street_address_1: '123 Any Street',
        street_address_2: 'P.O. Box 12345',
        street_address_3: 'c/o Some Person',
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
  },
  mtoShipments: [
    {
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
  ],
  mtoServiceItems: [],
  mtoAgents: [],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const counselingCompletedMoveDetailsQuery = {
  ...newMoveDetailsQuery,
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
    status: MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED,
  },
};

describe('MoveDetails page', () => {
  it('renders the h1', () => {
    const wrapper = mount(
      <MockProviders initialEntries={[`counseling/moves/${mockRequestedMoveCode}/details`]}>
        <ServicesCounselingMoveDetails />
      </MockProviders>,
    );
    expect(wrapper.find({ 'data-testid': 'sc-move-details' }).exists()).toBe(true);
    expect(wrapper.containsMatchingElement(<h1>Move details</h1>)).toBe(true);
  });

  describe('new move - needs service counseling', () => {
    useMoveDetailsQueries.mockImplementation(() => newMoveDetailsQuery);

    const wrapper = mount(
      <MockProviders initialEntries={[`counseling/moves/${mockRequestedMoveCode}/details`]}>
        <ServicesCounselingMoveDetails />
      </MockProviders>,
    );

    it('submit move details button is on page', () => {
      expect(wrapper.find('button[data-testid="submitMoveDetailsBtn"]').length).toBe(1);
    });
  });

  describe('service counseling completed', () => {
    useMoveDetailsQueries.mockImplementation(() => counselingCompletedMoveDetailsQuery);

    const wrapper = mount(
      <MockProviders initialEntries={[`counseling/moves/${mockRequestedMoveCode}/details`]}>
        <ServicesCounselingMoveDetails />
      </MockProviders>,
    );

    it('submit move details button is not on page', () => {
      expect(wrapper.find('[data-testid="submitMoveDetailsBtn"]').length).toBe(0);
    });
  });
});
