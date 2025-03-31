import React from 'react';
import { screen } from '@testing-library/react';
import AcknowledgeMove from './AcknowledgeMove';
import { renderWithProviders } from 'testUtils';
import { usePrimeSimulatorGetMove } from 'hooks/queries';

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
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

describe('PrimeUI AcknowledgeMove', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  it('renders the form with the data from an acknowledged move', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(acknowledgedMoveReturnValue);
    renderWithProviders(<AcknowledgeMove />);

    const moveCodeElement = screen.getByText('Move Code:');
    expect(moveCodeElement).toBeInTheDocument();
    expect(moveCodeElement.nextSibling).toHaveTextContent('DEPPRQ');

    const moveIdElement = screen.getByText('Move Id:');
    expect(moveIdElement).toBeInTheDocument();
    expect(moveIdElement.nextSibling).toHaveTextContent('1');

    const primeAcknowledgedAtText = 'Prime Acknowledged At';
    const primeAcknowledgedAtLabel = screen.getByText(primeAcknowledgedAtText);
    expect(primeAcknowledgedAtLabel).toBeInTheDocument();
    const dateInput = screen.getByLabelText(primeAcknowledgedAtText);
    expect(dateInput).toBeInTheDocument();
    expect(dateInput).toHaveValue('13 Apr 2025');
  });

  it('renders the form with the data from an unacknowledged move', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(unacknowledgedMoveReturnValue);
    renderWithProviders(<AcknowledgeMove />);

    const moveCodeElement = screen.getByText('Move Code:');
    expect(moveCodeElement).toBeInTheDocument();
    expect(moveCodeElement.nextSibling).toHaveTextContent('DEPPRZ');

    const moveIdElement = screen.getByText('Move Id:');
    expect(moveIdElement).toBeInTheDocument();
    expect(moveIdElement.nextSibling).toHaveTextContent('2');

    const primeAcknowledgedAtText = 'Prime Acknowledged At';
    const primeAcknowledgedAtLabel = screen.getByText(primeAcknowledgedAtText);
    expect(primeAcknowledgedAtLabel).toBeInTheDocument();
    const dateInput = screen.getByLabelText(primeAcknowledgedAtText);
    expect(dateInput).toBeInTheDocument();
    expect(dateInput).not.toHaveValue();
  });
});
