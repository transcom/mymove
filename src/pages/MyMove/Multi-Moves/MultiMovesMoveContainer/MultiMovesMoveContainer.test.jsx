import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect'; // For expect assertions

import { movesPCS } from '../MultiMovesTestData';

import MultiMovesMoveContainer from './MultiMovesMoveContainer';

describe('MultiMovesMoveContainer', () => {
  const mockMove = movesPCS.previousMoves;

  it('renders move list correctly', () => {
    render(<MultiMovesMoveContainer move={mockMove} />);

    // Check if move codes and status are rendered
    expect(screen.getByText('#SAMPLE')).toBeInTheDocument();
    expect(screen.getByText('#EXAMPL')).toBeInTheDocument();
  });
});
