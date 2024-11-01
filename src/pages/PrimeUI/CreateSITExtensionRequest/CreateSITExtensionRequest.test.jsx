import React from 'react';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import CreateSITExtensionRequest from './CreateSITExtensionRequest';

import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { MockProviders } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';
import { createSITExtensionRequest } from 'services/primeApi';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));
const routingParams = { moveCodeOrID: 'LN4T89', shipmentId: '2' };

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  createSITExtensionRequest: jest.fn(),
}));

const moveTaskOrder = {
  id: '1',
  moveCode: 'LN4T89',
  mtoShipments: [
    {
      id: '2',
      shipmentType: 'HHG',
      requestedPickupDate: '2021-11-26',
      pickupAddress: {
        streetAddress1: '100 1st Avenue',
        city: 'New York',
        state: 'NY',
        postalCode: '10001' },
        marketCode: 'd',
    },
  ],
  mtoServiceItems: [
    {
      reServiceCode: 'DDDSIT',
      reason: 'Holiday break',
      sitDestinationFinalAddress: {
        streetAddress1: '444 Main Ave',
        streetAddress2: 'Apartment 9000',
        streetAddress3: 'c/o Some Person',
        city: 'Anytown',
        state: 'AL',
        postalCode: '90210',
      },
      id: '45fe9475-d592-48f5-896a-45d4d6eb7e76',
    },
  ],
};

const moveReturnValue = {
  moveTaskOrder,
  isLoading: false,
  isError: false,
};

const renderWithProviders = () => {
  render(
    <MockProviders path={primeSimulatorRoutes.CREATE_SIT_EXTENSION_REQUEST_PATH} params={routingParams}>
      <CreateSITExtensionRequest setFlashMessage={jest.fn()} />
    </MockProviders>,
  );
};

describe('CreateSITExtensionRequest page', () => {
  describe('check loading and error component states', () => {
    const loadingReturnValue = {
      moveTaskOrder: undefined,
      isLoading: true,
      isError: false,
    };

    const errorReturnValue = {
      moveTaskOrder: undefined,
      isLoading: false,
      isError: true,
    };

    it('renders the loading placeholder when the query is still loading', () => {
      usePrimeSimulatorGetMove.mockReturnValue(loadingReturnValue);

      renderWithProviders();

      expect(screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 }));
    });

    it('renders the Something Went Wrong component when the query has an error', () => {
      usePrimeSimulatorGetMove.mockReturnValue(errorReturnValue);

      renderWithProviders();

      expect(screen.getByText(/Something went wrong./));
    });
  });

  describe('displaying header, shipment, and submit button', () => {
    it('displays the shipment information', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      renderWithProviders();

      const shipmentsHeading = screen.getByRole('heading', { name: 'Create SIT Extension Request', level: 1 });
      expect(shipmentsHeading).toBeInTheDocument();

      const shipmentsContainer = shipmentsHeading.parentElement;
      const hhgHeading = within(shipmentsContainer).getByRole('heading', {
        name: `${moveTaskOrder.mtoShipments[0].marketCode}HHG shipment`,
        level: 3,
      });

      expect(hhgHeading).toBeInTheDocument();
    });

    it('displays the submit button', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      renderWithProviders();

      expect(screen.getByRole('button', { name: 'Request SIT Extension' })).toBeEnabled();
    });
  });

  describe('successful submission of form', () => {
    it('calls history router back to move details', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      createSITExtensionRequest.mockReturnValue({});

      renderWithProviders();

      const requestReason = await screen.findByRole('combobox');
      await userEvent.selectOptions(requestReason, 'Other reason');
      const requestedDays = await screen.findByLabelText('Requested Days');
      await userEvent.clear(requestedDays);
      await userEvent.type(requestedDays, '13');

      const contractorRemarks = screen.getByLabelText('Contractor Remarks');
      await userEvent.type(contractorRemarks, 'testing contractor remarks');

      const submitBtn = await screen.getByRole('button', { name: 'Request SIT Extension' });

      await waitFor(() => {
        expect(submitBtn).toBeEnabled();
      });
      await userEvent.click(submitBtn);

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/simulator/moves/LN4T89/details');
      });
    });
  });
});
