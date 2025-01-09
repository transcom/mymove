import React from 'react';
import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { usePrimeSimulatorGetMove } from '../../../hooks/queries';
import { updateShipmentDestinationAddress } from '../../../services/primeApi';

import PrimeUIShipmentUpdateDestinationAddress from './PrimeUIShipmentUpdateDestinationAddress';

import { ReactQueryWrapper, MockProviders } from 'testUtils';
import { primeSimulatorRoutes } from 'constants/routes';

const mockNavigate = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  ...jest.requireActual('services/primeApi'),
  updateShipmentDestinationAddress: jest.fn(),
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
      shipmentType: 'HHG_INTO_NTS',
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

const testShipmentReturnValue = {
  moveTaskOrder: {
    id: '1',
    moveCode: 'LN4T89',
    mtoShipments: [
      {
        id: '4',
        shipmentType: 'HHG',
        requestedPickupDate: '2021-12-01',
        pickupAddress: null,
        destinationAddress: {
          id: '2',
          streetAddress1: '100 1st Avenue',
          city: 'New York',
          state: 'NY',
          county: 'New York',
          postalCode: '10001',
        },
      },
    ],
  },
  isLoading: false,
  isError: false,
};

const renderComponent = () => {
  render(
    <ReactQueryWrapper>
      <MockProviders path={primeSimulatorRoutes.SHIPMENT_UPDATE_DESTINATION_ADDRESS_PATH} params={routingParams}>
        <PrimeUIShipmentUpdateDestinationAddress />
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
    it('displays the delivery address form', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(testShipmentReturnValue);

      renderComponent();

      const pageHeading = await screen.getByRole('heading', {
        name: 'Update Shipment Delivery Address',
        level: 2,
      });
      expect(pageHeading).toBeInTheDocument();

      const pageDescription = await screen.getByTestId('destination-form-details');
      expect(pageDescription).toBeInTheDocument();

      const shipmentIndex = testShipmentReturnValue.moveTaskOrder.mtoShipments.findIndex(
        (mtoShipment) => mtoShipment.id === routingParams.shipmentId,
      );
      const shipment = testShipmentReturnValue.moveTaskOrder.mtoShipments[shipmentIndex];

      await waitFor(() => {
        expect(screen.getAllByLabelText('Address 1').length).toBe(1);
        expect(screen.getAllByLabelText('Address 1')[0]).toHaveValue(shipment.destinationAddress.streetAddress1);
        expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
        expect(screen.getAllByLabelText(/Address 3/)[0]).toBeInTheDocument();

        expect(screen.getAllByText('City')[0]).toBeInTheDocument();
        expect(screen.getAllByText(shipment.destinationAddress.city)[0]).toBeInTheDocument();
        expect(screen.getAllByText('State')[0]).toBeInTheDocument();
        expect(screen.getAllByText(shipment.destinationAddress.state)[0]).toBeInTheDocument();
        expect(screen.getAllByText('County')[0]).toBeInTheDocument();
        expect(screen.getAllByText(shipment.destinationAddress.county)[0]).toBeInTheDocument();
        expect(screen.getAllByText('ZIP')[0]).toBeInTheDocument();
        expect(screen.getAllByText(shipment.destinationAddress.postalCode)[0]).toBeInTheDocument();

        expect(screen.getAllByLabelText('Contractor Remarks')[0]).toHaveValue('');
      });
    });
    it('displays validation error when contractor remarks are blank', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      renderComponent();

      await act(async () => {
        expect(screen.getAllByRole('button', { name: 'Save' }).length).toBe(1);
        await userEvent.click(screen.getAllByRole('button', { name: 'Save' })[0]);
      });

      expect(screen.getByText('Contractor remarks are required to make these changes')).toBeInTheDocument();
    });
  });

  describe('successful submission of form', () => {
    it('calls history router back to move details', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      updateShipmentDestinationAddress.mockReturnValue({
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

      await userEvent.type(screen.getByLabelText('Contractor Remarks'), 'Test remarks');

      await act(async () => {
        expect(screen.getAllByRole('button', { name: 'Save' }).length).toBe(1);
        await userEvent.click(screen.getAllByRole('button', { name: 'Save' })[0]);
      });

      await waitFor(() => {
        expect(mockNavigate).toHaveBeenCalledWith('/simulator/moves/LN4T89/details');
      });
    });
  });

  describe('error alert display', () => {
    it('displays the error alert when the api submission returns an error', () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      updateShipmentDestinationAddress.mockRejectedValue({
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
      updateShipmentDestinationAddress.mockRejectedValue('malformed api error response');

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
});
