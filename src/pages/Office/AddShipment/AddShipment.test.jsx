/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import AddShipment from './AddShipment';

import { createMTOShipment } from 'services/ghcApi';
import { useEditShipmentQueries } from 'hooks/queries';
import { MockProviders } from 'testUtils';
import { tooRoutes } from 'constants/routes';

// Explicitly setup navigate mock so we can verify it was called with correct pathing in tests
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  createMTOShipment: jest.fn(),
}));

jest.mock('hooks/queries', () => ({
  ...jest.requireActual('@tanstack/react-query'),
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

const loadingReturnValue = {
  ...useEditShipmentQueriesReturnValue,
  isLoading: true,
  isError: false,
  isSuccess: false,
};

const errorReturnValue = {
  ...useEditShipmentQueriesReturnValue,
  isLoading: false,
  isError: true,
  isSuccess: false,
};

const renderWithMocks = () => {
  render(
    <MockProviders path={tooRoutes.BASE_SHIPMENT_ADD_PATH} params={{ moveCode: 'move123', shipmentType: 'HHG' }}>
      <AddShipment />
    </MockProviders>,
  );
};

beforeEach(() => {
  jest.resetAllMocks();
});

describe('AddShipment component', () => {
  describe('check different component states', () => {
    it('renders the Loading Placeholder when the query is still loading', async () => {
      useEditShipmentQueries.mockReturnValue(loadingReturnValue);
      renderWithMocks();

      const h2 = await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 });
      expect(h2).toBeInTheDocument();
    });

    it('renders the Something Went Wrong component when the query errors', async () => {
      useEditShipmentQueries.mockReturnValue(errorReturnValue);

      renderWithMocks();

      const errorMessage = await screen.getByText(/Something went wrong./);
      expect(errorMessage).toBeInTheDocument();
    });
  });

  describe('Basic rendering', () => {
    it('renders the Shipment Form', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      renderWithMocks();

      const h1 = await screen.getByRole('heading', { name: 'Add shipment details', level: 1 });
      await waitFor(() => {
        expect(h1).toBeInTheDocument();
      });
    });

    it('routes to the move details page when the save button is clicked', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      createMTOShipment.mockImplementation(() => Promise.resolve({}));

      renderWithMocks();

      const saveButton = screen.getByRole('button', { name: 'Save' });

      expect(saveButton).toBeInTheDocument();

      await waitFor(() => {
        expect(saveButton).toBeDisabled();
      });

      expect(screen.getByLabelText('Use current address')).not.toBeChecked();

      await userEvent.type(screen.getAllByLabelText('Address 1')[0], '812 S 129th St');
      await userEvent.type(screen.getAllByLabelText('City')[0], 'San Antonio');
      await userEvent.selectOptions(screen.getAllByLabelText('State')[0], ['TX']);
      await userEvent.type(screen.getAllByLabelText('ZIP')[0], '78234');
      await userEvent.type(screen.getByLabelText('Requested pickup date'), '01 Nov 2020');
      await userEvent.type(screen.getByLabelText('Requested delivery date'), '08 Nov 2020');

      await waitFor(() => {
        expect(saveButton).not.toBeDisabled();
      });

      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/moves/move123/details');
      });
    });

    it('routes to the move details page when the cancel button is clicked', async () => {
      useEditShipmentQueries.mockReturnValue(useEditShipmentQueriesReturnValue);
      renderWithMocks();

      const cancelButton = screen.getByRole('button', { name: 'Cancel' });

      expect(cancelButton).not.toBeDisabled();

      await userEvent.click(cancelButton);

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/moves/move123/details');
      });
    });
  });
});
