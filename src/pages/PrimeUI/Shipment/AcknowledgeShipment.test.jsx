import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import AcknowledgeShipment from './AcknowledgeShipment';

import { renderWithProviders } from 'testUtils';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { primeSimulatorRoutes } from 'constants/routes';
import { acknowledgeMovesAndShipments } from 'services/primeApi';

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  acknowledgeMovesAndShipments: jest.fn(),
}));

const moveReturnValue = {
  moveTaskOrder: {
    id: '1',
    mtoShipments: [
      {
        id: '2',
        primeAcknowledgedAt: '2025-04-13T14:15:22.000Z',
      },
      {
        id: '3',
      },
    ],
  },
  isLoading: false,
  isError: false,
};

const primeAcknowledgedAtText = 'Prime Acknowledged At *';

const renderShipmentComponent = (shipmentId) => {
  renderWithProviders(<AcknowledgeShipment />, {
    path: primeSimulatorRoutes.ACKNOWLEDGE_SHIPMENT_PATH,
    params: { moveCodeOrID: '1', shipmentId },
  });
};

describe('PrimeUI Acknowledge Shipment Page', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  it('renders the form with the data from an acknowledged shipment', async () => {
    renderShipmentComponent('2');

    const heading = screen.getByText('Acknowledge Shipment');
    expect(heading).toBeInTheDocument();

    const shipmentIdElement = screen.getByText('Shipment ID:');
    expect(shipmentIdElement).toBeInTheDocument();
    expect(shipmentIdElement.nextSibling).toHaveTextContent('2');

    const primeAcknowledgedAtLabel = screen.getByLabelText(primeAcknowledgedAtText);
    expect(primeAcknowledgedAtLabel).toBeInTheDocument();

    const dateInput = screen.getByLabelText(primeAcknowledgedAtText);
    expect(dateInput).toBeInTheDocument();
    expect(dateInput).toHaveValue('13 Apr 2025');
  });

  it('renders the form with the data from an unacknowledged shipment', async () => {
    renderShipmentComponent('3');

    const heading = screen.getByText('Acknowledge Shipment');
    expect(heading).toBeInTheDocument();

    const shipmentIdElement = screen.getByText('Shipment ID:');
    expect(shipmentIdElement).toBeInTheDocument();
    expect(shipmentIdElement.nextSibling).toHaveTextContent('3');

    const primeAcknowledgedAtLabel = screen.getByLabelText(primeAcknowledgedAtText);
    expect(primeAcknowledgedAtLabel).toBeInTheDocument();

    const dateInput = screen.getByLabelText(primeAcknowledgedAtText);
    expect(dateInput).toBeInTheDocument();
    expect(dateInput).not.toHaveValue();
  });

  it('calls acknowledgeMovesAndShipments when the form is submitted', async () => {
    renderShipmentComponent('3');

    // Input a prime acknowledged at date
    const dateInput = screen.getByLabelText(primeAcknowledgedAtText);
    userEvent.clear(dateInput);
    fireEvent.change(dateInput, { target: { value: '21 Nov 2024' } });

    // Submit the form
    const saveButton = screen.getByRole('button', { name: /Save/ });
    expect(saveButton).toBeInTheDocument();
    userEvent.click(saveButton);

    // Verify that the mutation function was called with the correct parameters
    await waitFor(() => {
      expect(acknowledgeMovesAndShipments).toHaveBeenCalledWith({
        body: [
          {
            id: '1',
            mtoShipments: [
              {
                id: '3',
                primeAcknowledgedAt: '2024-11-21',
              },
            ],
          },
        ],
      });
    });
  });

  it('disables the save button when an invalid date is inputted', async () => {
    renderShipmentComponent('3');

    const saveButton = screen.getByRole('button', { name: /Save/ });
    expect(saveButton).toBeInTheDocument();

    // Input an invalid prime acknowledged at date
    const dateInput = screen.getByLabelText(primeAcknowledgedAtText);
    userEvent.clear(dateInput);
    fireEvent.change(dateInput, { target: { value: '99 Nov 2024' } });

    await waitFor(() => {
      // Save button is disabled since we inputted an invalid date
      expect(saveButton).toBeDisabled();
    });
  });

  it('enables the save button when the user inputs a valid date', async () => {
    renderShipmentComponent('3');

    const dateInput = screen.getByLabelText(primeAcknowledgedAtText);
    expect(dateInput).toBeInTheDocument();
    expect(dateInput).not.toHaveValue();

    const saveButton = screen.getByRole('button', { name: /Save/ });
    expect(saveButton).toBeInTheDocument();

    // Save button is initially disabled
    expect(saveButton).toBeDisabled();

    // Set a date value
    fireEvent.change(dateInput, { target: { value: '2025-02-10' } });

    await waitFor(() => {
      // Save button is now enabled since we inputted a valid date
      expect(saveButton).toBeEnabled();
    });
  });
});
