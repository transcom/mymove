import React from 'react';
import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { usePrimeSimulatorGetMove } from '../../../hooks/queries';
import { updatePrimeMTOShipmentAddress } from '../../../services/primeApi';

import PrimeUIShipmentUpdateAddress from './PrimeUIShipmentUpdateAddress';

import { ReactQueryWrapper, MockProviders } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';

const mockNavigate = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  useLocation: () => ({ state: { addressType: 'pickupAddress' } }),
}));

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  updatePrimeMTOShipmentAddress: jest.fn(),
}));

const routingParams = { moveCodeOrID: 'LN4T89', shipmentId: '4' };

const moveTaskOrder = {
  id: '1',
  moveCode: 'LN4T89',
  mtoShipments: [
    {
      id: '2',
      shipmentType: 'HHG',
      requestedPickupDate: '2021-11-26',
      pickupAddress: { streetAddress1: '100 1st Avenue', city: 'New York', state: 'NY', postalCode: '10001' },
    },
    {
      id: '3',
      shipmentType: 'HHG_INTO_NTS_DOMESTIC',
      requestedPickupDate: '2021-12-01',
      pickupAddress: { streetAddress1: '800 Madison Avenue', city: 'New York', state: 'NY', postalCode: '10002' },
    },
    {
      id: '4',
      shipmentType: 'HHG',
      requestedPickupDate: '2021-12-01',
      pickupAddress: {
        id: '1',
        streetAddress1: '800 Madison Avenue',
        city: 'New York',
        state: 'NY',
        postalCode: '10002',
      },
      destinationAddress: {
        id: '2',
        streetAddress1: '100 1st Avenue',
        city: 'New York',
        state: 'NY',
        postalCode: '10001',
      },
    },
  ],
};

const moveReturnValue = {
  moveTaskOrder,
  isLoading: false,
  isError: false,
};

const renderComponent = () => {
  render(
    <ReactQueryWrapper>
      <MockProviders path={primeSimulatorRoutes.SHIPMENT_UPDATE_ADDRESS_PATH} params={routingParams}>
        <PrimeUIShipmentUpdateAddress />
      </MockProviders>
    </ReactQueryWrapper>,
  );
};

describe('PrimeUIShipmentUpdateAddress page', () => {
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

    it('renders the loading placeholder when the query is still loading', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(loadingReturnValue);

      renderComponent();

      expect(await screen.getByRole('heading', { name: 'Loading, please wait...', level: 2 }));
    });

    it('renders the Something Went Wrong component when the query has an error', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(errorReturnValue);

      renderComponent();

      expect(await screen.getByText(/Something went wrong./));
    });
  });

  describe('displaying shipment address information', () => {
    it('displays shipment pickup and destination address', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      renderComponent();

      const pageHeading = await screen.getByRole('heading', {
        name: 'Update Existing Pickup Address',
        level: 1,
      });
      expect(pageHeading).toBeInTheDocument();

      const shipmentIndex = moveTaskOrder.mtoShipments.findIndex(
        (mtoShipment) => mtoShipment.id === routingParams.shipmentId,
      );
      const shipment = moveTaskOrder.mtoShipments[shipmentIndex];

      await waitFor(() => {
        expect(screen.getAllByLabelText('Address 1')[0]).toHaveValue(shipment.pickupAddress.streetAddress1);
        expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
        expect(screen.getAllByLabelText('City')[0]).toHaveValue(shipment.pickupAddress.city);
        expect(screen.getAllByLabelText('State')[0]).toHaveValue(shipment.pickupAddress.state);
        expect(screen.getAllByLabelText('ZIP')[0]).toHaveValue(shipment.pickupAddress.postalCode);
      });
    });
  });

  describe('error alert display', () => {
    it('displays the error alert when the api submission returns an error', () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      updatePrimeMTOShipmentAddress.mockRejectedValue({
        response: { body: { title: 'Error title', detail: 'Error detail' } },
      });

      renderComponent();

      waitFor(async () => {
        expect(screen.getAllByRole('button', { name: 'Save' }).length).toBe(2);
        await userEvent.click(screen.getAllByRole('button', { name: 'Save' })[0]);
        expect(screen.getByText('Error title')).toBeInTheDocument();
        expect(screen.getByText('Error detail')).toBeInTheDocument();
      });
    });

    it('displays the unknown error when none is provided', () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      updatePrimeMTOShipmentAddress.mockRejectedValue('malformed api error response');

      renderComponent();

      waitFor(async () => {
        expect(screen.getAllByRole('button', { name: 'Save' }).length).toBe(2);
        await userEvent.click(screen.getAllByRole('button', { name: 'Save' })[0]);

        expect(screen.getByText('Unexpected error')).toBeInTheDocument();
        expect(
          screen.getByText('An unknown error has occurred, please check the address values used'),
        ).toBeInTheDocument();
      });
    });
  });

  describe('successful submission of form', () => {
    it('calls history router back to move details', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      updatePrimeMTOShipmentAddress.mockReturnValue({
        id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
        streetAddress1: '444 Main Ave',
        streetAddress2: 'Apartment 9000',
        streetAddress3: '',
        city: 'Anytown',
        state: 'AL',
        postalCode: '90210',
        country: 'USA',
        eTag: '1234567890',
      });

      renderComponent();

      await act(async () => {
        expect(screen.getAllByRole('button', { name: 'Save' }).length).toBe(1);
        await userEvent.click(screen.getAllByRole('button', { name: 'Save' })[0]);
      });

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/simulator/moves/LN4T89/details');
      });
    });
  });
});
