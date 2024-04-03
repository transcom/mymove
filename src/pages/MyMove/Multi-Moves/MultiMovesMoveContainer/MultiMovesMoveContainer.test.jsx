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
    // TODO commenting this out for now
    // expect(screen.getAllByRole('button', { name: 'Download' })).toHaveLength(2);
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

  it('renders Go to Move & Download buttons for current move', () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockCurrentMoves} />
      </MockProviders>,
    );

    expect(screen.getByTestId('headerBtns')).toBeInTheDocument();
    // TODO commenting this out for now
    // expect(screen.getByRole('button', { name: 'Download' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Go to Move' })).toBeInTheDocument();
  });

  it('renders Go to Move & Download buttons for previous moves exceeding one', () => {
    render(
      <MockProviders>
        <MultiMovesMoveContainer moves={mockPreviousMoves} />
      </MockProviders>,
    );

    // Check for the container that holds the buttons - there should be 2
    const headerBtnsElements = screen.getAllByTestId('headerBtns');
    expect(headerBtnsElements).toHaveLength(2);

    // TODO commenting these out for now
    // Check for Download buttons - there should be 2
    // const downloadButtons = screen.getAllByRole('button', { name: 'Download' });
    // expect(downloadButtons).toHaveLength(2);

    // Check for Go to Move buttons - there should be 2
    const goToMoveButtons = screen.getAllByRole('button', { name: 'Go to Move' });
    expect(goToMoveButtons).toHaveLength(2);
  });
});
