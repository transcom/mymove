/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';

import MoveAllowances from './MoveAllowances';

import { MockProviders } from 'testUtils';

const mockOriginDutyStation = {
  address: {
    city: 'Des Moines',
    country: 'US',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42OTg1OTha',
    id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
    postal_code: '50309',
    state: 'IA',
    street_address_1: '987 Other Avenue',
    street_address_2: 'P.O. Box 1234',
    street_address_3: 'c/o Another Person',
  },
  address_id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MDcxOTVa',
  id: 'a3ec2bdd-aa0a-434a-ba58-34c85f047704',
  name: 'XBc1KNi3pA',
};

const mockDestinationDutyStation = {
  address: {
    city: 'Augusta',
    country: 'United States',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
    id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    postal_code: '30813',
    state: 'GA',
    street_address_1: 'Fort Gordon',
  },
  address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
  id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
  name: 'Fort Gordon',
};

jest.mock('hooks/queries', () => ({
  useOrdersDocumentQueries: () => {
    return {
      orders: {
        1: {
          agency: 'ARMY',
          customerID: '6ac40a00-e762-4f5f-b08d-3ea72a8e4b63',
          date_issued: '2018-03-15',
          department_indicator: 'AIR_FORCE',
          destinationDutyStation: mockDestinationDutyStation,
          eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MTE0Nlo=',
          entitlement: {
            authorizedWeight: 5000,
            dependentsAuthorized: true,
            eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42ODAwOVo=',
            id: '0dbc9029-dfc5-4368-bc6b-dfc95f5fe317',
            nonTemporaryStorage: true,
            privatelyOwnedVehicle: true,
            proGearWeight: 2000,
            proGearWeightSpouse: 500,
            storageInTransit: 2,
            totalDependents: 1,
            totalWeight: 5000,
          },
          first_name: 'Leo',
          grade: 'E_1',
          id: '1',
          last_name: 'Spacemen',
          order_number: 'ORDER3',
          order_type: 'PERMANENT_CHANGE_OF_STATION',
          order_type_detail: 'HHG_PERMITTED',
          originDutyStation: mockOriginDutyStation,
          report_by_date: '2018-08-01',
          tac: 'F8E1',
          sac: 'E2P3',
        },
      },
    };
  },
}));

describe('MoveAllowances page', () => {
  const wrapper = mount(
    <MockProviders initialEntries={['moves/1000/allowances']}>
      <MoveAllowances />
    </MockProviders>,
  );

  it('renders the sidebar elements', () => {
    expect(wrapper.find({ 'data-testid': 'allowances-header' }).text()).toBe('View Allowances');
    // There is only 1 button, but mount-rendering react-uswds Button component has inner buttons
    expect(wrapper.find({ 'data-testid': 'view-orders' }).at(0).text()).toBe('View Orders');
  });

  it('renders displays the allowances in the sidebar form', () => {
    expect(wrapper.find({ 'data-testid': 'weightAllowance' }).text()).toBe('5,000 lbs');
    expect(wrapper.find({ 'data-testid': 'proGearWeight' }).text()).toBe('2,000 lbs');
    expect(wrapper.find({ 'data-testid': 'spouseProGearWeight' }).text()).toBe('500 lbs');
    expect(wrapper.find({ 'data-testid': 'storageInTransit' }).text()).toBe('2 days');
  });
});
