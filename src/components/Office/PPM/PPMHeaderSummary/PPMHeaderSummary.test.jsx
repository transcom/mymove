import React from 'react';
import { waitFor, screen, fireEvent } from '@testing-library/react';

import PPMHeaderSummary from './PPMHeaderSummary';

import { useEditShipmentQueries, usePPMShipmentDocsQueries } from 'hooks/queries';
import { renderWithProviders } from 'testUtils';
import { tooRoutes } from 'constants/routes';
import { getPPMTypeLabel, PPM_TYPES } from 'shared/constants';
import { formatWeight } from 'utils/formatters';

beforeEach(() => {
  jest.clearAllMocks();
});

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
const routingParams = { moveCode: 'move123', shipmentId: 'shipment123' };
const mockRoutingConfig = {
  path: tooRoutes.BASE_SHIPMENT_EDIT_PATH,
  params: routingParams,
};

jest.mock('hooks/queries', () => ({
  usePPMShipmentDocsQueries: jest.fn(),
  useEditShipmentQueries: jest.fn(),
}));

const useEditShipmentQueriesReturnValue = {
  move: {
    id: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    ordersId: '1',
    status: 'NEEDS SERVICE COUNSELING',
  },
  order: {
    id: '1',
    originDutyLocation: {
      address: {
        streetAddress1: '',
        city: 'Fort Knox',
        state: 'KY',
        postalCode: '40121',
      },
    },
    destinationDutyLocation: {
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
    order_type: 'PERMANENT_CHANGE_OF_STATION',
    order_type_detail: 'HHG_PERMITTED',
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
      id: 'shipment123',
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
      requestedDeliveryDate: '2018-04-15',
      scheduledDeliveryDate: '2014-04-16',
      shipmentType: 'HHG',
      status: 'SUBMITTED',
      updatedAt: '2020-06-10T15:58:02.404031Z',
    },
  ],
  isLoading: false,
  isError: false,
  isSuccess: true,
};

const defaultProps = {
  ppmShipmentInfo: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    expectedDepartureDate: '2022-12-02',
    actualMoveDate: '2022-12-06',
    pickupAddress: {
      streetAddress1: '812 S 129th St',
      streetAddress2: '#123',
      city: 'San Antonio',
      state: 'TX',
      postalCode: '78234',
    },
    destinationAddress: {
      streetAddress1: '456 Oak Ln.',
      streetAddress2: '#123',
      city: 'Oakland',
      state: 'CA',
      postalCode: '94611',
    },
    miles: 300,
    estimatedWeight: 3000,
    actualWeight: 3500,
    isActualExpenseReimbursement: true,
  },
  ppmNumber: '1',
  showAllFields: false,
  readOnly: false,
};

const smallPackagePPMInfo = {
  id: 'eyedee',
  ppmType: PPM_TYPES.SMALL_PACKAGE,
  expectedDepartureDate: '2022-12-02',
  actualMoveDate: '2022-12-06',
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '456 Oak Ln.',
    city: 'Oakland',
    state: 'CA',
    postalCode: '94611',
  },
  miles: 300,
  estimatedWeight: 3000,
  actualWeight: 3500,
  allowableWeight: 4000,
  isActualExpenseReimbursement: false,
  movingExpenses: [
    { weightShipped: 2000, isProGear: false },
    { weightShipped: 500, isProGear: true, proGearBelongsToSelf: false },
  ],
};

const smallPackageProps = {
  ppmShipmentInfo: smallPackagePPMInfo,
  order: { grade: 'ARMY' },
  ppmNumber: '1',
  showAllFields: false,
  readOnly: false,
};

describe('PPMHeaderSummary component', () => {
  describe('displays form', () => {
    it('renders default values', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      renderWithProviders(<PPMHeaderSummary {...defaultProps} />, mockRoutingConfig);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 3, name: 'PPM 1' })).toBeInTheDocument();
      });
      expect(screen.getByTestId('tag', { name: 'actual expense reimbursement' })).toBeInTheDocument();

      expect(screen.getByText('Expense Type')).toBeInTheDocument();
      expect(screen.getByText('Planned Move Start Date')).toBeInTheDocument();
      expect(screen.getByText('02-Dec-2022')).toBeInTheDocument();
      expect(screen.getByText('Actual Move Start Date')).toBeInTheDocument();
      expect(screen.getByText('06-Dec-2022')).toBeInTheDocument();
      expect(screen.getByText('Starting Address')).toBeInTheDocument();
      expect(screen.getByText('812 S 129th St, #123, San Antonio, TX 78234')).toBeInTheDocument();
      expect(screen.getByText('Ending Address')).toBeInTheDocument();
      expect(screen.getByText('456 Oak Ln., #123, Oakland, CA 94611')).toBeInTheDocument();
      expect(screen.getByText('Miles')).toBeInTheDocument();
      expect(screen.getByText('300')).toBeInTheDocument();
      expect(screen.getByText('Estimated Net Weight')).toBeInTheDocument();
      expect(screen.getByText('3,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Actual Net Weight')).toBeInTheDocument();
      expect(screen.getByText('3,500 lbs')).toBeInTheDocument();

      fireEvent.click(screen.getByTestId('shipmentInfo-showRequestDetailsButton'));
      await waitFor(() => {
        expect(screen.getByText('Show Details', { exact: false })).toBeInTheDocument();
      });
    });

    it('renders small package tag and small package details', async () => {
      usePPMShipmentDocsQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      renderWithProviders(<PPMHeaderSummary {...smallPackageProps} />, mockRoutingConfig);

      // verify that the header shows the small package tag.
      expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent('PPM 1');
      const smallPackageTag = screen.getByTestId('smallPackageTag');
      expect(smallPackageTag).toBeInTheDocument();
      expect(smallPackageTag).toHaveTextContent(getPPMTypeLabel(PPM_TYPES.SMALL_PACKAGE));

      expect(screen.getByText('Allowable Weight')).toBeInTheDocument();
      expect(screen.getByText(formatWeight(2500))).toBeInTheDocument();
      expect(screen.getByText('Pro-gear')).toBeInTheDocument();
      // should be two because we have pro gear (one yes) and one spouse pro-gear (yes)
      expect(screen.getAllByText('Yes')).toHaveLength(2);
      expect(screen.getByText('Spouse Pro-gear')).toBeInTheDocument();
    });
  });
});
