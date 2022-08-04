import React from 'react';
import { useParams } from 'react-router-dom';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import PrimeUIShipmentUpdateReweigh from './PrimeUIShipmentUpdateReweigh';

import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { updatePrimeMTOShipmentReweigh } from 'services/primeApi';
import { MockProviders } from 'testUtils';

const mockUseHistoryPush = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: jest.fn().mockReturnValue({ moveCodeOrID: 'LN4T89', shipmentId: '4', reweighId: '1' }),
  useHistory: () => ({
    push: mockUseHistoryPush,
  }),
}));

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  updatePrimeMTOShipmentReweigh: jest.fn(),
}));

const moveTaskOrder = {
  id: '1',
  moveCode: 'LN4T89',
  mtoShipments: [
    {
      id: '4',
      shipmentType: 'HHG',
      reweigh: {
        id: '1',
        weight: 12,
        verificationReason: 'Reweigh performed.',
      },
    },
  ],
};

const moveReturnValue = {
  moveTaskOrder,
  isLoading: false,
  isError: false,
};

const noReweigh = {
  moveTaskOrder: {
    id: '1',
    moveCode: 'LN4T89',
    mtoShipments: [
      {
        id: '4',
        shipmentType: 'HHG',
      },
    ],
  },
  isLoading: false,
  isError: false,
};

const noVerificationReason = {
  moveTaskOrder: {
    id: '1',
    moveCode: 'LN4T89',
    mtoShipments: [
      {
        id: '4',
        shipmentType: 'HHG',
        reweigh: {
          id: '1',
          weight: 12,
        },
      },
    ],
  },
  isLoading: false,
  isError: false,
};

describe('PrimeUIShipmentUpdateReweigh page', () => {
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

      render(
        <MockProviders>
          <PrimeUIShipmentUpdateReweigh />
        </MockProviders>,
      );

      expect(await screen.findByRole('heading', { name: 'Loading, please wait...', level: 2 }));
    });

    it('renders the Something Went Wrong component when the query has an error', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(errorReturnValue);

      render(
        <MockProviders>
          <PrimeUIShipmentUpdateReweigh />
        </MockProviders>,
      );

      expect(await screen.findByText(/Something went wrong./));
    });
  });

  describe('displaying shipment reweigh information', () => {
    it('displays the reweigh weight and verification reason', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);

      render(
        <MockProviders>
          <PrimeUIShipmentUpdateReweigh />
        </MockProviders>,
      );

      const pageHeading = await screen.findByRole('heading', {
        name: 'Edit Reweigh',
        level: 1,
      });
      expect(pageHeading).toBeInTheDocument();

      const { shipmentId } = useParams();
      const shipmentIndex = moveTaskOrder.mtoShipments.findIndex((mtoShipment) => mtoShipment.id === shipmentId);
      const shipment = moveTaskOrder.mtoShipments[shipmentIndex];

      expect(await screen.findByLabelText('Reweigh Weight (lbs)')).toHaveValue(String(shipment.reweigh.weight));

      expect(screen.getByTestId('remarks')).toHaveValue(shipment.reweigh.verificationReason);
    });

    it('displays only the reweigh weight', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(noVerificationReason);

      render(
        <MockProviders>
          <PrimeUIShipmentUpdateReweigh />
        </MockProviders>,
      );

      const pageHeading = await screen.findByRole('heading', {
        name: 'Edit Reweigh',
        level: 1,
      });
      expect(pageHeading).toBeInTheDocument();

      const { shipmentId } = useParams();
      const shipmentIndex = noVerificationReason.moveTaskOrder.mtoShipments.findIndex(
        (mtoShipment) => mtoShipment.id === shipmentId,
      );
      const shipment = noVerificationReason.moveTaskOrder.mtoShipments[shipmentIndex];

      expect(await screen.findByLabelText('Reweigh Weight (lbs)')).toHaveValue(String(shipment.reweigh.weight));
      expect(screen.getByTestId('remarks')).toHaveValue('');
    });

    it('uses the default values when there is no reweigh', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(noReweigh);

      render(
        <MockProviders>
          <PrimeUIShipmentUpdateReweigh />
        </MockProviders>,
      );

      const pageHeading = await screen.findByRole('heading', {
        name: 'Edit Reweigh',
        level: 1,
      });
      expect(pageHeading).toBeInTheDocument();

      expect(screen.getByLabelText('Reweigh Weight (lbs)')).toHaveValue('0');
      expect(screen.getByTestId('remarks')).toHaveValue('');
    });
  });

  describe('error alert display', () => {
    it('displays the error alert when the api submission returns an error', async () => {
      jest.spyOn(console, 'error').mockImplementation(() => {});
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      updatePrimeMTOShipmentReweigh.mockRejectedValue({
        response: { body: { title: 'Error title', detail: 'Error detail' } },
      });

      render(
        <MockProviders>
          <PrimeUIShipmentUpdateReweigh />
        </MockProviders>,
      );

      const saveButton = screen.getByRole('button', { name: 'Save' });
      await userEvent.click(saveButton);

      expect(await screen.findByText(/Error title/)).toBeInTheDocument();
      expect(screen.getByText('Error detail')).toBeInTheDocument();
    });

    it('displays the unknown error when none is provided', async () => {
      jest.spyOn(console, 'error').mockImplementation(() => {});
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      updatePrimeMTOShipmentReweigh.mockRejectedValue('malformed api error response');

      render(
        <MockProviders>
          <PrimeUIShipmentUpdateReweigh />
        </MockProviders>,
      );

      const saveButton = screen.getByRole('button', { name: 'Save' });
      await userEvent.click(saveButton);

      expect(await screen.findByText('Unexpected error')).toBeInTheDocument();
      expect(
        screen.getByText('An unknown error has occurred, please check the address values used'),
      ).toBeInTheDocument();
    });
  });

  describe('successful submission of form', () => {
    it('calls history router back to move details', async () => {
      usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
      updatePrimeMTOShipmentReweigh.mockReturnValue({
        id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
        weight: 123,
        verificationReason: 'Reweigh performed',
        eTag: '1234567890',
      });

      render(
        <MockProviders>
          <PrimeUIShipmentUpdateReweigh />
        </MockProviders>,
      );

      const saveButton = screen.getByRole('button', { name: 'Save' });
      await userEvent.click(saveButton);

      await waitFor(() => {
        expect(mockUseHistoryPush).toHaveBeenCalledWith('/simulator/moves/LN4T89/details');
      });
    });
  });
});
