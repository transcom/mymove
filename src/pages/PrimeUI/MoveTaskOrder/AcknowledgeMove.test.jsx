import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import AcknowledgeMove from './AcknowledgeMove';

import { renderWithProviders } from 'testUtils';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import { acknowledgeMovesAndShipments } from 'services/primeApi';

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

jest.mock('services/primeApi', () => ({
  acknowledgeMovesAndShipments: jest.fn(),
}));

const acknowledgedMoveReturnValue = {
  moveTaskOrder: {
    id: '1',
    moveCode: 'DEPPRQ',
    primeAcknowledgedAt: '2025-04-13T14:15:22.000Z',
  },
  isLoading: false,
  isError: false,
};

const unacknowledgedMoveReturnValue = {
  moveTaskOrder: {
    id: '2',
    moveCode: 'DEPPRZ',
    primeAcknowledgedAt: null,
  },
  isLoading: false,
  isError: false,
};

const primeAcknowledgedAtText = 'Prime Acknowledged At *';

describe('PrimeUI AcknowledgeMove Page', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  it('renders the form with the data from an acknowledged move and asterisks for required fields', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(acknowledgedMoveReturnValue);
    renderWithProviders(<AcknowledgeMove />);

    expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');

    // Verify fields are present and populated
    const moveCodeElement = screen.getByText('Move Code:');
    expect(moveCodeElement).toBeInTheDocument();
    expect(moveCodeElement.nextSibling).toHaveTextContent('DEPPRQ');

    const moveIdElement = screen.getByText('Move Id:');
    expect(moveIdElement).toBeInTheDocument();
    expect(moveIdElement.nextSibling).toHaveTextContent('1');

    const primeAcknowledgedAtLabel = screen.getByLabelText(primeAcknowledgedAtText);
    expect(primeAcknowledgedAtLabel).toBeInTheDocument();

    const dateInput = screen.getByLabelText(primeAcknowledgedAtText);
    expect(dateInput).toBeInTheDocument();
    expect(dateInput).toHaveValue('13 Apr 2025');

    // Prime Acknowledged Date field is disabled since a value is already present
    expect(dateInput).toBeDisabled();
  });

  it('renders the form with the data from an unacknowledged move ', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(unacknowledgedMoveReturnValue);
    renderWithProviders(<AcknowledgeMove />);

    const moveCodeElement = screen.getByText('Move Code:');
    expect(moveCodeElement).toBeInTheDocument();
    expect(moveCodeElement.nextSibling).toHaveTextContent('DEPPRZ');

    const moveIdElement = screen.getByText('Move Id:');
    expect(moveIdElement).toBeInTheDocument();
    expect(moveIdElement.nextSibling).toHaveTextContent('2');

    const primeAcknowledgedAtLabel = screen.getByLabelText(primeAcknowledgedAtText);
    expect(primeAcknowledgedAtLabel).toBeInTheDocument();

    // Verify prime acknowledged date is empty
    const dateInput = screen.getByLabelText(primeAcknowledgedAtText);
    expect(dateInput).toBeInTheDocument();
    expect(dateInput).not.toHaveValue();

    // Prime Acknowledged Date field is enabled since it is empty
    expect(dateInput).toBeEnabled();
  });

  it('calls acknowledgeMovesAndShipments when the form is submitted', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(unacknowledgedMoveReturnValue);
    renderWithProviders(<AcknowledgeMove />);

    // Input a prime acknowledged at date
    const dateInput = screen.getByLabelText(primeAcknowledgedAtText);
    userEvent.clear(dateInput);
    fireEvent.change(dateInput, { target: { value: '17 May 2024' } });

    // Submit the form
    const submitButton = screen.getByText('Save');
    userEvent.click(submitButton);

    // Verify that the mutation function was called with the correct parameters
    await waitFor(() => {
      expect(acknowledgeMovesAndShipments).toHaveBeenCalledWith({
        body: [
          {
            id: '2',
            primeAcknowledgedAt: '2024-05-17',
          },
        ],
      });
    });
  });

  it('enables the save button when the user inputs a valid date', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(unacknowledgedMoveReturnValue);
    renderWithProviders(<AcknowledgeMove />);

    const dateInput = screen.getByLabelText(primeAcknowledgedAtText);
    expect(dateInput).toBeInTheDocument();
    expect(dateInput).not.toHaveValue();

    const saveButton = screen.getByRole('button', { name: /Save/ });
    expect(saveButton).toBeInTheDocument();

    // Save button is initially disabled
    expect(saveButton).toBeDisabled();

    // Set a date value
    fireEvent.change(dateInput, { target: { value: '2025-04-01' } });

    await waitFor(() => {
      // Save button is now enabled since we inputted a valid date
      expect(saveButton).toBeEnabled();
    });
  });
});
