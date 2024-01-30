import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect'; // For expect assertions

import { mockMovesPCS } from '../MultiMovesTestData';

import MultiMovesMoveContainer from './MultiMovesMoveContainer';

import { MockProviders } from 'testUtils';

describe('MultiMovesMoveContainer', () => {
  const mockCurrentMoves = mockMovesPCS.currentMove;
  const mockPreviousMoves = mockMovesPCS.previousMoves;

  it('renders current move list correctly', () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockCurrentMoves} />
      </MockProviders>,
    );

    expect(screen.getByTestId('move-info-container')).toBeInTheDocument();
    expect(screen.getByText('#MOVECO')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Go to Move' })).toBeInTheDocument();
  });

  it('renders previous move list correctly', () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockPreviousMoves} />
      </MockProviders>,
    );

    expect(screen.queryByText('#SAMPLE')).toBeInTheDocument();
    expect(screen.queryByText('#EXAMPL')).toBeInTheDocument();
    expect(screen.getAllByRole('button', { name: 'Download' })).toHaveLength(2);
  });

  it('expands and collapses moves correctly', () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockCurrentMoves} />
      </MockProviders>,
    );

    // Initially, the move details should not be visible
    expect(screen.queryByText('Shipment')).not.toBeInTheDocument();

    fireEvent.click(screen.getByTestId('expand-icon'));

    // Now, the move details should be visible
    expect(screen.getByText('Shipments')).toBeInTheDocument();

    fireEvent.click(screen.getByTestId('expand-icon'));

    // The move details should be hidden again
    expect(screen.queryByText('Shipments')).not.toBeInTheDocument();
  });
});
