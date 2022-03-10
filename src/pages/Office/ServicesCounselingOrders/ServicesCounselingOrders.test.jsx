/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ServicesCounselingOrders from 'pages/Office/ServicesCounselingOrders/ServicesCounselingOrders';
import { MockProviders } from 'testUtils';
import { useOrdersDocumentQueries } from 'hooks/queries';

const mockOriginDutyLocation = {
  address: {
    city: 'Des Moines',
    country: 'US',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC42OTg1OTha',
    id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
    postalCode: '50309',
    state: 'IA',
    streetAddress1: '987 Other Avenue',
    streetAddress2: 'P.O. Box 1234',
    streetAddress3: 'c/o Another Person',
  },
  address_id: '2e26b066-aaca-4563-b284-d7f3f978fb3c',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MTozOC43MDcxOTVa',
  id: 'a3ec2bdd-aa0a-434a-ba58-34c85f047704',
  name: 'XBc1KNi3pA',
};

const mockDestinationDutyLocation = {
  address: {
    city: 'Augusta',
    country: 'United States',
    eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
    id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    postalCode: '30813',
    state: 'GA',
    streetAddress1: 'Fort Gordon',
  },
  address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
  eTag: 'MjAyMC0wOS0xNFQxNzo0MDo0OC44OTM3MDVa',
  id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
  name: 'Fort Gordon',
};

jest.mock('hooks/queries', () => ({
  useOrdersDocumentQueries: jest.fn(),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  getTacValid: ({ tac }) => {
    return {
      tac,
      isValid: tac === '1111' || tac === '2222',
    };
  },
}));

const useOrdersDocumentQueriesReturnValue = {
  orders: {
    1: {
      agency: 'ARMY',
      customerID: '6ac40a00-e762-4f5f-b08d-3ea72a8e4b63',
      date_issued: '2018-03-15',
      department_indicator: 'AIR_FORCE',
      destinationDutyLocation: mockDestinationDutyLocation,
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
      originDutyLocation: mockOriginDutyLocation,
      report_by_date: '2018-08-01',
      tac: 'F8E1',
      sac: 'E2P3',
      ntsTac: '1111',
      ntsSac: 'R6X1',
    },
  },
};

const loadingReturnValue = {
  ...useOrdersDocumentQueriesReturnValue,
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  ...useOrdersDocumentQueriesReturnValue,
  isLoading: false,
  isError: true,
  isSuccess: false,
};

describe('Orders page', () => {
  describe('check loading and error component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useOrdersDocumentQueries.mockReturnValue(loadingReturnValue);

      render(
        <MockProviders initialEntries={['moves/FP24I2/orders']}>
          <ServicesCounselingOrders />
        </MockProviders>,
      );

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useOrdersDocumentQueries.mockReturnValue(errorReturnValue);

      render(
        <MockProviders initialEntries={['moves/FP24I2/orders']}>
          <ServicesCounselingOrders />
        </MockProviders>,
      );

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Basic rendering', () => {
    it('renders the sidebar orders detail form', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders initialEntries={['moves/FP24I2/orders']}>
          <ServicesCounselingOrders />
        </MockProviders>,
      );

      expect(await screen.findByLabelText('Current duty location')).toBeInTheDocument();
    });

    it('renders the sidebar elements', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders initialEntries={['moves/FP24I2/orders']}>
          <ServicesCounselingOrders />
        </MockProviders>,
      );

      expect(await screen.findByTestId('view-orders-header')).toHaveTextContent('View orders');
      expect(screen.getByTestId('view-allowances')).toHaveTextContent('View allowances');
    });

    it('populates initial field values', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders initialEntries={['moves/FP24I2/orders']}>
          <ServicesCounselingOrders />
        </MockProviders>,
      );

      expect(await screen.findByText(mockOriginDutyLocation.name)).toBeInTheDocument();
      expect(screen.getByText(mockDestinationDutyLocation.name)).toBeInTheDocument();
      expect(screen.getByLabelText('Orders type')).toHaveValue('PERMANENT_CHANGE_OF_STATION');
      expect(screen.getByTestId('hhgTacInput')).toHaveValue('F8E1');
      expect(screen.getByTestId('hhgSacInput')).toHaveValue('E2P3');
      expect(screen.getByTestId('ntsTacInput')).toHaveValue('1111');
      expect(screen.getByTestId('ntsSacInput')).toHaveValue('R6X1');
    });
  });

  describe('TAC validation', () => {
    it('validates on load', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders initialEntries={['moves/FP24I2/orders']}>
          <ServicesCounselingOrders />
        </MockProviders>,
      );

      expect(await screen.findByText(/This TAC does not appear in TGET/)).toBeInTheDocument();
    });

    it('validates on user input', async () => {
      useOrdersDocumentQueries.mockReturnValue(useOrdersDocumentQueriesReturnValue);

      render(
        <MockProviders initialEntries={['moves/FP24I2/orders']}>
          <ServicesCounselingOrders />
        </MockProviders>,
      );

      const hhgTacInput = screen.getByTestId('hhgTacInput');
      userEvent.clear(hhgTacInput);
      userEvent.type(hhgTacInput, '2222');

      await waitFor(() => {
        expect(screen.queryByText(/This TAC does not appear in TGET/)).not.toBeInTheDocument();
      });

      userEvent.clear(hhgTacInput);
      userEvent.type(hhgTacInput, '3333');

      await waitFor(() => {
        expect(screen.getByText(/This TAC does not appear in TGET/)).toBeInTheDocument();
      });
    });
  });
});
