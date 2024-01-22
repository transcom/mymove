import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect'; // For expect assertions

import { mockMovesPCS } from '../MultiMovesTestData';

import MultiMovesMoveContainer from './MultiMovesMoveContainer';

describe('MultiMovesMoveContainer', () => {
  const mockCurrentMoves = mockMovesPCS.currentMove;
  const mockPreviousMoves = mockMovesPCS.previousMoves;

  it('renders current move list correctly', () => {
    render(<MultiMovesMoveContainer moves={mockCurrentMoves} />);

    expect(screen.getByText('#MOVECO')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Go to Move' })).toBeInTheDocument();
  });

  it('renders previous move list correctly', () => {
    render(<MultiMovesMoveContainer moves={mockPreviousMoves} />);

    expect(screen.getByText('#SAMPLE')).toBeInTheDocument();
    expect(screen.getByText('#EXAMPL')).toBeInTheDocument();
    // should render two buttons
    expect(screen.getAllByRole('button', { name: 'Download' })).toHaveLength(2);
  });
});
