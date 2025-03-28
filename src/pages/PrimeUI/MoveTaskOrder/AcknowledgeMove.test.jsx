import React from 'react';
import { screen } from '@testing-library/react';
import AcknowledgeMove from './AcknowledgeMove';
import { renderWithProviders } from 'testUtils';
import { usePrimeSimulatorGetMove } from 'hooks/queries';

jest.mock('hooks/queries', () => ({
  usePrimeSimulatorGetMove: jest.fn(),
}));

const moveTaskOrder = {
  id: '1',
  moveCode: 'DEPPRQ',
};

const moveReturnValue = {
  moveTaskOrder,
  isLoading: false,
  isError: false,
};

describe('PrimeUIRequestSITDestAddressChangeForm', () => {
  it('renders the address change request form', async () => {
    usePrimeSimulatorGetMove.mockReturnValue(moveReturnValue);
    renderWithProviders(<AcknowledgeMove moveTaskOrder={moveTaskOrder} />);

    const moveCodeElement = screen.getByText('Move Code:');
    expect(moveCodeElement).toBeInTheDocument();
    expect(moveCodeElement.nextSibling).toHaveTextContent('DEPPRQ');

    const moveIdElement = screen.getByText('Move Id:');
    expect(moveIdElement).toBeInTheDocument();
    expect(moveIdElement.nextSibling).toHaveTextContent('1');
  });
});
