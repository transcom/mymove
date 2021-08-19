/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { mount } from 'enzyme';
import { render, screen } from '@testing-library/react';

import ServicesCounselingMoveAllowances from 'pages/Office/ServicesCounselingMoveAllowances/ServicesCounselingMoveAllowances';
import { MockProviders } from 'testUtils';
import { useOrdersDocumentQueries } from 'hooks/queries';

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
  useOrdersDocumentQueries: jest.fn(),
}));

const useOrdersDocumentQueriesReturnValue = {
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
        requiredMedicalEquipmentWeight: 1000,
        organizationalClothingAndIndividualEquipment: true,
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

describe('MoveAllowances page', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useOrdersDocumentQueries.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders initialEntries={['/counseling/moves/1000/allowances']}>
          <ServicesCounselingMoveAllowances />
        </MockProviders>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useOrdersDocumentQueries.mockReturnValue(errorReturnValue);

      render(
        <MockProviders initialEntries={['/counseling/moves/1000/allowances']}>
          <ServicesCounselingMoveAllowances />
        </MockProviders>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Basic rendering', () => {
    useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

    const wrapper = mount(
      <MockProviders initialEntries={['/counseling/moves/1000/allowances']}>
        <ServicesCounselingMoveAllowances />
      </MockProviders>,
    );

    it('renders the sidebar elements', () => {
      expect(wrapper.find({ 'data-testid': 'allowances-header' }).text()).toBe('View allowances');
      expect(wrapper.find({ 'data-testid': 'view-orders' }).at(0).text()).toBe('View orders');
      expect(wrapper.find({ 'data-testid': 'header' }).text()).toBe('Counseling');
    });

    it('renders displays the allowances in the sidebar form', () => {
      // Pro-gear
      expect(wrapper.find(`input[data-testid="proGearWeightInput"]`).getDOMNode().value).toBe('2,000');

      // Pro-gear spouse
      expect(wrapper.find(`input[data-testid="proGearWeightSpouseInput"]`).getDOMNode().value).toBe('500');

      // RME
      expect(wrapper.find(`input[data-testid="rmeInput"]`).getDOMNode().value).toBe('1,000');

      // Branch
      expect(wrapper.find(`select[data-testid="branchInput"]`).getDOMNode().value).toBe('ARMY');

      // Rank
      expect(wrapper.find(`select[data-testid="rankInput"]`).getDOMNode().value).toBe('E_1');

      // OCIE
      expect(
        wrapper.find(`input[name="organizationalClothingAndIndividualEquipment"]`).getDOMNode().checked,
      ).toBeTruthy();

      // Authorized weight
      expect(wrapper.find('dd').at(0).text()).toBe('5,000 lbs');

      // Weight allowance
      expect(wrapper.find('dd').at(1).text()).toBe('5,000 lbs');

      // Storage in-transit
      expect(wrapper.find('dd').at(2).text()).toBe('2 days');

      // Dependents authorized
      expect(wrapper.find(`input[name="dependentsAuthorized"]`).getDOMNode().checked).toBeTruthy();
    });
  });
});
